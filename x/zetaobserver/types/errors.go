package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/zetaobserver module sentinel errors
var (
	ErrUnableToAddVote = sdkerrors.Register(ModuleName, 1100, "Unable to add vote ")
	ErrParamsThreshold = sdkerrors.Register(ModuleName, 1101, "Threshold cannot be more than 1")
)
