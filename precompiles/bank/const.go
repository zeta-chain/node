package bank

const (
	// ZEVM cosmos coins prefix.
	ZEVMDenom = "zevm/"

	// Write methods.
	DepositMethodName = "deposit"
	DepositMethodGas  = 200_000
	DepositEventName  = "Deposit"

	WithdrawMethodName = "withdraw"
	WithdrawMethodGas  = 200_000
	WithdrawEventName  = "Withdraw"

	// Read methods.
	BalanceOfMethodName = "balanceOf"
	BalanceOfGas        = 10_000

	// Default gas for unknown methods.
	DefaultGas = 0
)
