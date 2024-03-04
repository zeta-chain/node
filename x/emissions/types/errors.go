package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrEmissionsNotFound                   = errorsmod.Register(ModuleName, 1000, "emissions not found")
	ErrUnableToWithdrawEmissions           = errorsmod.Register(ModuleName, 1002, "unable to withdraw emissions")
	ErrInvalidAddress                      = errorsmod.Register(ModuleName, 1003, "invalid address")
	ErrRewardsPoolDoesNotHaveEnoughBalance = errorsmod.Register(ModuleName, 1004, "rewards pool does not have enough balance")
	ErrInvalidAmount                       = errorsmod.Register(ModuleName, 1005, "invalid amount")
)
