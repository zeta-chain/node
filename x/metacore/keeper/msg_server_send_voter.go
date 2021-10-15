package keeper

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateSendVoter(goCtx context.Context, msg *types.MsgCreateSendVoter) (*types.MsgCreateSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetSendVoter(ctx, msg.Index)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("index %v already set", msg.Index))
	}

	var sendVoter = types.SendVoter{
		Index:           msg.Index,
		Creator:         msg.Creator,
		Sender:          msg.Sender,
		SenderChainId:   msg.SenderChainId,
		Receiver:        msg.Receiver,
		ReceiverChainId: msg.ReceiverChainId,
		MBurnt:          msg.MBurnt,
		Message:         msg.Message,
		TxHash:          msg.TxHash,
		BlockHeight:     msg.BlockHeight,
	}

	k.SetSendVoter(
		ctx,
		sendVoter,
	)

	sendVoter.Index = ""
	sendVoter.Creator = ""
	hashSend := crypto.Keccak256Hash([]byte(sendVoter.String()))

	send, isFound := k.GetSend(ctx, hashSend.Hex());
	if isFound {
		for _, s := range send.Signers {
			if s == msg.Creator {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("send index %s already set from signer %s", msg.TxHash, msg.Creator))
			}
		}
		send.Signers = append(send.Signers, msg.Creator)
	} else {
		send = types.Send{
			Index:           msg.Index,
			Creator:         msg.Creator,
			Sender:          msg.Sender,
			SenderChainId:   msg.SenderChainId,
			Receiver:        msg.Receiver,
			ReceiverChainId: msg.ReceiverChainId,
			MBurnt:          msg.MBurnt,
			Message:         msg.Message,
			TxHash:          msg.TxHash,
			Signers: 		 []string{msg.Creator},
			FinalizedHeight: 0,
		}
	}

	if len(send.Signers) == 2 {
		send.FinalizedHeight = uint64(ctx.BlockHeader().Height)
		//TODO: generate RECV item.
	}
	return &types.MsgCreateSendVoterResponse{}, nil
}
