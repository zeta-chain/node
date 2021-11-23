package keeper

import (
	"context"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ReceiveConfirmation(goCtx context.Context, msg *types.MsgReceiveConfirmation) (*types.MsgReceiveConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator"))
	}

	index := msg.SendHash
	send, isFound := k.GetSend(ctx, index)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("sendHash %s does not exist", index))
	}

	if msg.Status != common.ReceiveStatus_Failed {
		if msg.MMint != send.MMint {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("MMint %s does not match send MMint %s", msg.MMint, send.MMint))
		}
	}

	receiveIndex := msg.Digest()
	receive, isFound := k.GetReceive(ctx, receiveIndex)

	if isDuplicateSigner(msg.Creator, receive.Signers) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s double signing!!", msg.Creator))
	}

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
		receive.Signers = append(receive.Signers, msg.Creator)
	}

	if hasSuperMajorityValidators(len(receive.Signers), validators){
		receive.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		lastblock, isFound := k.GetLastBlockHeight(ctx, send.ReceiverChain)
		if !isFound {
			lastblock = types.LastBlockHeight{
				Creator:           msg.Creator,
				Index:             send.ReceiverChain,
				Chain:             send.ReceiverChain,
				LastSendHeight:    0,
				LastReceiveHeight: msg.OutBlockHeight,
			}
		} else {
			lastblock.LastSendHeight = msg.OutBlockHeight
		}
		k.SetLastBlockHeight(ctx, lastblock)

		if receive.Status == common.ReceiveStatus_Success {
			if send.Status == types.SendStatus_Abort {
				send.Status = types.SendStatus_Reverted
			} else {
				send.Status = types.SendStatus_Mined
			}
		} else if receive.Status == common.ReceiveStatus_Failed {
			if send.Status == types.SendStatus_Finalized {
				send.Status = types.SendStatus_Abort
			}
		}
		k.SetSend(ctx, send)

	}

	k.SetReceive(ctx, receive)
	return &types.MsgReceiveConfirmationResponse{}, nil
}
