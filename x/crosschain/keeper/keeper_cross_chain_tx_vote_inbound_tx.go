package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerKeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// FIXME: use more specific error types & codes

// Casts a vote on an inbound transaction observed on a connected chain. If this
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
	observationType := observerTypes.ObservationType_InBoundTx
	if !k.zetaObserverKeeper.IsInboundEnabled(ctx) {
		return nil, types.ErrNotEnoughPermissions
	}
	// GetChainFromChainID makes sure we are getting only supported chains , if a chain support has been turned on using gov proposal, this function returns nil
	observationChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.SenderChainId)
	if observationChain == nil {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.SenderChainId, observationType.String()))
	}
	receiverChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ReceiverChain)
	if receiverChain == nil {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.ReceiverChain, observationType.String()))
	}
	tssPub := ""
	tss, tssFound := k.GetTSS(ctx)
	if tssFound {
		tssPub = tss.TssPubkey
	}
	// IsAuthorized does various checks against the list of observer mappers
	ok, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, observationChain)
	if !ok {
		return nil, err
	}

	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	ballot, isNew, err := k.zetaObserverKeeper.FindBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	if isNew {
		observerKeeper.EmitEventBallotCreated(ctx, ballot, msg.InTxHash, observationChain.String())
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.zetaObserverKeeper.AddVoteToBallot(ctx, ballot, msg.Creator, observerTypes.VoteType_SuccessObservation)
	if err != nil {
		return nil, err
	}

	_, isFinalized := k.zetaObserverKeeper.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	// Validation if we want to send ZETA to external chain, but there is no ZETA token.
	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, receiverChain.ChainId)
	if !found {
		return nil, types.ErrNotFoundCoreParams
	}
	if receiverChain.IsExternalChain() && coreParams.ZetaTokenContractAddress == "" && msg.CoinType == common.CoinType_Zeta {
		return nil, types.ErrUnableToSendCoinType
	}

	// ******************************************************************************
	// below only happens when ballot is finalized: exactly when threshold vote is in
	// ******************************************************************************

	// Inbound Ballot has been finalized , Create CCTX
	cctx := k.CreateNewCCTX(ctx, msg, index, tssPub, types.CctxStatus_PendingInbound, observationChain, receiverChain)
	defer func() {
		EmitEventInboundFinalized(ctx, &cctx)
		cctx.InboundTxParams.InboundTxFinalizedZetaHeight = uint64(ctx.BlockHeight())
		k.RemoveInTxTrackerIfExists(ctx, cctx.InboundTxParams.SenderChainId, cctx.InboundTxParams.InboundTxObservedHash)
		k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	}()
	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail
	if receiverChain.IsZetaChain() {
		tmpCtx, commit := ctx.CacheContext()
		isContractReverted, err := k.HandleEVMDeposit(tmpCtx, &cctx, *msg, observationChain)
		if err != nil && !isContractReverted { // exceptional case; internal error; should abort CCTX
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		} else if err != nil && isContractReverted { // contract call reverted; should refund
			revertMessage := err.Error()
			chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(cctx.InboundTxParams.SenderChainId)
			if chain == nil {
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "invalid sender chain")
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
			if !isFound {
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "cannot find gas price")
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			// create new OutboundTxParams for the revert
			cctx.OutboundTxParams = append(cctx.OutboundTxParams, &types.OutboundTxParams{
				Receiver:           cctx.InboundTxParams.Sender,
				ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
				Amount:             cctx.InboundTxParams.Amount,
				CoinType:           cctx.InboundTxParams.CoinType,
				OutboundTxGasLimit: 0, // for fungible refund, use default gasLimit
				OutboundTxGasPrice: medianGasPrice.MulUint64(2).String(),
			})
			if err = k.UpdateNonce(ctx, chain.ChainId, &cctx); err != nil {
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			// do not commit() here;
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, revertMessage)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		} else { // successful HandleEVMDeposit;
			commit()
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "Remote omnichain contract call completed")
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
	} else { // Cross Chain SWAP
		tmpCtx, commit := ctx.CacheContext()
		err = func() error {
			err := k.PayGasInZetaAndUpdateCctx(tmpCtx, receiverChain.ChainId, &cctx, false)
			if err != nil {
				return err
			}
			err = k.UpdateNonce(tmpCtx, receiverChain.ChainId, &cctx)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			// do not commit anything here as the CCTX should be aborted
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
		commit()
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingOutbound, "")
		k.RemoveInTxTrackerIfExists(ctx, cctx.InboundTxParams.SenderChainId, cctx.InboundTxParams.InboundTxObservedHash)
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}
}
