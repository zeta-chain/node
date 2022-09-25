package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k msgServer) DeployFungibleCoinZRC4(goCtx context.Context, msg *types.MsgDeployFungibleCoinZRC4) (*types.MsgDeployFungibleCoinZRC4Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDeployFungibleCoinZRC4Response{}, nil
}
