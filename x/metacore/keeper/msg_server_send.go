package keeper

import (
	"context"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateSend(goCtx context.Context, msg *types.MsgCreateSend) (*types.MsgCreateSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	send, isFound := k.GetSend(ctx, msg.Index)
	if !isFound {
		send = types.Send{
			Index:          msg.Index,
			Creator:        "",
			Sender:         msg.Sender,
			SenderChain:    msg.SenderChain,
			Receiver:       msg.Receiver,
			ReceiverChain:  msg.ReceiverChain,
			MBurnt:         msg.MBurnt,
			MMint:          "",
			Message:        msg.Message,
			InTxHash:       msg.InTxHash,
			InBlockHeight:  msg.InBlockHeight,
			OutTxHash:      "",
			OutBlockHeight: 0,
			Signers:        []string{msg.Creator},
			InTxFinalizedHeight: 0,
			OutTxFianlizedHeight: 0,
		}
	} else {
		//TODO: check if msg.Creator is already in signers;
		//TODO: verify that msg.Creator is one of the validators;
		send.Signers = append(send.Signers, msg.Creator)

	}




	k.SetSend(ctx, send)


	return &types.MsgCreateSendResponse{}, nil
}
