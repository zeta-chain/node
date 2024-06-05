package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrUnauthorized              = errorsmod.Register(ModuleName, 1102, "sender not authorized")
	ErrInvalidAuthorizationList  = errorsmod.Register(ModuleName, 1103, "invalid authorization list")
	ErrAuthorizationNotFound     = errorsmod.Register(ModuleName, 1104, "authorization not found")
	ErrAuthorizationListNotFound = errorsmod.Register(ModuleName, 1105, "authorization list not found")
)
