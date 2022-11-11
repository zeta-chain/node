package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrFloatParseError             = sdkerrors.Register(ModuleName, 1101, "float parse error")
	ErrUnsupportedChain            = sdkerrors.Register(ModuleName, 1102, "chain parse error")
	ErrDuplicateMsg                = sdkerrors.Register(ModuleName, 1103, "duplicate msg error")
	ErrNotBondedValidator          = sdkerrors.Register(ModuleName, 1104, "not a bonded validator error")
	ErrOutOfBound                  = sdkerrors.Register(ModuleName, 1105, "out of bound of array")
	ErrInvalidPubKeySet            = sdkerrors.Register(ModuleName, 1106, "invalid pubkeyset")
	ErrUnableToGetGasPrice         = sdkerrors.Register(ModuleName, 1107, "unable to get gas price")
	ErrUnableToGetConversionRate   = sdkerrors.Register(ModuleName, 1108, "zeta conversion rate not found")
	ErrNotEnoughZetaBurnt          = sdkerrors.Register(ModuleName, 1109, "not enough zeta burnt")
	ErrCannotFindReceiverNonce     = sdkerrors.Register(ModuleName, 1110, "cannot find receiver chain nonce")
	ErrStatusTransitionNotPossible = sdkerrors.Register(ModuleName, 1111, "cannot transition status for CCTX")
	ErrNotAuthorized               = sdkerrors.Register(ModuleName, 1112, "Err not authorized")
	ErrGasCoinNotFound             = sdkerrors.Register(ModuleName, 1113, "Err gas coin not found for SenderChain")
)
