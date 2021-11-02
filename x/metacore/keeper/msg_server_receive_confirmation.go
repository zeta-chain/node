package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ReceiveConfirmation(goCtx context.Context, msg *types.MsgReceiveConfirmation) (*types.MsgReceiveConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	index := msg.SendHash
	send, isFound := k.GetSend(ctx, index)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("sendHash %s does not exist", index))
	}

	if msg.MMint != send.MMint {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("MMint %s does not match send MMint %s", msg.MMint, send.MMint))
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
		}
	} else {
		receive.Signers = append(receive.Signers, msg.Creator)
	}

	//TODO: do proper super majority check
	if len(receive.Signers) == 2 {
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

		send.Status = types.SendStatus_Mined
		k.SetSend(ctx, send)
	}

	k.SetReceive(ctx, receive)
	return &types.MsgReceiveConfirmationResponse{}, nil
}
