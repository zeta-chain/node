package keeper

import (
	"context"
	"fmt"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTxinVoter(goCtx context.Context, msg *types.MsgCreateTxinVoter) (*types.MsgCreateTxinVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetTxinVoter(ctx, msg.Index)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("index %v already set", msg.Index))
	}

	var txinVoter = types.TxinVoter{
		Index:            msg.Index,
		Creator:          msg.Creator,
		TxHash:           msg.TxHash,
		SourceAsset:      msg.SourceAsset,
		SourceAmount:     msg.SourceAmount,
		MBurnt:           msg.MBurnt,
		DestinationAsset: msg.DestinationAsset,
		FromAddress:      msg.FromAddress,
		ToAddress:        msg.ToAddress,
		BlockHeight:      msg.BlockHeight,
	}

	k.SetTxinVoter(
		ctx,
		txinVoter,
	)
	return &types.MsgCreateTxinVoterResponse{}, nil
}
