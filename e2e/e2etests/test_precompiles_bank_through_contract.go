package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testbank"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
)

func TestPrecompilesBankThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	var (
		spender        = r.EVMAddress()
		bankAddress    = bank.ContractAddress
		zrc20Address   = r.ERC20ZRC20Addr
		oneThousand    = big.NewInt(1e3)
		oneThousandOne = big.NewInt(1001)
		fiveHundred    = big.NewInt(500)
		fiveHundredOne = big.NewInt(501)
		zero           = big.NewInt(0)
	)

	// Get ERC20ZRC20.
	txHash := r.DepositERC20WithAmountAndMessage(r.EVMAddress(), oneThousand, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	bankPrecompileCaller, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank precompile caller")

	// Deploy the TestBank. Ensure the transaction is successful.
	_, tx, testBank, err := testbank.DeployTestBank(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Deployment of TestBank contract failed")

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit

		// Reset the allowance to 0; this is needed when running upgrade tests where this test runs twice.
		approveAllowance(r, bank.ContractAddress, big.NewInt(0))

		// Reset balance to 0; this is needed when running upgrade tests where this test runs twice.
		tx, err = r.ERC20ZRC20.Transfer(
			r.ZEVMAuth,
			common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
			oneThousand,
		)
		require.NoError(r, err)
		receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, "Resetting balance failed")
	}()

	// Check initial balances.
	balanceShouldBe(r, zero, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, oneThousand, checkZRC20Balance(r, spender))
	balanceShouldBe(r, zero, checkZRC20Balance(r, bankAddress))

	// Deposit without previous alllowance should fail.
	receipt = depositThroughTestBank(r, testBank, zrc20Address, oneThousand)
	utils.RequiredTxFailed(r, receipt, "Deposit ERC20ZRC20 without allowance should fail")

	// Check balances, should be the same.
	balanceShouldBe(r, zero, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, oneThousand, checkZRC20Balance(r, spender))
	balanceShouldBe(r, zero, checkZRC20Balance(r, bankAddress))

	// Allow 500 ZRC20 to bank precompile.
	approveAllowance(r, bankAddress, fiveHundred)

	// Deposit 501 ERC20ZRC20 tokens to the bank contract, through TestBank.
	// It's higher than allowance but lower than balance, should fail.
	receipt = depositThroughTestBank(r, testBank, zrc20Address, fiveHundredOne)
	utils.RequiredTxFailed(r, receipt, "Depositting an amount higher than allowed should fail")

	// Balances shouldn't change.
	balanceShouldBe(r, zero, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, oneThousand, checkZRC20Balance(r, spender))
	balanceShouldBe(r, zero, checkZRC20Balance(r, bankAddress))

	// Allow 1000 ZRC20 to bank precompile.
	approveAllowance(r, bankAddress, oneThousand)

	// Deposit 1001 ERC20ZRC20 tokens to the bank contract.
	// It's higher than spender balance but within approved allowance, should fail.
	receipt = depositThroughTestBank(r, testBank, zrc20Address, oneThousandOne)
	utils.RequiredTxFailed(r, receipt, "Depositting an amount higher than balance should fail")

	// Balances shouldn't change.
	balanceShouldBe(r, zero, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, oneThousand, checkZRC20Balance(r, spender))
	balanceShouldBe(r, zero, checkZRC20Balance(r, bankAddress))

	// Deposit 500 ERC20ZRC20 tokens to the bank contract, it's within allowance and balance. Should pass.
	receipt = depositThroughTestBank(r, testBank, zrc20Address, fiveHundred)
	utils.RequireTxSuccessful(r, receipt, "Depositting a correct amount should pass")

	// Balances should be transferred. Bank now locks 500 ZRC20 tokens.
	balanceShouldBe(r, fiveHundred, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, fiveHundred, checkZRC20Balance(r, spender))
	balanceShouldBe(r, fiveHundred, checkZRC20Balance(r, bankAddress))

	// Check the deposit event.
	eventDeposit, err := bankPrecompileCaller.ParseDeposit(*receipt.Logs[0])
	require.NoError(r, err, "Parse Deposit event")
	require.Equal(r, r.EVMAddress(), eventDeposit.Zrc20Depositor, "Deposit event token should be r.EVMAddress()")
	require.Equal(r, r.ERC20ZRC20Addr, eventDeposit.Zrc20Token, "Deposit event token should be ERC20ZRC20Addr")
	require.Equal(r, fiveHundred, eventDeposit.Amount, "Deposit event amount should be 500")

	// Should faild to withdraw more than cosmos balance.
	receipt = withdrawThroughTestBank(r, testBank, zrc20Address, fiveHundredOne)
	utils.RequiredTxFailed(r, receipt, "Withdrawing an amount higher than balance should fail")

	// Balances shouldn't change.
	balanceShouldBe(r, fiveHundred, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, fiveHundred, checkZRC20Balance(r, spender))
	balanceShouldBe(r, fiveHundred, checkZRC20Balance(r, bankAddress))

	// Try to withdraw 500 ERC20ZRC20 tokens. Should pass.
	receipt = withdrawThroughTestBank(r, testBank, zrc20Address, fiveHundred)
	utils.RequireTxSuccessful(r, receipt, "Withdraw correct amount should pass")

	// Balances should be reverted to initial state.
	balanceShouldBe(r, zero, checkCosmosBalanceThroughBank(r, testBank, zrc20Address, spender))
	balanceShouldBe(r, oneThousand, checkZRC20Balance(r, spender))
	balanceShouldBe(r, zero, checkZRC20Balance(r, bankAddress))

	// Check the withdraw event.
	eventWithdraw, err := bankPrecompileCaller.ParseWithdraw(*receipt.Logs[0])
	require.NoError(r, err, "Parse Withdraw event")
	require.Equal(r, r.EVMAddress(), eventWithdraw.Zrc20Withdrawer, "Withdrawer should be r.EVMAddress()")
	require.Equal(r, r.ERC20ZRC20Addr, eventWithdraw.Zrc20Token, "Withdraw event token should be ERC20ZRC20Addr")
	require.Equal(r, fiveHundred, eventWithdraw.Amount, "Withdraw event amount should be 500")
}

