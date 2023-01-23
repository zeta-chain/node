package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k msgServer) AddTokenEmission(goCtx context.Context, msg *types.MsgAddTokenEmission) (*types.MsgAddTokenEmissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	tracker := types.EmissionTracker{}
	tracker, found := k.GetEmissionTracker(ctx, msg.Category)
	if !found {
		return &types.MsgAddTokenEmissionResponse{}, types.ErrEmissionTrackerNotFound
	}
	tracker.AmountLeft = tracker.AmountLeft.Add(msg.Amount)
	k.SetEmissionTracker(ctx, &tracker)
	return &types.MsgAddTokenEmissionResponse{}, nil
}
