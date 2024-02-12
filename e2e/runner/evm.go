package runner

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	"github.com/zeta-chain/zetacore/e2e/utils"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var blockHeaderETHTimeout = 5 * time.Minute

// WaitForTxReceiptOnEvm waits for a tx receipt on EVM
func (runner *E2ERunner) WaitForTxReceiptOnEvm(tx *ethtypes.Transaction) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
}

// MintUSDTOnEvm mints USDT on EVM
// amountUSDT is a multiple of 1e18
func (runner *E2ERunner) MintUSDTOnEvm(amountUSDT int64) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountUSDT))

	tx, err := runner.USDTERC20.Mint(runner.GoerliAuth, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("mint failed")
	}
	runner.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())
}

// SendUSDTOnEvm sends USDT to an address on EVM
// this allows the USDT contract deployer to funds other accounts on EVM
// amountUSDT is a multiple of 1e18
func (runner *E2ERunner) SendUSDTOnEvm(address ethcommon.Address, amountUSDT int64) *ethtypes.Transaction {
	// the deployer might be sending USDT in different goroutines
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountUSDT))

	// transfer
	tx, err := runner.USDTERC20.Transfer(runner.GoerliAuth, address, amount)
	if err != nil {
		panic(err)
	}
	return tx
}

func (runner *E2ERunner) DepositERC20() ethcommon.Hash {
	runner.Logger.Print("⏳ depositing ERC20 into ZEVM")

	return runner.DepositERC20WithAmountAndMessage(big.NewInt(1e18), []byte{})
}

func (runner *E2ERunner) DepositERC20WithAmountAndMessage(amount *big.Int, msg []byte) ethcommon.Hash {
	// reset allowance, necessary for USDT
	tx, err := runner.USDTERC20.Approve(runner.GoerliAuth, runner.ERC20CustodyAddr, big.NewInt(0))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	runner.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = runner.USDTERC20.Approve(runner.GoerliAuth, runner.ERC20CustodyAddr, amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	runner.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = runner.ERC20Custody.Deposit(runner.GoerliAuth, runner.DeployerAddress.Bytes(), runner.USDTERC20Addr, amount, msg)
	runner.Logger.Print("TX: %v", tx)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	runner.Logger.Info("Deposit receipt tx hash: %s, status %d", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := runner.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		runner.Logger.Info("Deposited event:")
		runner.Logger.Info("  Recipient address: %x", event.Recipient)
		runner.Logger.Info("  ERC20 address: %s", event.Asset.Hex())
		runner.Logger.Info("  Amount: %d", event.Amount)
		runner.Logger.Info("  Message: %x", event.Message)
	}
	return tx.Hash()
}

// DepositEther sends Ethers into ZEVM
func (runner *E2ERunner) DepositEther(testHeader bool) ethcommon.Hash {
	return runner.DepositEtherWithAmount(testHeader, big.NewInt(1000000000000000000)) // in wei (1 eth)
}

// DepositEtherWithAmount sends Ethers into ZEVM
func (runner *E2ERunner) DepositEtherWithAmount(testHeader bool, amount *big.Int) ethcommon.Hash {
	runner.Logger.Print("⏳ depositing Ethers into ZEVM")

	signedTx, err := runner.SendEther(runner.TSSAddress, amount, nil)
	if err != nil {
		panic(err)
	}
	runner.Logger.EVMTransaction(*signedTx, "send to TSS")

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, signedTx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	runner.Logger.EVMReceipt(*receipt, "send to TSS")

	// due to the high block throughput in localnet, ZetaClient might catch up slowly with the blocks
	// to optimize block header proof test, this test is directly executed here on the first deposit instead of having a separate test
	if testHeader {
		runner.ProveEthTransaction(receipt)
	}

	return signedTx.Hash()
}

// SendEther sends ethers to the TSS on Goerli
func (runner *E2ERunner) SendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	goerliClient := runner.GoerliClient

	nonce, err := goerliClient.PendingNonceAt(runner.Ctx, runner.DeployerAddress)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(30000) // in units
	gasPrice, err := goerliClient.SuggestGasPrice(runner.Ctx)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(nonce, runner.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(runner.Ctx)
	if err != nil {
		return nil, err
	}

	deployerPrivkey, err := crypto.HexToECDSA(runner.DeployerPrivateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = goerliClient.SendTransaction(runner.Ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// ProveEthTransaction proves an ETH transaction on ZetaChain
func (runner *E2ERunner) ProveEthTransaction(receipt *ethtypes.Receipt) {
	startTime := time.Now()

	txHash := receipt.TxHash
	blockHash := receipt.BlockHash

	// #nosec G701 test - always in range
	txIndex := int(receipt.TransactionIndex)

	block, err := runner.GoerliClient.BlockByHash(runner.Ctx, blockHash)
	if err != nil {
		panic(err)
	}
	for {
		// check timeout
		if time.Since(startTime) > blockHeaderETHTimeout {
			panic("timeout waiting for block header")
		}

		_, err := runner.ObserverClient.GetBlockHeaderByHash(runner.Ctx, &observertypes.QueryGetBlockHeaderByHashRequest{
			BlockHash: blockHash.Bytes(),
		})
		if err != nil {
			runner.Logger.Info("WARN: block header not found; retrying... error: %s", err.Error())
		} else {
			runner.Logger.Info("OK: block header found")
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
	res, err := runner.ObserverClient.Prove(runner.Ctx, &observertypes.QueryProveRequest{
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
	runner.Logger.Info("OK: txProof verified")
}
