package keeper

import (
	"context"
	"fmt"

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
	return &types.MsgCreateSendVoterResponse{}, nil
}
