package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"math/big"
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
	send, isFound := k.GetSendMultipleStatus(ctx, index, []types.SendStatus{types.SendStatus_PendingOutbound, types.SendStatus_PendingRevert})
	if !isFound {
		log.Error().Msgf("Cannot find broadcast tx hash %s on %s chain", index, msg.Chain)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("Cannot find broadcast tx hash %s on %s chain", index, msg.Chain))
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

		zetaBurnt, ok := big.NewInt(0).SetString(send.ZetaBurnt, 10)
		if !ok {
			log.Error().Msgf("ReceiveConfirmation: failed to parse ZetaBurnt: %s", send.ZetaBurnt)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse ZetaBurnt: %s", send.ZetaBurnt))
		}
		zetaMinted, ok := big.NewInt(0).SetString(send.ZetaMint, 10)
		if !ok {
			log.Error().Msgf("ReceiveConfirmation: failed to parse ZetaMint: %s", send.ZetaMint)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse ZetaMint: %s", send.ZetaMint))
		}

		var oldstatus, newstatus types.SendStatus

		if receive.Status == common.ReceiveStatus_Success {
			oldstatus = send.Status
			if send.Status == types.SendStatus_PendingRevert {
				send.Status = types.SendStatus_Reverted
			} else if send.Status == types.SendStatus_PendingOutbound {
				send.Status = types.SendStatus_OutboundMined
			}

			err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(common.ZETADenom, sdk.NewIntFromBigInt(zetaBurnt.Sub(zetaBurnt, zetaMinted)))))
			if err != nil {
				log.Error().Msgf("ReceiveConfirmation: failed to mint coins: %s", err.Error())
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to mint coins: %s", err.Error()))
			}
			newstatus = send.Status
			event := sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
				sdk.NewAttribute(types.SubTypeKey, string(types.OutboundTxSuccessful)),
				sdk.NewAttribute(types.SendHash, receive.SendHash),
				sdk.NewAttribute(types.OutTxHash, receive.OutTxHash),
				sdk.NewAttribute(types.ZetaMint, msg.MMint),
				sdk.NewAttribute(types.Chain, msg.Chain),
				sdk.NewAttribute(types.OldStatus, oldstatus.String()),
				sdk.NewAttribute(types.NewStatus, newstatus.String()),
			)
			ctx.EventManager().EmitEvent(event)
		} else if receive.Status == common.ReceiveStatus_Failed {
			oldstatus = send.Status
			if send.Status == types.SendStatus_PendingOutbound {
				send.Status = types.SendStatus_PendingRevert
				send.StatusMessage = fmt.Sprintf("destination tx %s failed", msg.OutTxHash)
				chain := send.SenderChain
				k.updateSend(ctx, chain, &send)
			} else if send.Status == types.SendStatus_PendingRevert {
				send.Status = types.SendStatus_Aborted
				send.StatusMessage = fmt.Sprintf("revert tx %s failed", msg.OutTxHash)
			}
			newstatus = send.Status
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
					sdk.NewAttribute(types.SubTypeKey, types.OutboundTxFailed),
					sdk.NewAttribute(types.SendHash, receive.SendHash),
					sdk.NewAttribute(types.OutTxHash, receive.OutTxHash),
					sdk.NewAttribute(types.ZetaMint, send.ZetaMint),
					sdk.NewAttribute(types.Chain, msg.Chain),
					sdk.NewAttribute(types.OldStatus, oldstatus.String()),
					sdk.NewAttribute(types.NewStatus, newstatus.String()),
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
		k.SendMigrateStatus(ctx, send, oldstatus)

	}
	k.SetReceive(ctx, receive)
	return &types.MsgReceiveConfirmationResponse{}, nil
}
