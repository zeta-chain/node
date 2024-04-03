package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// ProcessInbound processes the inbound CCTX.
// It does a conditional dispatch to ProcessZEVMDeposit or ProcessCrosschainMsgPassing based on the receiver chain.
// The internal functions handle the state changes and error handling.
func (k Keeper) ProcessInbound(ctx sdk.Context, cctx *types.CrossChainTx) {
	if chains.IsZetaChain(cctx.GetCurrentOutTxParam().ReceiverChainId) {
		k.ProcessZEVMDeposit(ctx, cctx)
	} else {
		k.ProcessCrosschainMsgPassing(ctx, cctx)
	}
}

/*
ProcessZEVMDeposit processes the EVM deposit CCTX. A deposit is a cctx which has Zetachain as the receiver chain.It trasnsitions state according to the following rules:
  - If the deposit is successful, the CCTX status is changed to OutboundMined.
  - If the deposit returns an internal error i.e if HandleEVMDeposit() returns an error, but isContractReverted is false, the CCTX status is changed to Aborted.
  - If the deposit is reverted, the function tries to create a revert cctx with status PendingRevert.
  - If the creation of revert tx also fails it changes the status to Aborted.

Note : Aborted CCTXs are not refunded in this function. The refund is done using a separate refunding mechanism.
We do not return an error from this function , as all changes need to be persisted to the state.
Instead we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to OutboundMined.
*/
func (k Keeper) ProcessZEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx) {
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := k.HandleEVMDeposit(tmpCtx, cctx)

	if err != nil && !isContractReverted { // exceptional case; internal error; should abort CCTX
		cctx.SetAbort(err.Error())
		return
	} else if err != nil && isContractReverted { // contract call reverted; should refund
		revertMessage := err.Error()
		senderChain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.InboundTxParams.SenderChainId)
		if senderChain == nil {
			cctx.SetAbort(fmt.Sprintf("invalid sender chain id %d", cctx.InboundTxParams.SenderChainId))
			return
		}

		gasLimit, err := k.GetRevertGasLimit(ctx, *cctx)
		if err != nil {
			cctx.SetAbort(fmt.Sprintf("revert gas limit error: %s", err.Error()))
			return
		}
		if gasLimit == 0 {
			// use same gas limit of outbound as a fallback -- should not happen
			gasLimit = cctx.GetCurrentOutTxParam().OutboundTxGasLimit
		}

		err = cctx.AddRevertOutbound(gasLimit)
		if err != nil {
			cctx.SetAbort(fmt.Sprintf("revert outbound error: %s", err.Error()))
			return
		}

		// we create a new cached context, and we don't commit the previous one with EVM deposit
		tmpCtxRevert, commitRevert := ctx.CacheContext()
		err = func() error {
			err := k.PayGasAndUpdateCctx(
				tmpCtxRevert,
				senderChain.ChainId,
				cctx,
				cctx.InboundTxParams.Amount,
				false,
			)
			if err != nil {
				return err
			}
			// Update nonce using senderchain id as this is a revert tx and would go back to the original sender
			return k.UpdateNonce(tmpCtxRevert, senderChain.ChainId, cctx)
		}()
		if err != nil {
			cctx.SetAbort(fmt.Sprintf("deposit revert message: %s err : %s", revertMessage, err.Error()))
			return
		}
		commitRevert()
		cctx.SetPendingRevert(revertMessage)
		return
	}
	// successful HandleEVMDeposit;
	commit()
	cctx.SetOutBoundMined("Remote omnichain contract call completed")
	return
}

/*
ProcessCrosschainMsgPassing processes the CCTX for crosschain message passing. A crosschain message passing is a cctx which has a non-Zetachain as the receiver chain.It trasnsitions state according to the following rules:
  - If the crosschain message passing is successful, the CCTX status is changed to PendingOutbound.
  - If the crosschain message passing returns an error, the CCTX status is changed to Aborted.
    We do not return an error from this function, as all changes need to be persisted to the state.
    Instead, we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to PendingOutbound.
*/
func (k Keeper) ProcessCrosschainMsgPassing(ctx sdk.Context, cctx *types.CrossChainTx) {
	tmpCtx, commit := ctx.CacheContext()
	outboundReceiverChainID := cctx.GetCurrentOutTxParam().ReceiverChainId
	err := func() error {
		err := k.PayGasAndUpdateCctx(
			tmpCtx,
			outboundReceiverChainID,
			cctx,
			cctx.InboundTxParams.Amount,
			false,
		)
		if err != nil {
			return err
		}
		return k.UpdateNonce(tmpCtx, outboundReceiverChainID, cctx)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted
		cctx.SetAbort(err.Error())
		return
	}
	commit()
	cctx.SetPendingOutbound("")
	return
}
