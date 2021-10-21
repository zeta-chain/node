package keeper

import (
	"context"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateSend(goCtx context.Context, msg *types.MsgCreateSend) (*types.MsgCreateSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// index should be inTxHash
	send, isFound := k.GetSend(ctx, msg.Index)
	if !isFound {
		send = types.Send{
			Index:          msg.Index,
			Creator:        msg.Creator,
			Sender:         msg.Sender,
			SenderChain:    msg.SenderChain,
			Receiver:       msg.Receiver,
			ReceiverChain:  msg.ReceiverChain,
			MBurnt:         msg.MBurnt,
			MMint:          msg.MMint,
			Message:        msg.Message,
			InTxHash:       msg.InTxHash,
			InBlockHeight:  msg.InBlockHeight,
			OutTxHash:      msg.OutTxHash,
			OutBlockHeight: msg.OutBlockHeight,
		}
	}


	k.SetSend(ctx, send)


	return &types.MsgCreateSendResponse{}, nil
}
