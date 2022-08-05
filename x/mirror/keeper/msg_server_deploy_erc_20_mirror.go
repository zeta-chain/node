package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

func (k msgServer) DeployERC20Mirror(goCtx context.Context, msg *types.MsgDeployERC20Mirror) (*types.MsgDeployERC20MirrorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDeployERC20MirrorResponse{}, nil
}
