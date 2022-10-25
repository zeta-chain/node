package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnableToAddVote = sdkerrors.Register(ModuleName, 1100, "Unable to add vote ")
	ErrParamsThreshold = sdkerrors.Register(ModuleName, 1101, "Threshold cannot be more than 1")
	ErrSupportedChains = sdkerrors.Register(ModuleName, 1102, "Err supported Chains")
	ErrInvalidStatus   = sdkerrors.Register(ModuleName, 1103, "Invalid Voting Status")
	ErrNotValidator    = sdkerrors.Register(ModuleName, 1104, "User needs to be a validator before applying to become an observer")
	ErrInvalidAddress  = sdkerrors.Register(ModuleName, 1105, "Invalid Address")
)
