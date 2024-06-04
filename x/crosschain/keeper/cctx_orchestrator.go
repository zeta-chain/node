package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) ValidateOutboundZEVM(ctx sdk.Context, cctx *types.CrossChainTx, zevmError error, isContractReverted bool) (newCCTXStatus types.CctxStatus) {
	if zevmError != nil && !isContractReverted {
		// exceptional case; internal error; should abort CCTX
		cctx.SetAbort(zevmError.Error())
		return types.CctxStatus_Aborted
	}

	if zevmError != nil && isContractReverted {
		// contract call reverted; should refund via a revert tx
		err := k.tryRevertOutbound(ctx, cctx, zevmError)
		if err != nil {
			cctx.SetAbort(err.Error())
			return types.CctxStatus_Aborted
		}

		return types.CctxStatus_PendingRevert
	}
	_, commit := ctx.CacheContext()
	commit()
	cctx.SetOutBoundMined("Remote omnichain contract call completed")
	return types.CctxStatus_OutboundMined
}

func (k Keeper) tryRevertOutbound(ctx sdk.Context, cctx *types.CrossChainTx, zevmError error) error {
	senderChain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.InboundParams.SenderChainId)
	if senderChain == nil {
		return fmt.Errorf("invalid sender chain id %d", cctx.InboundParams.SenderChainId)
	}
	// use same gas limit of outbound as a fallback -- should not be required
	gasLimit, err := k.GetRevertGasLimit(ctx, *cctx)
	if err != nil {
		return fmt.Errorf("revert gas limit error: %s", err.Error())
	}
	if gasLimit == 0 {
		gasLimit = cctx.GetCurrentOutboundParam().GasLimit
	}

	revertMessage := zevmError.Error()
	err = cctx.AddRevertOutbound(gasLimit)
	if err != nil {
		return fmt.Errorf("revert outbound error: %s", err.Error())
	}

	// we create a new cached context, and we don't commit the previous one with EVM deposit
	tmpCtxRevert, commitRevert := ctx.CacheContext()
	err = func() error {
		err := k.PayGasAndUpdateCctx(
			tmpCtxRevert,
			senderChain.ChainId,
			cctx,
			cctx.InboundParams.Amount,
			false,
		)
		if err != nil {
			return err
		}

		// update nonce using senderchain id as this is a revert tx and would go back to the original sender
		return k.UpdateNonce(tmpCtxRevert, senderChain.ChainId, cctx)
	}()
	if err != nil {
		return fmt.Errorf("deposit revert message: %s err : %s", revertMessage, err.Error())
	}
	commitRevert()
	cctx.SetPendingRevert(revertMessage)
	return nil
}

func (k Keeper) ValidateOutboundObservers(ctx sdk.Context, cctx *types.CrossChainTx, ballotStatus observertypes.BallotStatus, valueReceived string) error {
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
