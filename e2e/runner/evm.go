package runner

import (
	"log"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/pkg/proofs/ethereum"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
)

var blockHeaderETHTimeout = 5 * time.Minute

// WaitForTxReceiptOnEvm waits for a tx receipt on EVM
func (r *E2ERunner) WaitForTxReceiptOnEvm(tx *ethtypes.Transaction) {
	r.Lock()
	defer r.Unlock()

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)
}

// MintERC20OnEvm mints ERC20 on EVM
// amount is a multiple of 1e18
func (r *E2ERunner) MintERC20OnEvm(amountERC20 int64) {
	r.Lock()
	defer r.Unlock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountERC20))

	tx, err := r.ERC20.Mint(r.EVMAuth, amount)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())
}

// SendERC20OnEvm sends ERC20 to an address on EVM
// this allows the ERC20 contract deployer to funds other accounts on EVM
// amountERC20 is a multiple of 1e18
func (r *E2ERunner) SendERC20OnEvm(address ethcommon.Address, amountERC20 int64) *ethtypes.Transaction {
	// the deployer might be sending ERC20 in different goroutines
	r.Lock()
	defer r.Unlock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountERC20))

	// transfer
	tx, err := r.ERC20.Transfer(r.EVMAuth, address, amount)
	require.NoError(r, err)

	return tx
}

func (r *E2ERunner) DepositERC20() ethcommon.Hash {
	r.Logger.Print("⏳ depositing ERC20 into ZEVM")

	return r.DepositERC20WithAmountAndMessage(r.EVMAddress(), big.NewInt(1e18), []byte{})
}

func (r *E2ERunner) DepositERC20WithAmountAndMessage(to ethcommon.Address, amount *big.Int, msg []byte) ethcommon.Hash {
	// reset allowance, necessary for USDT
	tx, err := r.ERC20.Approve(r.EVMAuth, r.ERC20CustodyAddr, big.NewInt(0))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = r.ERC20.Approve(r.EVMAuth, r.ERC20CustodyAddr, amount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = r.ERC20Custody.Deposit(r.EVMAuth, to.Bytes(), r.ERC20Addr, amount, msg)
	require.NoError(r, err)

	r.Logger.Info("TX: %v", tx)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("Deposit receipt tx hash: %s, status %d", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("Deposited event:")
		r.Logger.Info("  Recipient address: %x", event.Recipient)
		r.Logger.Info("  ERC20 address: %s", event.Asset.Hex())
		r.Logger.Info("  Amount: %d", event.Amount)
		r.Logger.Info("  Message: %x", event.Message)
	}
	return tx.Hash()
}

// DepositEther sends Ethers into ZEVM
func (r *E2ERunner) DepositEther(testHeader bool) ethcommon.Hash {
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100)) // 100 eth
	return r.DepositEtherWithAmount(testHeader, amount)
}

// DepositEtherWithAmount sends Ethers into ZEVM
func (r *E2ERunner) DepositEtherWithAmount(testHeader bool, amount *big.Int) ethcommon.Hash {
	r.Logger.Print("⏳ depositing Ethers into ZEVM")

	signedTx, err := r.SendEther(r.TSSAddress, amount, nil)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*signedTx, "send to TSS")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, "deposit failed")

	r.Logger.EVMReceipt(*receipt, "send to TSS")

	// due to the high block throughput in localnet, ZetaClient might catch up slowly with the blocks
	// to optimize block header proof test, this test is directly executed here on the first deposit instead of having a separate test
	if testHeader {
		r.ProveEthTransaction(receipt)
	}

	return signedTx.Hash()
}

// SendEther sends ethers to the TSS on EVM
func (r *E2ERunner) SendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	evmClient := r.EVMClient

	nonce, err := evmClient.PendingNonceAt(r.Ctx, r.EVMAddress())
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(30000) // in units
	gasPrice, err := evmClient.SuggestGasPrice(r.Ctx)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.ChainID(r.Ctx)
	if err != nil {
		return nil, err
	}

	deployerPrivkey, err := r.Account.PrivateKey()
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// ProveEthTransaction proves an ETH transaction on ZetaChain
func (r *E2ERunner) ProveEthTransaction(receipt *ethtypes.Receipt) {
	startTime := time.Now()

	txHash := receipt.TxHash
	blockHash := receipt.BlockHash

	// #nosec G115 test - always in range
	txIndex := int(receipt.TransactionIndex)

	block, err := r.EVMClient.BlockByHash(r.Ctx, blockHash)
	require.NoError(r, err)

	for {
		// check timeout
		reachedTimeout := time.Since(startTime) > blockHeaderETHTimeout
		require.False(r, reachedTimeout, "timeout waiting for block header")

		_, err := r.LightclientClient.BlockHeader(r.Ctx, &lightclienttypes.QueryGetBlockHeaderRequest{
			BlockHash: blockHash.Bytes(),
		})
		if err != nil {
			r.Logger.Info("WARN: block header not found; retrying... error: %s", err.Error())
		} else {
			r.Logger.Info("OK: block header found")
			break
		}

		time.Sleep(2 * time.Second)
	}

	trie := ethereum.NewTrie(block.Transactions())
	require.Equal(r, trie.Hash(), block.Header().TxHash, "tx root hash & block tx root mismatch")

	txProof, err := trie.GenerateProof(txIndex)
	require.NoError(r, err, "error generating txProof")

	val, err := txProof.Verify(block.TxHash(), txIndex)
	require.NoError(r, err, "error verifying txProof")

	var txx ethtypes.Transaction
	require.NoError(r, txx.UnmarshalBinary(val))

	res, err := r.LightclientClient.Prove(r.Ctx, &lightclienttypes.QueryProveRequest{
		BlockHash: blockHash.Hex(),
		TxIndex:   int64(txIndex),
		TxHash:    txHash.Hex(),
		Proof:     proofs.NewEthereumProof(txProof),
		ChainId:   chains.GoerliLocalnet.ChainId,
	})

	// FIXME: @lumtis: don't do this in production
	require.NoError(r, err)
	require.True(r, res.Valid, "txProof invalid")

	r.Logger.Info("OK: txProof verified")
}

// AnvilMineBlocks mines blocks on Anvil localnet
// the block time is provided in seconds
// the method returns a function to stop the mining
func (r *E2ERunner) AnvilMineBlocks(url string, blockTime int) (func(), error) {
	stop := make(chan struct{})

	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				time.Sleep(time.Duration(blockTime) * time.Second)

				var result interface{}
				err = client.CallContext(r.Ctx, &result, "evm_mine")
				if err != nil {
					log.Fatalf("Failed to mine a new block: %v", err)
				}
			}
		}
	}()
	return func() {
		close(stop)
	}, nil
}
