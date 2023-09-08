package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) InTxTrackerAllByChain(goCtx context.Context, request *types.QueryAllInTxTrackerByChainRequest) (*types.QueryAllInTxTrackerByChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	inTxTrackers := k.GetAllInTxTrackerForChain(ctx, request.ChainId)
	return &types.QueryAllInTxTrackerByChainResponse{InTxTracker: inTxTrackers}, nil
}
