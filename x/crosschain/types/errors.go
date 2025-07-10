package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnsupportedChain        = errorsmod.Register(ModuleName, 1102, "chain parse error")
	ErrInvalidChainID          = errorsmod.Register(ModuleName, 1101, "chain id cannot be negative")
	ErrUnableToGetGasPrice     = errorsmod.Register(ModuleName, 1107, "unable to get gas price")
	ErrNotEnoughZetaBurnt      = errorsmod.Register(ModuleName, 1109, "not enough zeta burnt")
	ErrCannotFindReceiverNonce = errorsmod.Register(ModuleName, 1110, "cannot find receiver chain nonce")
	ErrGasCoinNotFound         = errorsmod.Register(ModuleName, 1113, "gas coin not found for sender chain")
	ErrUnableToParseAddress    = errorsmod.Register(ModuleName, 1115, "cannot parse address and data")
	ErrCannotProcessWithdrawal = errorsmod.Register(ModuleName, 1116, "cannot process withdrawal event")
	ErrForeignCoinNotFound     = errorsmod.Register(ModuleName, 1118, "foreign coin not found for sender chain")
	ErrCannotFindPendingNonces = errorsmod.Register(ModuleName, 1121, "cannot find pending nonces")
	ErrCannotFindTSSKeys       = errorsmod.Register(ModuleName, 1122, "cannot find TSS keys")
	ErrNonceMismatch           = errorsmod.Register(ModuleName, 1123, "nonce mismatch")
	ErrUnableToSendCoinType    = errorsmod.Register(
		ModuleName,
		1127,
		"unable to send this coin type to a receiver chain",
	)
	ErrInvalidAddress                = errorsmod.Register(ModuleName, 1128, "invalid address")
	ErrDeployContract                = errorsmod.Register(ModuleName, 1129, "unable to deploy contract")
	ErrUnableToUpdateTss             = errorsmod.Register(ModuleName, 1130, "unable to update TSS address")
	ErrNotEnoughGas                  = errorsmod.Register(ModuleName, 1131, "not enough gas")
	ErrNotEnoughFunds                = errorsmod.Register(ModuleName, 1132, "not enough funds")
	ErrProofVerificationFail         = errorsmod.Register(ModuleName, 1133, "proof verification fail")
	ErrCannotFindCctx                = errorsmod.Register(ModuleName, 1134, "cannot find cctx")
	ErrStatusNotPending              = errorsmod.Register(ModuleName, 1135, "Status not pending")
	ErrCannotFindGasParams           = errorsmod.Register(ModuleName, 1136, "cannot find gas params")
	ErrInvalidGasAmount              = errorsmod.Register(ModuleName, 1137, "invalid gas amount")
	ErrNoLiquidityPool               = errorsmod.Register(ModuleName, 1138, "no liquidity pool")
	ErrInvalidCoinType               = errorsmod.Register(ModuleName, 1139, "invalid coin type")
	ErrCannotMigrateTssFunds         = errorsmod.Register(ModuleName, 1140, "cannot migrate TSS funds")
	ErrTxBodyVerificationFail        = errorsmod.Register(ModuleName, 1141, "transaction body verification fail")
	ErrReceiverIsEmpty               = errorsmod.Register(ModuleName, 1142, "receiver is empty")
	ErrUnsupportedStatus             = errorsmod.Register(ModuleName, 1143, "unsupported status")
	ErrObservedTxAlreadyFinalized    = errorsmod.Register(ModuleName, 1144, "observed tx already finalized")
	ErrInsufficientFundsTssMigration = errorsmod.Register(ModuleName, 1145, "insufficient funds for TSS migration")
	ErrInvalidIndexValue             = errorsmod.Register(ModuleName, 1146, "invalid index hash")
	ErrInvalidStatus                 = errorsmod.Register(ModuleName, 1147, "invalid cctx status")
	ErrUnableProcessRefund           = errorsmod.Register(ModuleName, 1148, "unable to process refund")
	ErrUnableToFindZetaAccounting    = errorsmod.Register(ModuleName, 1149, "unable to find zeta accounting")
	ErrInsufficientZetaAmount        = errorsmod.Register(ModuleName, 1150, "insufficient zeta amount")
	ErrUnableToDecodeMessageString   = errorsmod.Register(ModuleName, 1151, "unable to decode message string")
	ErrInvalidRateLimiterFlags       = errorsmod.Register(ModuleName, 1152, "invalid rate limiter flags")
	ErrMaxTxOutTrackerHashesReached  = errorsmod.Register(ModuleName, 1153, "max tx out tracker hashes reached")
	ErrInitiatitingOutbound          = errorsmod.Register(ModuleName, 1154, "cannot initiate outbound")
	ErrInvalidWithdrawalAmount       = errorsmod.Register(ModuleName, 1155, "invalid withdrawal amount")
	ErrMigrationFromOldTss           = errorsmod.Register(
		ModuleName,
		1156,
		"migration tx from an old tss address detected",
	)
	ErrValidatingInbound           = errorsmod.Register(ModuleName, 1157, "unable to validate inbound")
	ErrInvalidGasLimit             = errorsmod.Register(ModuleName, 1158, "invalid gas limit")
	ErrUnableToSetOutboundInfo     = errorsmod.Register(ModuleName, 1159, "unable to set outbound info")
	ErrCCTXAlreadyFinalized        = errorsmod.Register(ModuleName, 1160, "cctx already finalized")
	ErrUnableToParseCCTXIndexBytes = errorsmod.Register(ModuleName, 1161, "unable to parse cctx index bytes")
	ErrInvalidPriorityFee          = errorsmod.Register(ModuleName, 1162, "invalid priority fee")
	ErrInvalidWithdrawalEvent      = errorsmod.Register(ModuleName, 1163, "invalid withdrawal event")
)
