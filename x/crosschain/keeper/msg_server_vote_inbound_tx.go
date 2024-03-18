package keeper

import (
	"context"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// FIXME: use more specific error types & codes

// VoteOnObservedInboundTx casts a vote on an inbound transaction observed on a connected chain. If this
// is the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized. When a ballot is finalized, a new CCTX is
// created.
//
// If the receiver chain is ZetaChain, `HandleEVMDeposit` is called. If the
// tokens being deposited are ZETA, `MintZetaToEVMAccount` is called and the
// tokens are minted to the receiver account on ZetaChain. If the tokens being
// deposited are gas tokens or ERC20 of a connected chain, ZRC20's `deposit`
// method is called and the tokens are deposited to the receiver account on
// ZetaChain. If the message is not empty, system contract's `depositAndCall`
// method is also called and an omnichain contract on ZetaChain is executed.
// Omnichain contract address and arguments are passed as part of the message.
// If everything is successful, the CCTX status is changed to `OutboundMined`.
//
// If the receiver chain is a connected chain, the `FinalizeInbound` method is
// called to prepare the CCTX to be processed as an outbound transaction. To
// cover the outbound transaction fee, the required amount of tokens submitted
// with the CCTX are swapped using a Uniswap V2 contract instance on ZetaChain
// for the ZRC20 of the gas token of the receiver chain. The ZRC20 tokens are
// then burned. The nonce is updated. If everything is successful, the CCTX
// status is changed to `PendingOutbound`.
//
// ```mermaid
// stateDiagram-v2
//
//	state evm_deposit_success <<choice>>
//	state finalize_inbound <<choice>>
//	state evm_deposit_error <<choice>>
//	PendingInbound --> evm_deposit_success: Receiver is ZetaChain
//	evm_deposit_success --> OutboundMined: EVM deposit success
//	evm_deposit_success --> evm_deposit_error: EVM deposit error
//	evm_deposit_error --> PendingRevert: Contract error
//	evm_deposit_error --> Aborted: Internal error, invalid chain, gas, nonce
//	PendingInbound --> finalize_inbound: Receiver is connected chain
//	finalize_inbound --> Aborted: Finalize inbound error
//	finalize_inbound --> PendingOutbound: Finalize inbound success
//
// ```
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	index := msg.Digest()

	// vote on inbound ballot
	// use a temporary context to not commit any ballot state change in case of error
	tmpCtx, commit := ctx.CacheContext()
	finalized, isNew, err := k.zetaObserverKeeper.VoteOnInboundBallot(
		tmpCtx,
		msg.SenderChainId,
		msg.ReceiverChain,
		msg.CoinType,
		msg.Creator,
		index,
		msg.InTxHash,
	)
	if err != nil {
		return nil, err
	}

	// If it is a new ballot, check if an inbound with the same hash, sender chain and event index has already been finalized
	// This may happen if the same inbound is observed twice where msg.Digest gives a different index
	// This check prevents double spending
	if isNew {
		if k.IsFinalizedInbound(tmpCtx, msg.InTxHash, msg.SenderChainId, msg.EventIndex) {
			return nil, cosmoserrors.Wrap(
				types.ErrObservedTxAlreadyFinalized,
				fmt.Sprintf("InTxHash:%s, SenderChainID:%d, EventIndex:%d", msg.InTxHash, msg.SenderChainId, msg.EventIndex),
			)
		}
	}
	commit()
	// If the ballot is not finalized return nil here to add vote to commit state
	if !finalized {
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	inboundCctx := k.GetInbound(ctx, msg)
	err = inboundCctx.Validate()
	if err != nil {
		return nil, err
	}
	k.ProcessInbound(ctx, &inboundCctx)
	k.SaveInbound(ctx, &inboundCctx, msg.EventIndex)
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}

// GetInbound returns a new CCTX from a given inbound message.
func (k Keeper) GetInbound(ctx sdk.Context, msg *types.MsgVoteOnObservedInboundTx) types.CrossChainTx {

	// get the latest TSS to set the TSS public key in the CCTX
	tssPub := ""
	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if tssFound {
		tssPub = tss.TssPubkey
	}
	return CreateNewCCTX(ctx, msg, msg.Digest(), tssPub, types.CctxStatus_PendingInbound, msg.SenderChainId, msg.ReceiverChain)
}

// CreateNewCCTX creates a new CCTX with the given parameters.
func CreateNewCCTX(
	ctx sdk.Context,
	msg *types.MsgVoteOnObservedInboundTx,
	index string,
	tssPubkey string,
	s types.CctxStatus,
	senderChainID,
	receiverChainID int64,
) types.CrossChainTx {
	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &types.InboundTxParams{
		Sender:                          msg.Sender,
		SenderChainId:                   senderChainID,
		TxOrigin:                        msg.TxOrigin,
		Asset:                           msg.Asset,
		Amount:                          msg.Amount,
		InboundTxObservedHash:           msg.InTxHash,
		InboundTxObservedExternalHeight: msg.InBlockHeight,
		InboundTxFinalizedZetaHeight:    0,
		InboundTxBallotIndex:            index,
		CoinType:                        msg.CoinType,
	}

	outBoundParams := &types.OutboundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChainId:                  receiverChainID,
		OutboundTxHash:                   "",
		OutboundTxTssNonce:               0,
		OutboundTxGasLimit:               msg.GasLimit,
		OutboundTxGasPrice:               "",
		OutboundTxBallotIndex:            "",
		OutboundTxObservedExternalHeight: 0,
		Amount:                           sdkmath.ZeroUint(),
		TssPubkey:                        tssPubkey,
		CoinType:                         msg.CoinType,
	}
	status := &types.Status{
		Status:              s,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		IsAbortRefunded:     false,
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaFees:         sdkmath.ZeroUint(),
		RelayedMessage:   msg.Message,
		CctxStatus:       status,
		InboundTxParams:  inboundParams,
		OutboundTxParams: []*types.OutboundTxParams{outBoundParams},
	}
	return newCctx
}

