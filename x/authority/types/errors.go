package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrUnauthorized              = errorsmod.Register(ModuleName, 1102, "sender not authorized")
	ErrInvalidAuthorizationList  = errorsmod.Register(ModuleName, 1103, "invalid authorization list")
	ErrAuthorizationNotFound     = errorsmod.Register(ModuleName, 1104, "authorization not found")
	ErrAuthorizationListNotFound = errorsmod.Register(ModuleName, 1105, "authorization list not found")
	ErrSigners                   = errorsmod.Register(ModuleName, 1106, "policy transactions must have only one signer")
	ErrMsgNotAuthorized          = errorsmod.Register(ModuleName, 1107, "msg type is not authorized")
	ErrPoliciesNotFound          = errorsmod.Register(ModuleName, 1108, "policies not found")
	ErrSignerDoesntMatch         = errorsmod.Register(ModuleName, 1109, "signer doesn't match required policy")
	ErrInvalidPolicyType         = errorsmod.Register(ModuleName, 1110, "invalid policy type")
)
