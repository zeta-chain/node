package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnableToAddVote = errorsmod.Register(ModuleName, 1100, "Unable to add vote ")
	ErrParamsThreshold = errorsmod.Register(ModuleName, 1101, "Threshold cannot be more than 1")
	ErrSupportedChains = errorsmod.Register(ModuleName, 1102, "Err chain not supported")
	ErrInvalidStatus   = errorsmod.Register(ModuleName, 1103, "Invalid Voting Status")

	ErrObserverNotPresent      = errorsmod.Register(ModuleName, 1105, "Observer for type and observation does not exist")
	ErrNotValidator            = errorsmod.Register(ModuleName, 1106, "User needs to be a validator before applying to become an observer")
	ErrValidatorStatus         = errorsmod.Register(ModuleName, 1107, "Corresponding validator needs to be bonded and not jailerd")
	ErrInvalidAddress          = errorsmod.Register(ModuleName, 1108, "Invalid Address")
	ErrSelfDelegation          = errorsmod.Register(ModuleName, 1109, "Self Delegation for operator not found")
	ErrCheckObserverDelegation = errorsmod.Register(ModuleName, 1110, "Observer delegation not sufficient")
	ErrNotAuthorizedPolicy     = errorsmod.Register(ModuleName, 1111, "Msg Sender is not the authorized policy")
	ErrCoreParamsNotSet        = errorsmod.Register(ModuleName, 1112, "Core params has not been set")
	ErrKeygenNotFound          = errorsmod.Register(ModuleName, 1113, "Err Keygen not found, Keygen block can only be updated,New keygen cannot be set")
	ErrKeygenBlockTooLow       = errorsmod.Register(ModuleName, 1114, "Please set a block number at-least 10 blocks higher than the current block number")
	ErrKeygenCompleted         = errorsmod.Register(ModuleName, 1115, "Keygen already completed")
)