// ProcessInbound processes the inbound CCTX.
// It does a conditional dispatch to ProcessZEVMDeposit or ProcessCrosschainMsgPassing based on the receiver chain.
func (k Keeper) ProcessInbound(ctx sdk.Context, cctx *types.CrossChainTx) {
	if common.IsZetaChain(cctx.GetCurrentOutTxParam().ReceiverChainId) {
		k.ProcessZEVMDeposit(ctx, cctx)
	} else {
		k.ProcessCrosschainMsgPassing(ctx, cctx)
	}
}

// ProcessZEVMDeposit processes the EVM deposit CCTX. A deposit is a cctx which has Zetachain as the receiver chain.
// If the deposit is successful, the CCTX status is changed to OutboundMined.
// If the deposit returns an internal error i.e if HandleEVMDeposit() returns an error, but isContractReverted is false, the CCTX status is changed to Aborted.
// If the deposit is reverted, the function tries to create a revert cctx with status PendingRevert.
// If the creation of revert tx also fails it changes the status to Aborted.
// Note : Aborted CCTXs are not refunded in this function. The refund is done using a separate refunding mechanism.
// We do not return an error from this function , as all changes need to be persisted to the state.
// Instead we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to OutboundMined.
func (k Keeper) ProcessZEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx) {
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := k.HandleEVMDeposit(tmpCtx, cctx)

	if err != nil && !isContractReverted { // exceptional case; internal error; should abort CCTX
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		return
	} else if err != nil && isContractReverted { // contract call reverted; should refund
		revertMessage := err.Error()
		senderChain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.InboundTxParams.SenderChainId)
		if senderChain == nil {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "invalid sender chain")
			return
		}

		gasLimit, err := k.GetRevertGasLimit(ctx, cctx)
		if err != nil {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, fmt.Sprintf("can't get revert tx gas limit,%s", err.Error()))
			return
		}
		if gasLimit == 0 {
			// use same gas limit of outbound as a fallback -- should not happen
			gasLimit = cctx.GetCurrentOutTxParam().OutboundTxGasLimit
		}

		// create new OutboundTxParams for the revert
		revertTxParams := &types.OutboundTxParams{
			Receiver:           cctx.InboundTxParams.Sender,
			ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
			Amount:             cctx.InboundTxParams.Amount,
			OutboundTxGasLimit: gasLimit,
		}
		cctx.OutboundTxParams = append(cctx.OutboundTxParams, revertTxParams)

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
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, fmt.Sprintf("deposit revert message: %s err : %s", revertMessage, err.Error()))
			return
		}
		commitRevert()
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, revertMessage)
		return
	}
	// successful HandleEVMDeposit;
	commit()
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "Remote omnichain contract call completed")
	return
}

// ProcessCrosschainMsgPassing processes the CCTX for crosschain message passing. A crosschain message passing is a cctx which has a non-Zetachain as the receiver chain.
// If the crosschain message passing is successful, the CCTX status is changed to PendingOutbound.
// If the crosschain message passing returns an error, the CCTX status is changed to Aborted.
// We do not return an error from this function , as all changes need to be persisted to the state.
// Instead we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to PendingOutbound.
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
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		return
	}
	commit()
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingOutbound, "")
	return
}

func (k Keeper) SaveInbound(ctx sdk.Context, cctx *types.CrossChainTx, eventIndex uint64) {
	EmitEventInboundFinalized(ctx, cctx)
	k.AddFinalizedInbound(ctx,
		cctx.GetInboundTxParams().InboundTxObservedHash,
		cctx.GetInboundTxParams().SenderChainId,
		eventIndex)
	// #nosec G701 always positive
	cctx.InboundTxParams.InboundTxFinalizedZetaHeight = uint64(ctx.BlockHeight())
	cctx.InboundTxParams.TxFinalizationStatus = types.TxFinalizationStatus_Executed
	k.RemoveInTxTrackerIfExists(ctx, cctx.InboundTxParams.SenderChainId, cctx.InboundTxParams.InboundTxObservedHash)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *cctx)
}
