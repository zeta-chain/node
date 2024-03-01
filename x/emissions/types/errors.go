package types

import errorsmod "cosmossdk.io/errors"

// DONTCOVER

var (
	ErrEmissionsNotFound                   = errorsmod.Register(ModuleName, 1000, "Emissions not found")
	ErrUnableToWithdrawEmissions           = errorsmod.Register(ModuleName, 1002, "Unable to withdraw emissions")
	ErrInvalidAddress                      = errorsmod.Register(ModuleName, 1003, "Invalid address")
	ErrRewardsPoolDoesNotHaveEnoughBalance = errorsmod.Register(ModuleName, 1004, "Rewards pool does not have enough balance")
	ErrInvalidAmount                       = errorsmod.Register(ModuleName, 1005, "Invalid amount")
)
