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

func (k msgServer) ReceiveConfirmation(goCtx context.Context, msg *types.MsgReceiveConfirmation) (*types.MsgReceiveConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	log.Info().Msgf("ReceiveConfirmation: %s", msg.String())

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		log.Error().Msgf("signer %s is not a bonded validator", msg.Creator)
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	index := msg.SendHash
	send, isFound := k.GetSend(ctx, index)
	if !isFound {
		log.Error().Msgf("send not found: %v", index)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("sendHash %s does not exist", index))
	}

	if msg.Status != common.ReceiveStatus_Failed {
		if msg.MMint != send.ZetaMint {
			log.Error().Msgf("ReceiveConfirmation: Mint mismatch: %s vs %s", msg.MMint, send.ZetaMint)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("MMint %s does not match send ZetaMint %s", msg.MMint, send.ZetaMint))
		}
	}

	receiveIndex := msg.Digest()
	receive, isFound := k.GetReceive(ctx, receiveIndex)

	if !isFound {
		receive = types.Receive{
			Creator:             "",
			Index:               receiveIndex,
			SendHash:            index,
			OutTxHash:           msg.OutTxHash,
			OutBlockHeight:      msg.OutBlockHeight,
			FinalizedMetaHeight: 0,
			Signers:             []string{msg.Creator},
			Status:              msg.Status,
			Chain:               msg.Chain,
		}
	} else {
		if isDuplicateSigner(msg.Creator, receive.Signers) {
			log.Info().Msgf("ReceiveConfirmation: TX %s has already been signed by %s", receiveIndex, msg.Creator)
			return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("ReceiveConfirmation: TX %s has already been signed by %s", receiveIndex, msg.Creator))
		}
		receive.Signers = append(receive.Signers, msg.Creator)
	}

	if hasSuperMajorityValidators(len(receive.Signers), validators) {
		receive.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)

		if receive.Status == common.ReceiveStatus_Success {
			oldstatus := send.Status.String()
			if send.Status == types.SendStatus_PendingRevert {
				send.Status = types.SendStatus_Reverted
			} else if send.Status == types.SendStatus_PendingOutbound {
				send.Status = types.SendStatus_OutboundMined
			}
			newstatus := send.Status.String()
			event := sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
				sdk.NewAttribute(types.SubTypeKey, string(types.OutboundTxSuccessful)),
				sdk.NewAttribute(types.SendHash, receive.SendHash),
				sdk.NewAttribute(types.OutTxHash, receive.OutTxHash),
				sdk.NewAttribute(types.ZetaMint, msg.MMint),
				sdk.NewAttribute(types.Chain, msg.Chain),
				sdk.NewAttribute(types.OldStatus, oldstatus),
				sdk.NewAttribute(types.NewStatus, newstatus),
			)
			ctx.EventManager().EmitEvent(event)
		} else if receive.Status == common.ReceiveStatus_Failed {
			oldstatus := send.Status.String()
			if send.Status == types.SendStatus_PendingOutbound {
				send.Status = types.SendStatus_PendingRevert
				send.StatusMessage = fmt.Sprintf("destination tx %s failed", msg.OutTxHash)
				chain := send.SenderChain
				k.updateSend(ctx, chain, &send)
			} else if send.Status == types.SendStatus_PendingRevert {
				send.Status = types.SendStatus_Aborted
				send.StatusMessage = fmt.Sprintf("revert tx %s failed", msg.OutTxHash)
			}
			newstatus := send.Status.String()
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
					sdk.NewAttribute(types.SubTypeKey, types.OutboundTxFailed),
					sdk.NewAttribute(types.SendHash, receive.SendHash),
					sdk.NewAttribute(types.OutTxHash, receive.OutTxHash),
					sdk.NewAttribute(types.ZetaMint, send.ZetaMint),
					sdk.NewAttribute(types.Chain, msg.Chain),
					sdk.NewAttribute(types.OldStatus, oldstatus),
					sdk.NewAttribute(types.NewStatus, newstatus),
					sdk.NewAttribute(types.StatusMessage, send.StatusMessage),
				),
			)
		}

		if receive.Status == common.ReceiveStatus_Success || receive.Status == common.ReceiveStatus_Failed {
			index := fmt.Sprintf("%s/%s", msg.Chain, strconv.Itoa(int(msg.OutTxNonce)))
			k.RemoveOutTxTracker(ctx, index)
		}

		send.RecvHash = receive.Index
		send.OutTxHash = receive.OutTxHash
		send.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
		k.SetSend(ctx, send)

	}
	k.SetReceive(ctx, receive)
	return &types.MsgReceiveConfirmationResponse{}, nil
}
