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
)
