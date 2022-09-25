package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k msgServer) DeployGasPriceOracle(goCtx context.Context, msg *types.MsgDeployGasPriceOracle) (*types.MsgDeployGasPriceOracleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := k.DeployGasPriceOracleContract(ctx)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("action", "DeployGasPriceOracle"),
			sdk.NewAttribute("contract", addr.String()),
		),
	)
	return &types.MsgDeployGasPriceOracleResponse{}, nil
}
