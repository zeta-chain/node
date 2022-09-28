package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/zetaobserver module sentinel errors
var (
	ErrVoter = sdkerrors.Register(ModuleName, 1100, "Error Voter")
)
