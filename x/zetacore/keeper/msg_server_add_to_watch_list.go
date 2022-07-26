package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) AddToWatchList(goCtx context.Context, msg *types.MsgAddToWatchList) (*types.MsgAddToWatchListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddToWatchListResponse{}, nil
}
