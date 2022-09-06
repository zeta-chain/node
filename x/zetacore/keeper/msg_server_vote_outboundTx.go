package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"strconv"
)

func (k msgServer) VoteOnObservedOutboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedOutboundTx) (*types.MsgVoteOnObservedOutboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	log.Info().Msgf("ReceiveConfirmation: %s", msg.String())

	if !k.isAuthorized(ctx, msg.Creator) {
		return nil, sdkerrors.Wrap(types.ErrNotBondedValidator, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	CctxIndex := msg.CctxHash
	cctx, isFound := k.GetCrossChainTx(ctx, CctxIndex)
	if !isFound {
		log.Error().Msgf("Cannot find Incoming Broadcast broadcast tx hash %s on %s chain", CctxIndex, msg.OutTxChain)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("Cannot find broadcast tx hash %s on %s chain", CctxIndex, msg.OutTxChain))
	}

	if msg.Status != common.ReceiveStatus_Failed {
		if msg.MMint != cctx.ZetaMint {
			log.Error().Msgf("ReceiveConfirmation: Mint mismatch: %s vs %s", msg.MMint, cctx.ZetaMint)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("MMint %s does not match send ZetaMint %s", msg.MMint, cctx.ZetaMint))
		}
	}

	receiveIndex := msg.Digest()
	receive, isFound := k.GetReceive(ctx, receiveIndex)

	if isFound {
		if isDuplicateSigner(msg.Creator, receive.Signers) {
			log.Info().Msgf("ReceiveConfirmation: TX %s has already been signed by %s", receiveIndex, msg.Creator)
			return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("ReceiveConfirmation: TX %s has already been signed by %s", receiveIndex, msg.Creator))
		}
		receive.Signers = append(receive.Signers, msg.Creator)
	} else {
		if !k.IsChainSupported(ctx, msg.OutTxChain) {
			return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, "Receiving chain is not supported")
		}
		receive = types.Receive{
			Creator:             "",
			Index:               receiveIndex,
			SendHash:            CctxIndex,
			OutTxHash:           msg.ObservedOutTxHash,
			OutBlockHeight:      msg.ObservedOutTxBlockHeight,
			FinalizedMetaHeight: 0,
			Signers:             []string{msg.Creator},
			Status:              msg.Status, // TODO : Check how this status is set
			Chain:               msg.OutTxChain,
		}
	}

	hasEnoughVotes := k.hasSuperMajorityValidators(ctx, receive.Signers)
	if hasEnoughVotes {
		// Finalize Receive Struct
		err := FinalizeReceive(k, ctx, &cctx, msg, &receive)
		if err != nil {
			return nil, err
		}

		// Remove OutTX tracker
		if receive.Status == common.ReceiveStatus_Success || receive.Status == common.ReceiveStatus_Failed {
			index := fmt.Sprintf("%s/%s", msg.OutTxChain, strconv.Itoa(int(msg.OutTxTssNonce)))
			k.RemoveOutTxTracker(ctx, index)
		}

		cctx.OutBoundTxParams.OutBoundTXReceiveIndex = receive.Index
		cctx.OutBoundTxParams.OutBoundTxHash = receive.OutTxHash
		cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
		k.SetCrossChainTx(ctx, cctx)
	}
	k.SetReceive(ctx, receive)
	return &types.MsgVoteOnObservedOutboundTxResponse{}, nil
}

func HandleFeeBalances(k msgServer, ctx sdk.Context, balanceAmount sdk.Uint) error {
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(common.ZETADenom, sdk.NewIntFromBigInt(balanceAmount.BigInt()))))
	if err != nil {
		log.Error().Msgf("ReceiveConfirmation: failed to mint coins: %s", err.Error())
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to mint coins: %s", err.Error()))
	}
	return nil
}

func FinalizeReceive(k msgServer, ctx sdk.Context, cctx *types.CrossChainTx, msg *types.MsgVoteOnObservedOutboundTx, receive *types.Receive) error {
	receive.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
	zetaBurnt := cctx.ZetaBurnt
	zetaMinted := cctx.ZetaMint

	switch receive.Status {
	case common.ReceiveStatus_Success:
		oldStatus := cctx.CctxStatus.Status.String()
		if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Reverted)
		}
		if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined)
		}
		newStatus := cctx.CctxStatus.Status.String()
		if zetaBurnt.LT(zetaMinted) {
			// TODO :Handle Error ?
		}
		balanceAmount := zetaBurnt.Sub(zetaMinted)
		err := HandleFeeBalances(k, ctx, balanceAmount)
		if err != nil {
			return err
		}
		EmitSuccess(ctx, msg, receive, oldStatus, newStatus)
	case common.ReceiveStatus_Failed:
		oldStatus := cctx.CctxStatus.Status.String()
		if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
			chain := cctx.InBoundTxParams.SenderChain
			err := k.updateCctx(ctx, chain, cctx)
			if err != nil {
				return err
			}
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert)
		} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted)
		}
		newstatus := cctx.CctxStatus.Status.String()
		EmitFailure(ctx, msg, receive, oldStatus, newstatus)
	}
	return nil
}

func EmitSuccess(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, receive *types.Receive, oldStatus, newStatus string) {
	event := sdk.NewEvent(types.OutboundTxSuccessful,
		sdk.NewAttribute(types.CctxIndex, receive.SendHash),
		sdk.NewAttribute(types.OutTxHash, receive.OutTxHash),
		sdk.NewAttribute(types.ZetaMint, msg.MMint.String()),
		sdk.NewAttribute(types.OutBoundChain, msg.OutTxChain),
		sdk.NewAttribute(types.OldStatus, oldStatus),
		sdk.NewAttribute(types.NewStatus, newStatus),
	)
	ctx.EventManager().EmitEvent(event)
}

func EmitFailure(ctx sdk.Context, msg *types.MsgVoteOnObservedOutboundTx, receive *types.Receive, oldStatus, newStatus string) {
	event := sdk.NewEvent(types.OutboundTxFailed,
		sdk.NewAttribute(types.CctxIndex, receive.SendHash),
		sdk.NewAttribute(types.OutTxHash, receive.OutTxHash),
		sdk.NewAttribute(types.ZetaMint, msg.MMint.String()),
		sdk.NewAttribute(types.OutBoundChain, msg.OutTxChain),
		sdk.NewAttribute(types.OldStatus, oldStatus),
		sdk.NewAttribute(types.NewStatus, newStatus),
	)
	ctx.EventManager().EmitEvent(event)
}
