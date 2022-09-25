package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k msgServer) DeployGasPriceOracle(goCtx context.Context, msg *types.MsgDeployGasPriceOracle) (*types.MsgDeployGasPriceOracleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDeployGasPriceOracleResponse{}, nil
}
