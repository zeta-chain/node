package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/emissions module sentinel errors
var (
	ErrEmissionTrackerNotFound = sdkerrors.Register(ModuleName, 1100, "Emission Tracker Not found")
	ErrParsingSenderAddress    = sdkerrors.Register(ModuleName, 1101, "Unable to parse address of sender")
	ErrAddingCoinstoTracker    = sdkerrors.Register(ModuleName, 1102, "Unable to add coins to emissionTracker ")
)
