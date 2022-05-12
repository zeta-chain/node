package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/zetacore module sentinel errors
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
	// this line is used by starport scaffolding # ibc/errors

	ErrFloatParseError    = sdkerrors.Register(ModuleName, 1101, "float parse error")
	ErrUnsupportedChain   = sdkerrors.Register(ModuleName, 1102, "chain parse error")
	ErrDuplicateMsg       = sdkerrors.Register(ModuleName, 1103, "duplicate msg error")
	ErrNotBondedValidator = sdkerrors.Register(ModuleName, 1104, "not a bonded validator error")
)
