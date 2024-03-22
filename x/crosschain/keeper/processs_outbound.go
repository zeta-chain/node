package keeper

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

/* ProcessSuccessfulOutbound processes a successful outbound transaction. It does the following things in one function:

	1. Change the status of the CCTX from
	 - PendingRevert to Reverted
     - PendingOutbound to OutboundMined

	2. Set the finalization status of the current outbound tx to executed

	3. Emit an event for the successful outbound transaction
*/

func (k Keeper) ProcessSuccessfulOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) {
	oldStatus := cctx.CctxStatus.Status
	switch oldStatus {
	case types.CctxStatus_PendingRevert:
		cctx.SetReverted("Outbound succeeded, revert executed")
	case types.CctxStatus_PendingOutbound:
		cctx.SetOutBoundMined("Outbound succeeded, mined")
	default:
		return
	}
	cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundSuccess(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
}

/*
ProcessFailedOutbound processes a failed outbound transaction. It does the following things in one function:

 1. For Admin Tx or a withdrawal from Zeta chain, it aborts the CCTX

 2. For other CCTX
    - If the CCTX is in PendingOutbound, it creates a revert tx and sets the finalization status of the current outbound tx to executed
    - If the CCTX is in PendingRevert, it sets the Status to Aborted

 3. Emit an event for the failed outbound transaction

 4. Set the finalization status of the current outbound tx to executed. If a revert tx is is created, the finalization status is not set, it would get set when the revert is processed via a subsequent transaction
*/
func (k Keeper) ProcessFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) error {
	oldStatus := cctx.CctxStatus.Status
	if cctx.InboundTxParams.CoinType == common.CoinType_Cmd || common.IsZetaChain(cctx.InboundTxParams.SenderChainId) {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "")
	} else {
		switch oldStatus {
		case types.CctxStatus_PendingOutbound:

			gasLimit, err := k.GetRevertGasLimit(ctx, *cctx)
			if err != nil {
				return cosmoserrors.Wrap(err, "GetRevertGasLimit")
			}
			if gasLimit == 0 {
				// use same gas limit of outbound as a fallback -- should not happen
				gasLimit = cctx.OutboundTxParams[0].OutboundTxGasLimit
			}

			// create new OutboundTxParams for the revert
			err = cctx.AddRevertOutbound(gasLimit)
			if err != nil {
				return cosmoserrors.Wrap(err, "AddRevertOutbound")
			}

			err = k.PayGasAndUpdateCctx(
				ctx,
				cctx.InboundTxParams.SenderChainId,
				cctx,
				cctx.OutboundTxParams[0].Amount,
				false,
			)
			if err != nil {
				return err
			}
			err = k.UpdateNonce(ctx, cctx.InboundTxParams.SenderChainId, cctx)
			if err != nil {
				return err
			}
			// Not setting the finalization status here, the required changes have been made while creating the revert tx
			cctx.SetPendingRevert("Outbound failed, start revert")
		case types.CctxStatus_PendingRevert:
			cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
			cctx.SetAbort("Outbound failed: revert failed; abort TX")
		}
	}
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundFailure(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
	return nil
}

// ProcessOutbound processes the finalization of an outbound transaction based on the ballot status
// The state is committed only if the individual steps are successful
func (k Keeper) ProcessOutbound(ctx sdk.Context, cctx *types.CrossChainTx, ballotStatus observertypes.BallotStatus, valueReceived string) error {
	tmpCtx, commit := ctx.CacheContext()
	err := func() error {
		switch ballotStatus {
		case observertypes.BallotStatus_BallotFinalized_SuccessObservation:
			k.ProcessSuccessfulOutbound(tmpCtx, cctx, valueReceived)
		case observertypes.BallotStatus_BallotFinalized_FailureObservation:
			err := k.ProcessFailedOutbound(tmpCtx, cctx, valueReceived)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}
	err = cctx.Validate()
	if err != nil {
		return err
	}
	commit()
	return nil
}
