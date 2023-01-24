package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"

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

	coins := sdk.NewCoin(config.BaseDenom, msg.Amount)
	senderAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgAddTokenEmissionResponse{}, errors.Wrap(types.ErrParsingSenderAddress, err.Error())
	}
	err = k.bankkeeper.SendCoinsFromAccountToModule(ctx, senderAddress, types.ModuleName, sdk.NewCoins(coins))
	if err != nil {
		return &types.MsgAddTokenEmissionResponse{}, errors.Wrap(types.ErrAddingCoinstoTracker, err.Error())
	}
	tracker.AmountLeft = tracker.AmountLeft.Add(msg.Amount)
	k.SetEmissionTracker(ctx, &tracker)
	return &types.MsgAddTokenEmissionResponse{}, nil
}
