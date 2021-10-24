package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	index := msg.Digest()
	send, isFound := k.GetSend(ctx, index)
	if isFound { // send exists; add creator to signers
		send.Signers = append(send.Signers, msg.Creator)
	} else {
		send = types.Send{
			Creator:             msg.Creator,
			Index:               index,
			Sender:              msg.Sender,
			SenderChain:         msg.SenderChain,
			Receiver:            msg.Receiver,
			ReceiverChain:       msg.ReceiverChain,
			MBurnt:              msg.MBurnt,
			MMint:               msg.MMint,
			Message:             msg.Message,
			InTxHash:            msg.InTxHash,
			InBlockHeight:       msg.InBlockHeight,
			FinalizedMetaHeight: 0,
			Signers:             []string{msg.Creator},
		}
	}

	//TODO: proper super majority needed
	if len(send.Signers) == 2 {
		send.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		lastblock, isFound := k.GetLastBlockHeight(ctx, msg.SenderChain)
		if !isFound {
			lastblock = types.LastBlockHeight{
				Creator:           msg.Creator,
				Index:             msg.SenderChain,
				Chain:             msg.SenderChain,
				LastSendHeight:    msg.InBlockHeight,
				LastReceiveHeight: 0,
			}
		} else {
			lastblock.LastSendHeight = msg.InBlockHeight
		}
		k.SetLastBlockHeight(ctx, lastblock)
	}

	k.SetSend(ctx, send)

	return &types.MsgSendVoterResponse{}, nil
}
