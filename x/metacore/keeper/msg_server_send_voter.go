package keeper

import (
	"context"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain, err := common.ParseChain(msg.ReceiverChain)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("cannot parse chain %s", msg.ReceiverChain))
	}
	nonce, isFound := k.GetChainNonces(ctx, chain.String())
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("no chain nonce"))
	}
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
			Status:              types.SendStatus_Created,
			Nonce:               0,
		}
	}

	//TODO: proper super majority needed
	if len(send.Signers) == 2 {
		send.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		send.Status = types.SendStatus_Finalized
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

		send.Broadcaster = uint64(rand.Intn(len(send.Signers)))
		// TODO: subtract gas fee from here
		send.MMint = send.MBurnt
		send.Nonce = nonce.Nonce
		nonce.Nonce++
		k.SetChainNonces(ctx, nonce)
	}

	k.SetSend(ctx, send)

	return &types.MsgSendVoterResponse{}, nil
}
