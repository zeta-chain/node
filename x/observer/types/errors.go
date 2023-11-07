package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnableToAddVote = errorsmod.Register(ModuleName, 1100, "unable to add vote ")
	ErrParamsThreshold = errorsmod.Register(ModuleName, 1101, "threshold cannot be more than 1")
	ErrSupportedChains = errorsmod.Register(ModuleName, 1102, "chain not supported")
	ErrInvalidStatus   = errorsmod.Register(ModuleName, 1103, "invalid Voting Status")

	ErrObserverNotPresent      = errorsmod.Register(ModuleName, 1105, "observer for type and observation does not exist")
	ErrNotValidator            = errorsmod.Register(ModuleName, 1106, "user needs to be a validator before applying to become an observer")
	ErrValidatorStatus         = errorsmod.Register(ModuleName, 1107, "corresponding validator needs to be bonded and not jailerd")
	ErrInvalidAddress          = errorsmod.Register(ModuleName, 1108, "invalid Address")
	ErrSelfDelegation          = errorsmod.Register(ModuleName, 1109, "self Delegation for operator not found")
	ErrCheckObserverDelegation = errorsmod.Register(ModuleName, 1110, "observer delegation not sufficient")
	ErrNotAuthorizedPolicy     = errorsmod.Register(ModuleName, 1111, "msg Sender is not the authorized policy")
	ErrCoreParamsNotSet        = errorsmod.Register(ModuleName, 1112, "core params has not been set")
	ErrKeygenNotFound          = errorsmod.Register(ModuleName, 1113, "Keygen not found, Keygen block can only be updated,New keygen cannot be set")
	ErrKeygenBlockTooLow       = errorsmod.Register(ModuleName, 1114, "please set a block number at-least 10 blocks higher than the current block number")
	ErrKeygenCompleted         = errorsmod.Register(ModuleName, 1115, "keygen already completed")
	ErrNotAuthorized           = errorsmod.Register(ModuleName, 1116, "not authorized")

	ErrBlockHeaderNotFound       = errorsmod.Register(ModuleName, 1117, "block header not found")
	ErrUnrecognizedBlockHeader   = errorsmod.Register(ModuleName, 1118, "unrecognized block header")
	ErrBlockAlreadyExist         = errorsmod.Register(ModuleName, 1119, "block already exists")
	ErrNoParentHash              = errorsmod.Register(ModuleName, 1120, "no parent hash")
	ErrInvalidTimestamp          = errorsmod.Register(ModuleName, 1121, "invalid timestamp")
	ErrLastObserverCountNotFound = errorsmod.Register(ModuleName, 1122, "last observer count not found")
	ErrUpdateObserver            = errorsmod.Register(ModuleName, 1123, "unable to update observer")
	ErrNodeAccountNotFound       = errorsmod.Register(ModuleName, 1124, "node account not found")
)
