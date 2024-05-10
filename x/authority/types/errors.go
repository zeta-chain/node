package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrUnauthorized      = errorsmod.Register(ModuleName, 1102, "sender not authorized")
	ErrSigners           = errorsmod.Register(ModuleName, 1103, "policy transactions must have only one signer")
	ErrMsgNotAuthorized  = errorsmod.Register(ModuleName, 1104, "msg type is not authorized")
	ErrPoliciesNotFound  = errorsmod.Register(ModuleName, 1105, "policies not found")
	ErrSignerDoesntMatch = errorsmod.Register(ModuleName, 1106, "signer doesn't match required policy")
)
