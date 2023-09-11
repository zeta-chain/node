package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnsupportedChain         = errorsmod.Register(ModuleName, 1102, "chain parse error")
	ErrInvalidChainID           = errorsmod.Register(ModuleName, 1101, "chain id cannot be negative")
	ErrInvalidPubKeySet         = errorsmod.Register(ModuleName, 1106, "invalid pubkeyset")
	ErrUnableToGetGasPrice      = errorsmod.Register(ModuleName, 1107, "unable to get gas price")
	ErrNotEnoughZetaBurnt       = errorsmod.Register(ModuleName, 1109, "not enough zeta burnt")
	ErrCannotFindReceiverNonce  = errorsmod.Register(ModuleName, 1110, "cannot find receiver chain nonce")
	ErrCannotFindPendingTxQueue = errorsmod.Register(ModuleName, 1111, "cannot find pending tx queue")

	ErrGasCoinNotFound         = errorsmod.Register(ModuleName, 1113, "Err gas coin not found for SenderChain")
	ErrUnableToDepositZRC20    = errorsmod.Register(ModuleName, 1114, "Unable to deposit ZRC20 ")
	ErrUnableToParseContract   = errorsmod.Register(ModuleName, 1115, "Cannot parse contract and data")
	ErrCannotProcessWithdrawal = errorsmod.Register(ModuleName, 1116, "Cannot process withdrawal event")
	ErrForeignCoinNotFound     = errorsmod.Register(ModuleName, 1118, "Err gas coin not found for SenderChain")
	ErrNotEnoughPermissions    = errorsmod.Register(ModuleName, 1119, "Not enough permissions for current actions")

	ErrCannotFindPendingNonces = errorsmod.Register(ModuleName, 1121, "Err Cannot find pending nonces")
	ErrCannotFindTSSKeys       = errorsmod.Register(ModuleName, 1122, "Err Cannot find TSS keys")
	ErrNonceMismatch           = errorsmod.Register(ModuleName, 1123, "Err Nonce mismatch")
	ErrNotFoundCoreParams      = errorsmod.Register(ModuleName, 1126, "Not found chain core params")
	ErrUnableToSendCoinType    = errorsmod.Register(ModuleName, 1127, "Unable to send this coin type to a receiver chain")

	ErrInvalidAddress    = errorsmod.Register(ModuleName, 1128, "Invalid address")
	ErrDeployContract    = errorsmod.Register(ModuleName, 1129, "Unable to deploy contract")
	ErrUnableToUpdateTss = errorsmod.Register(ModuleName, 1130, "Unable to update TSS address")
	ErrNotEnoughGas      = errorsmod.Register(ModuleName, 1131, "not enough gas")
	ErrNotEnoughFunds    = errorsmod.Register(ModuleName, 1132, "Not enough funds")
)
