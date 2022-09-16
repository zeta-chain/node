package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/fungible module sentinel errors
var (
	ErrSample            = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrABIPack           = sdkerrors.Register(ModuleName, 1101, "abi pack error")
	ErrABIGet            = sdkerrors.Register(ModuleName, 1102, "abi get error")
	ErrUnexpectedEvent   = sdkerrors.Register(ModuleName, 1103, "unexpected event")
	ErrABIUnpack         = sdkerrors.Register(ModuleName, 1104, "abi unpack error")
	ErrBlanceQuery       = sdkerrors.Register(ModuleName, 1105, "balance query error")
	ErrBalanceInvariance = sdkerrors.Register(ModuleName, 1106, "balance invariance error")
)
