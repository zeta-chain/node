package keeper

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

/* ProcessSuccessfulOutbound processes a successful outbound transaction. It does the following things in one function:

	1. Change the status of the CCTX from
	 - PendingRevert to Reverted
     - PendingOutbound to OutboundMined

	2. Set the finalization status of the current outbound tx to executed

	3. Emit an event for the successful outbound transaction
*/

// This function sets CCTX status , in cases where the outbound tx is successful, but tx itself fails
// This is done because SaveSuccessfulOutbound does not set the cctx status
// For cases where the outbound tx is unsuccessful, the cctx status is automatically set to Aborted in the ProcessFailedOutbound function, so we can just return and error to trigger that

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

// This function sets CCTX status , in cases where the outbound tx is successful, but tx itself fails
// This is done because SaveSuccessfulOutbound does not set the cctx status
// For cases where the outbound tx is unsuccessful, the cctx status is automatically set to Aborted in the ProcessFailedOutbound function, so we can just return and error to trigger that
func (k Keeper) ProcessFailedOutbound(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) error {
	oldStatus := cctx.CctxStatus.Status
	if cctx.InboundTxParams.CoinType == coin.CoinType_Cmd {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.SetAbort("Outbound failed")
	} else if chains.IsZetaChain(cctx.InboundTxParams.SenderChainId) {
		// Fetch the original sender and receiver from the CCTX , since this is a revert the sender with be the receiver in the new tx
		originalSender := ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
		originalReceiver := ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver)
		data := []byte(cctx.RelayedMessage)
		indexBytes, err := cctx.GetCCTXIndexBytes()
		if err != nil {
			// Return err to save the failed outbound ad set to aborted
			return fmt.Errorf("failed reverting GetCCTXIndexBytes: %s", err.Error())
		}
		// Finalize the older outbound tx
		cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed

		// create new OutboundTxParams for the revert. We use the fixed gas limit for revert when calling zEVM
		err = cctx.AddRevertOutbound(fungiblekeeper.ZEVMGasLimitDepositAndCall.Uint64())
		if err != nil {
			// Return err to save the failed outbound ad set to aborted
			return fmt.Errorf("failed AddRevertOutbound: %s", err.Error())
		}

		// Trying to revert the transaction this would get set to a finalized status in the same block as this does not need a TSS singing
		cctx.SetPendingRevert("Outbound failed, trying revert")

		// Call evm to revert the transaction
		_, err = k.fungibleKeeper.ZEVMRevertAndCallContract(ctx,
			originalSender,
			originalReceiver,
			cctx.InboundTxParams.SenderChainId,
			cctx.GetCurrentOutTxParam().ReceiverChainId,
			cctx.GetCurrentOutTxParam().Amount.BigInt(), data, indexBytes)
		// If revert fails, we set it to abort directly there is no way to refund here as the revert failed
		if err != nil {
			return fmt.Errorf("failed ZEVMRevertAndCallContract: %s", err.Error())
		}
		cctx.SetReverted("Outbound failed, revert executed")
		cctx.GetCurrentOutTxParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
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
