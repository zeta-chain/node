package types

import errorsmod "cosmossdk.io/errors"

// DONTCOVER

var (
	ErrEmissionsNotFound                   = errorsmod.Register(ModuleName, 1000, "Emissions not found")
	ErrNotEnoughEmissionsAvailable         = errorsmod.Register(ModuleName, 1001, "Not enough emissions available to withdraw")
	ErrUnableToCreateWithdrawEmissions     = errorsmod.Register(ModuleName, 1002, "Unable to create withdraw emissions")
	ErrInvalidAddress                      = errorsmod.Register(ModuleName, 1003, "Invalid address")
	ErrRewardsPoolDoesNotHaveEnoughBalance = errorsmod.Register(ModuleName, 1004, "Rewards pool does not have enough balance")

	ErrInvalidAmount = errorsmod.Register(ModuleName, 1005, "Invalid amount")
)
