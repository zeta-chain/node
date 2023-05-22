package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnableToAddVote = sdkerrors.Register(ModuleName, 1100, "Unable to add vote ")
	ErrParamsThreshold = sdkerrors.Register(ModuleName, 1101, "Threshold cannot be more than 1")
	ErrSupportedChains = sdkerrors.Register(ModuleName, 1102, "Err chain not supported")
	ErrInvalidStatus   = sdkerrors.Register(ModuleName, 1103, "Invalid Voting Status")

	ErrObserverNotPresent      = sdkerrors.Register(ModuleName, 1105, "Observer for type and observation does not exist")
	ErrNotValidator            = sdkerrors.Register(ModuleName, 1106, "User needs to be a validator before applying to become an observer")
	ErrValidatorStatus         = sdkerrors.Register(ModuleName, 1107, "Corresponding validator needs to be bonded and not jailerd")
	ErrInvalidAddress          = sdkerrors.Register(ModuleName, 1108, "Invalid Address")
	ErrSelfDelegation          = sdkerrors.Register(ModuleName, 1109, "Self Delegation for operator not found")
	ErrCheckObserverDelegation = sdkerrors.Register(ModuleName, 1110, "Observer delegation not sufficient")
	ErrNotAuthorizedPolicy     = sdkerrors.Register(ModuleName, 1111, "Msg Sender is not the authorized policy")
	ErrCoreParamsNotSet        = sdkerrors.Register(ModuleName, 1112, "Core params has not been set")
)
