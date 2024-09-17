package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
)

func TestPrecompilesBank(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	spender, bankAddr := r.EVMAddress(), bank.ContractAddress

	// Create a bank contract caller.
	bankContract, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Deposit and approve 50 WZETA for the test.
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.DepositAndApproveWZeta(approveAmount)

	// Initial WZETA spender balance should be 50.
	initialBalance, err := r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	require.NoError(r, err, "Error approving allowance for bank contract")
	require.EqualValues(r, approveAmount, initialBalance, "spender balance should be 50")

	// Initial cosmos coin spender balance should be 0.
	retBalanceOf, err := bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.WZetaAddr, spender)
	require.NoError(r, err, "Error calling bank.balanceOf()")
	require.EqualValues(r, uint64(0), retBalanceOf.Uint64(), "Initial cosmos coins balance has to be 0")

	// Allow the bank contract to spend 25 WZeta tokens.
	tx, err := r.WZeta.Approve(r.ZEVMAuth, bankAddr, big.NewInt(25))
	require.NoError(r, err, "Error approving allowance for bank contract")
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, uint64(1), receipt.Status, "approve allowance tx failed")

	// Check the allowance of the bank in WZeta tokens. Should be 25.
	balance, err := r.WZeta.Allowance(&bind.CallOpts{Context: r.Ctx}, spender, bankAddr)
	require.NoError(r, err, "Error retrieving bank allowance")
	require.EqualValues(r, uint64(25), balance.Uint64(), "Error allowance for bank contract")

	// Call Deposit with 25 coins.
	tx, err = bankContract.Deposit(r.ZEVMAuth, r.WZetaAddr, big.NewInt(25))
	require.NoError(r, err, "Error calling bank.deposit()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// Deposit event should be emitted.
	depositEvent, err := bankContract.ParseDeposit(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, big.NewInt(25).Uint64(), depositEvent.Amount.Uint64())
	require.Equal(r, common.BytesToAddress(spender.Bytes()), depositEvent.Zrc20Depositor)
	require.Equal(r, r.WZetaAddr, depositEvent.Zrc20Token)

	// After deposit, cosmos coin spender balance should be 25.
	retBalanceOf, err = bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.WZetaAddr, spender)
	require.NoError(r, err, "Error calling bank.balanceOf()")
	require.EqualValues(r, uint64(25), retBalanceOf.Uint64(), "balanceOf result has to be 25")

	// After deposit, WZeta spender balance should be 25 less than initial.
	finalBalance, err := r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	require.NoError(r, err, "Error retrieving final owner balance")
	require.EqualValues(
		r,
		initialBalance.Uint64()-25, // expected
		finalBalance.Uint64(),      // actual
		"Final balance should be initial - 25",
	)

	// After deposit, WZeta bank balance should be 25.
	balance, err = r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, bankAddr)
	require.NoError(r, err, "Error retrieving bank's balance")
	require.EqualValues(r, uint64(25), balance.Uint64(), "Wrong locked WZeta amount in bank contract")

	// Withdraw 15 coins to spender.
	tx, err = bankContract.Withdraw(r.ZEVMAuth, r.WZetaAddr, big.NewInt(15))
	require.NoError(r, err, "Error calling bank.withdraw()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// Withdraw event should be emitted.
	withdrawEvent, err := bankContract.ParseWithdraw(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, big.NewInt(15).Uint64(), withdrawEvent.Amount.Uint64())
	require.Equal(r, common.BytesToAddress(spender.Bytes()), withdrawEvent.Zrc20Withdrawer)
	require.Equal(r, r.WZetaAddr, withdrawEvent.Zrc20Token)

	// After withdraw, WZeta spender balance should be only 10 less than initial. (25 - 15 = 10e2e/e2etests/test_precompiles_bank.go )
	afterWithdraw, err := r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	require.NoError(r, err, "Error retrieving final owner balance")
	require.EqualValues(
		r,
		initialBalance.Uint64()-10, // expected
		afterWithdraw.Uint64(),     // actual
		"Balance after withdraw should be initial - 10",
	)

	// After withdraw, cosmos coin spender balance should be 10.
	retBalanceOf, err = bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.WZetaAddr, spender)
	require.NoError(r, err, "Error calling bank.balanceOf()")
	require.EqualValues(r, uint64(10), retBalanceOf.Uint64(), "balanceOf result has to be 10")

	// Final WZETA bank balance should be 10.
	balance, err = r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, bankAddr)
	require.NoError(r, err, "Error retrieving bank's allowance")
	require.EqualValues(r, uint64(10), balance.Uint64(), "Wrong locked WZeta amount in bank contract")
}
