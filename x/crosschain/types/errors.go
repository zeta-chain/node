package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnsupportedChain        = errorsmod.Register(ModuleName, 1102, "chain parse error")
	ErrInvalidChainID          = errorsmod.Register(ModuleName, 1101, "chain id cannot be negative")
	ErrInvalidPubKeySet        = errorsmod.Register(ModuleName, 1106, "invalid pubkeyset")
	ErrUnableToGetGasPrice     = errorsmod.Register(ModuleName, 1107, "unable to get gas price")
	ErrNotEnoughZetaBurnt      = errorsmod.Register(ModuleName, 1109, "not enough zeta burnt")
	ErrCannotFindReceiverNonce = errorsmod.Register(ModuleName, 1110, "cannot find receiver chain nonce")

	ErrGasCoinNotFound         = errorsmod.Register(ModuleName, 1113, "gas coin not found for sender chain")
	ErrUnableToParseAddress    = errorsmod.Register(ModuleName, 1115, "cannot parse address and data")
	ErrCannotProcessWithdrawal = errorsmod.Register(ModuleName, 1116, "cannot process withdrawal event")
	ErrForeignCoinNotFound     = errorsmod.Register(ModuleName, 1118, "foreign coin not found for sender chain")
	ErrNotEnoughPermissions    = errorsmod.Register(ModuleName, 1119, "not enough permissions for current actions")

	ErrCannotFindPendingNonces = errorsmod.Register(ModuleName, 1121, "cannot find pending nonces")
	ErrCannotFindTSSKeys       = errorsmod.Register(ModuleName, 1122, "cannot find TSS keys")
	ErrNonceMismatch           = errorsmod.Register(ModuleName, 1123, "nonce mismatch")
	ErrNotFoundChainParams     = errorsmod.Register(ModuleName, 1126, "not found chain chain params")
	ErrUnableToSendCoinType    = errorsmod.Register(ModuleName, 1127, "unable to send this coin type to a receiver chain")

	ErrInvalidAddress    = errorsmod.Register(ModuleName, 1128, "invalid address")
	ErrDeployContract    = errorsmod.Register(ModuleName, 1129, "unable to deploy contract")
	ErrUnableToUpdateTss = errorsmod.Register(ModuleName, 1130, "unable to update TSS address")
	ErrNotEnoughGas      = errorsmod.Register(ModuleName, 1131, "not enough gas")
	ErrNotEnoughFunds    = errorsmod.Register(ModuleName, 1132, "not enough funds")

	ErrProofVerificationFail = errorsmod.Register(ModuleName, 1133, "proof verification fail")
	ErrCannotFindCctx        = errorsmod.Register(ModuleName, 1134, "cannot find cctx")
	ErrStatusNotPending      = errorsmod.Register(ModuleName, 1135, "Status not pending")

	ErrCannotFindGasParams        = errorsmod.Register(ModuleName, 1136, "cannot find gas params")
	ErrInvalidGasAmount           = errorsmod.Register(ModuleName, 1137, "invalid gas amount")
	ErrNoLiquidityPool            = errorsmod.Register(ModuleName, 1138, "no liquidity pool")
	ErrInvalidCoinType            = errorsmod.Register(ModuleName, 1139, "invalid coin type")
	ErrCannotMigrateTssFunds      = errorsmod.Register(ModuleName, 1140, "cannot migrate TSS funds")
	ErrTxBodyVerificationFail     = errorsmod.Register(ModuleName, 1141, "transaction body verification fail")
	ErrReceiverIsEmpty            = errorsmod.Register(ModuleName, 1142, "receiver is empty")
	ErrUnsupportedStatus          = errorsmod.Register(ModuleName, 1143, "unsupported status")
	ErrObservedTxAlreadyFinalized = errorsmod.Register(ModuleName, 1144, "observed tx already finalized")

	ErrInsufficientFundsTssMigration = errorsmod.Register(ModuleName, 1145, "insufficient funds for TSS migration")
)
