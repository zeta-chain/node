package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrFloatParseError    = sdkerrors.Register(ModuleName, 1101, "float parse error")
	ErrUnsupportedChain   = sdkerrors.Register(ModuleName, 1102, "chain parse error")
	ErrDuplicateMsg       = sdkerrors.Register(ModuleName, 1103, "duplicate msg error")
	ErrNotBondedValidator = sdkerrors.Register(ModuleName, 1104, "not a bonded validator error")
	ErrOutOfBound         = sdkerrors.Register(ModuleName, 1105, "out of bound of array")
	ErrInvalidPubKeySet   = sdkerrors.Register(ModuleName, 1106, "invalid pubkeyset")
)
