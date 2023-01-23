package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/emissions module sentinel errors
var (
	ErrEmissionTrackerNotFound = sdkerrors.Register(ModuleName, 1100, "Emission Tracker Not found")
)