func approveAllowance(r *runner.E2ERunner, target common.Address, amount *big.Int) {
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, target, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Approve ERC20ZRC20 allowance tx failed")
}

func balanceShouldBe(r *runner.E2ERunner, expected *big.Int, balance *big.Int) {
	require.Equal(r, expected.Uint64(), balance.Uint64(), "Balance should be %d, got: %d", expected, balance.Uint64())
}

func checkZRC20Balance(r *runner.E2ERunner, target common.Address) *big.Int {
	bankZRC20Balance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, target)
	require.NoError(r, err, "Call ERC20ZRC20.BalanceOf")
	return bankZRC20Balance
}

func checkCosmosBalanceThroughBank(
	r *runner.E2ERunner,
	bank *testbank.TestBank,
	zrc20, target common.Address,
) *big.Int {
	balance, err := bank.BalanceOf(&bind.CallOpts{Context: r.Ctx, From: r.ZEVMAuth.From}, zrc20, target)
	require.NoError(r, err)
	return balance
}

func depositThroughTestBank(
	r *runner.E2ERunner,
	bank *testbank.TestBank,
	zrc20Address common.Address,
	amount *big.Int,
) *types.Receipt {
	tx, err := bank.Deposit(r.ZEVMAuth, zrc20Address, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	return receipt
}

func withdrawThroughTestBank(
	r *runner.E2ERunner,
	bank *testbank.TestBank,
	zrc20Address common.Address,
	amount *big.Int,
) *types.Receipt {
	tx, err := bank.Withdraw(r.ZEVMAuth, zrc20Address, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	return receipt
}
