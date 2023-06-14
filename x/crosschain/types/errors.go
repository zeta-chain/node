package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnsupportedChain         = sdkerrors.Register(ModuleName, 1102, "chain parse error")
	ErrInvalidChainID           = sdkerrors.Register(ModuleName, 1101, "chain id cannot be negative")
	ErrInvalidPubKeySet         = sdkerrors.Register(ModuleName, 1106, "invalid pubkeyset")
	ErrUnableToGetGasPrice      = sdkerrors.Register(ModuleName, 1107, "unable to get gas price")
	ErrNotEnoughZetaBurnt       = sdkerrors.Register(ModuleName, 1109, "not enough zeta burnt")
	ErrCannotFindReceiverNonce  = sdkerrors.Register(ModuleName, 1110, "cannot find receiver chain nonce")
	ErrCannotFindPendingTxQueue = sdkerrors.Register(ModuleName, 1111, "cannot find pending tx queue")
	ErrNotAuthorized            = sdkerrors.Register(ModuleName, 1112, "Err not authorized")
	ErrGasCoinNotFound          = sdkerrors.Register(ModuleName, 1113, "Err gas coin not found for SenderChain")
	ErrUnableToDepositZRC20     = sdkerrors.Register(ModuleName, 1114, "Unable to deposit ZRC20 ")
	ErrUnableToParseContract    = sdkerrors.Register(ModuleName, 1115, "Cannot parse contract and data")
	ErrCannotProcessWithdrawal  = sdkerrors.Register(ModuleName, 1116, "Cannot process withdrawal event")
	ErrForeignCoinNotFound      = sdkerrors.Register(ModuleName, 1118, "Err gas coin not found for SenderChain")
	ErrNotEnoughPermissions     = sdkerrors.Register(ModuleName, 1119, "Not enough permissions for current actions")
	ErrKeygenNotFound           = sdkerrors.Register(ModuleName, 1120, "Err Keygen not found, Keygen block can only be updated,New keygen cannot be set")
	ErrSuccessfullCompleted     = sdkerrors.Register(ModuleName, 1120, "Err Keygen not found, Keygen block can only be updated,New keygen cannot be set")
	ErrCannotFindPendingNonces  = sdkerrors.Register(ModuleName, 1121, "Err Cannot find pending nonces")
	ErrCannotFindTSSKeys        = sdkerrors.Register(ModuleName, 1122, "Err Cannot find TSS keys")
	ErrNonceMismatch            = sdkerrors.Register(ModuleName, 1123, "Err Nonce mismatch")
	ErrKeygenBlockTooLow        = sdkerrors.Register(ModuleName, 1125, "Please set a block number at-least 10 blocks higher than the current block number")
	ErrNotFoundCoreParams       = sdkerrors.Register(ModuleName, 1126, "Not found chain core params")
	ErrUnableToSendCoinType     = sdkerrors.Register(ModuleName, 1127, "Unable to send this coin type to a receiver chain")
)
