package keeper

import (
	"context"
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
	finalized, err := k.zetaObserverKeeper.VoteOnInboundBallot(
		ctx,
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
	if !finalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	// get the latest TSS to set the TSS public key in the CCTX
	tssPub := ""
	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if tssFound {
		tssPub = tss.TssPubkey
	}

	// create the CCTX
	cctx := k.CreateNewCCTX(
		ctx,
		msg,
		index,
		tssPub,
		types.CctxStatus_PendingInbound,
		msg.SenderChainId,
		msg.ReceiverChain,
	)

	defer func() {
		EmitEventInboundFinalized(ctx, &cctx)
		// #nosec G701 always positive
		cctx.InboundTxParams.InboundTxFinalizedZetaHeight = uint64(ctx.BlockHeight())
		k.RemoveInTxTrackerIfExists(ctx, cctx.InboundTxParams.SenderChainId, cctx.InboundTxParams.InboundTxObservedHash)
		k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	}()

	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail
	if common.IsZetaChain(msg.ReceiverChain) {
		tmpCtx, commit := ctx.CacheContext()
		isContractReverted, err := k.HandleEVMDeposit(tmpCtx, &cctx, *msg, msg.SenderChainId)

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

			gasLimit, err := k.GetRevertGasLimit(ctx, cctx)
			if err != nil {
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "can't get revert tx gas limit"+err.Error())
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			if gasLimit == 0 {
				// use same gas limit of outbound as a fallback -- should not happen
				gasLimit = msg.GasLimit
			}

			// create new OutboundTxParams for the revert
			revertTxParams := &types.OutboundTxParams{
				Receiver:           cctx.InboundTxParams.Sender,
				ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
				Amount:             cctx.InboundTxParams.Amount,
				CoinType:           cctx.InboundTxParams.CoinType,
				OutboundTxGasLimit: gasLimit,
			}
			cctx.OutboundTxParams = append(cctx.OutboundTxParams, revertTxParams)

			// we create a new cached context, and we don't commit the previous one with EVM deposit
			tmpCtx, commit := ctx.CacheContext()
			err = func() error {
				err := k.PayGasAndUpdateCctx(
					tmpCtx,
					chain.ChainId,
					&cctx,
					cctx.InboundTxParams.Amount,
					false,
				)
				if err != nil {
					return err
				}
				return k.UpdateNonce(tmpCtx, chain.ChainId, &cctx)
			}()
			if err != nil {
				// do not commit anything here as the CCTX should be aborted

				// gas payment for erc20 type might fail because no liquidity pool is defined to swap the zrc20 token into the gas token
				// in this gas we should refund the sender on ZetaChain
				if cctx.InboundTxParams.CoinType == common.CoinType_ERC20 {

					if err := k.RefundAmountOnZetaChain(ctx, cctx, cctx.InboundTxParams.Amount); err != nil {
						// log the error
						k.Logger(ctx).Error("failed to refund amount of aborted cctx on ZetaChain",
							"error", err,
							"sender", cctx.InboundTxParams.Sender,
							"amount", cctx.InboundTxParams.Amount.String(),
						)
					}
				}

				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error()+" deposit revert message: "+revertMessage)
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			commit()
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, revertMessage)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil

		}
		// successful HandleEVMDeposit;
		commit()
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "Remote omnichain contract call completed")
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	// Receiver is not ZetaChain: Cross Chain SWAP
	tmpCtx, commit := ctx.CacheContext()
	err = func() error {
		err := k.PayGasAndUpdateCctx(
			tmpCtx,
			msg.ReceiverChain,
			&cctx,
			cctx.InboundTxParams.Amount,
			false,
		)
		if err != nil {
			return err
		}
		return k.UpdateNonce(tmpCtx, msg.ReceiverChain, &cctx)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}
	commit()
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingOutbound, "")
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}
