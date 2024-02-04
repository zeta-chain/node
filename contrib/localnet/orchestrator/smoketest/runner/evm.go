package runner

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var blockHeaderETHTimeout = 5 * time.Minute

// WaitForTxReceiptOnEvm waits for a tx receipt on EVM
func (sm *SmokeTestRunner) WaitForTxReceiptOnEvm(tx *ethtypes.Transaction) {
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
}

// MintUSDTOnEvm mints USDT on EVM
// amountUSDT is a multiple of 1e18
func (sm *SmokeTestRunner) MintUSDTOnEvm(amountUSDT int64) {
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountUSDT))

	tx, err := sm.USDTERC20.Mint(sm.GoerliAuth, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("mint failed")
	}
	sm.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())
}

// SendUSDTOnEvm sends USDT to an address on EVM
// this allows the USDT contract deployer to funds other accounts on EVM
// amountUSDT is a multiple of 1e18
func (sm *SmokeTestRunner) SendUSDTOnEvm(address ethcommon.Address, amountUSDT int64) *ethtypes.Transaction {
	// the deployer might be sending USDT in different goroutines
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountUSDT))

	// transfer
	tx, err := sm.USDTERC20.Transfer(sm.GoerliAuth, address, amount)
	if err != nil {
		panic(err)
	}
	return tx
}

func (sm *SmokeTestRunner) DepositERC20() ethcommon.Hash {
	sm.Logger.Print("⏳ depositing ERC20 into ZEVM")

	return sm.DepositERC20WithAmountAndMessage(big.NewInt(1e18), []byte{})
}

func (sm *SmokeTestRunner) DepositERC20WithAmountAndMessage(amount *big.Int, msg []byte) ethcommon.Hash {
	// reset allowance, necessary for USDT
	tx, err := sm.USDTERC20.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, big.NewInt(0))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = sm.USDTERC20.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(sm.GoerliAuth, sm.DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, msg)
	sm.Logger.Print("TX: %v", tx)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	sm.Logger.Info("Deposit receipt tx hash: %s, status %d", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info("Deposited event:")
		sm.Logger.Info("  Recipient address: %x", event.Recipient)
		sm.Logger.Info("  ERC20 address: %s", event.Asset.Hex())
		sm.Logger.Info("  Amount: %d", event.Amount)
		sm.Logger.Info("  Message: %x", event.Message)
	}
	return tx.Hash()
}

// DepositEther sends Ethers into ZEVM
func (sm *SmokeTestRunner) DepositEther(testHeader bool) ethcommon.Hash {
	return sm.DepositEtherWithAmount(testHeader, big.NewInt(1000000000000000000)) // in wei (1 eth)
}

// DepositEtherWithAmount sends Ethers into ZEVM
func (sm *SmokeTestRunner) DepositEtherWithAmount(testHeader bool, amount *big.Int) ethcommon.Hash {
	sm.Logger.Print("⏳ depositing Ethers into ZEVM")

	signedTx, err := sm.SendEther(sm.TSSAddress, amount, nil)
	if err != nil {
		panic(err)
	}
	sm.Logger.EVMTransaction(*signedTx, "send to TSS")

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	sm.Logger.EVMReceipt(*receipt, "send to TSS")

	// due to the high block throughput in localnet, ZetaClient might catch up slowly with the blocks
	// to optimize block header proof test, this test is directly executed here on the first deposit instead of having a separate test
	if testHeader {
		sm.ProveEthTransaction(receipt)
	}

	return signedTx.Hash()
}

// SendEther sends ethers to the TSS on Goerli
func (sm *SmokeTestRunner) SendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	goerliClient := sm.GoerliClient

	nonce, err := goerliClient.PendingNonceAt(sm.Ctx, sm.DeployerAddress)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(30000) // in units
	gasPrice, err := goerliClient.SuggestGasPrice(sm.Ctx)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(nonce, sm.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(sm.Ctx)
	if err != nil {
		return nil, err
	}

	deployerPrivkey, err := crypto.HexToECDSA(sm.DeployerPrivateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = goerliClient.SendTransaction(sm.Ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// ProveEthTransaction proves an ETH transaction on ZetaChain
func (sm *SmokeTestRunner) ProveEthTransaction(receipt *ethtypes.Receipt) {
	startTime := time.Now()

	txHash := receipt.TxHash
	blockHash := receipt.BlockHash

	// #nosec G701 smoketest - always in range
	txIndex := int(receipt.TransactionIndex)

	block, err := sm.GoerliClient.BlockByHash(sm.Ctx, blockHash)
	if err != nil {
		panic(err)
	}
	for {
		// check timeout
		if time.Since(startTime) > blockHeaderETHTimeout {
			panic("timeout waiting for block header")
		}

		_, err := sm.ObserverClient.GetBlockHeaderByHash(sm.Ctx, &observertypes.QueryGetBlockHeaderByHashRequest{
			BlockHash: blockHash.Bytes(),
		})
		if err != nil {
			sm.Logger.Info("WARN: block header not found; retrying... error: %s", err.Error())
		} else {
			sm.Logger.Info("OK: block header found")
			break
		}

		time.Sleep(2 * time.Second)
	}

	trie := ethereum.NewTrie(block.Transactions())
	if trie.Hash() != block.Header().TxHash {
		panic("tx root hash & block tx root mismatch")
	}
	txProof, err := trie.GenerateProof(txIndex)
	if err != nil {
		panic("error generating txProof")
	}
	val, err := txProof.Verify(block.TxHash(), txIndex)
	if err != nil {
		panic("error verifying txProof")
	}
	var txx ethtypes.Transaction
	err = txx.UnmarshalBinary(val)
	if err != nil {
		panic("error unmarshalling txProof'd tx")
	}
	res, err := sm.ObserverClient.Prove(sm.Ctx, &observertypes.QueryProveRequest{
		BlockHash: blockHash.Hex(),
		TxIndex:   int64(txIndex),
		TxHash:    txHash.Hex(),
		Proof:     common.NewEthereumProof(txProof),
		ChainId:   common.GoerliLocalnetChain().ChainId,
	})
	if err != nil {
		panic(err)
	}
	if !res.Valid {
		panic("txProof invalid") // FIXME: don't do this in production
	}
	sm.Logger.Info("OK: txProof verified")
}
