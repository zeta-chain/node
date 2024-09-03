package runner

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

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
func (r *E2ERunner) DepositEther() ethcommon.Hash {
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100)) // 100 eth
	return r.DepositEtherWithAmount(amount)
}

// DepositEtherWithAmount sends Ethers into ZEVM
func (r *E2ERunner) DepositEtherWithAmount(amount *big.Int) ethcommon.Hash {
	r.Logger.Print("⏳ depositing Ethers into ZEVM")

	signedTx, err := r.SendEther(r.TSSAddress, amount, nil)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*signedTx, "send to TSS")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, "deposit failed")

	r.Logger.EVMReceipt(*receipt, "send to TSS")

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

// ApproveERC20OnEVM approves ERC20 on EVM to a specific address
// check if allowance is zero before calling this method
// allow a high amount to avoid multiple approvals
func (r *E2ERunner) ApproveERC20OnEVM(allowed ethcommon.Address) {
	allowance, err := r.ERC20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), r.GatewayEVMAddr)
	require.NoError(r, err)

	// approve 1M*1e18 if allowance is zero
	if allowance.Cmp(big.NewInt(0)) == 0 {
		tx, err := r.ERC20.Approve(r.EVMAuth, allowed, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}
}

// ApproveETHZRC20 approves ETH ZRC20 on EVM to a specific address
// check if allowance is zero before calling this method
// allow a high amount to avoid multiple approvals
func (r *E2ERunner) ApproveETHZRC20(allowed ethcommon.Address) {
	allowance, err := r.ETHZRC20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), allowed)
	require.NoError(r, err)

	// approve 1M*1e18 if allowance is below 1k
	thousand := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
	if allowance.Cmp(thousand) < 0 {
		tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, allowed, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}
}

// ApproveERC20ZRC20 approves ERC20 ZRC20 on EVM to a specific address
// check if allowance is zero before calling this method
// allow a high amount to avoid multiple approvals
func (r *E2ERunner) ApproveERC20ZRC20(allowed ethcommon.Address) {
	allowance, err := r.ERC20ZRC20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), allowed)
	require.NoError(r, err)

	// approve 1M*1e18 if allowance is below 1k
	thousand := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
	if allowance.Cmp(thousand) < 0 {
		tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, allowed, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}
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
