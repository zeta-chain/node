package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k msgServer) DeployFungibleCoinZRC4(goCtx context.Context, msg *types.MsgDeployFungibleCoinZRC4) (*types.MsgDeployFungibleCoinZRC4Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := k.DeployZRC4Contract(ctx, msg.Name, msg.Symbol, uint8(msg.Decimals), msg.ForeignChain, msg.CoinType, msg.ERC20)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("action", "DeployFungibleCoinZRC4"),
			sdk.NewAttribute("chain", msg.ForeignChain),
			sdk.NewAttribute("contract", addr.String()),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("symbol", msg.Symbol),
			sdk.NewAttribute("decimals", fmt.Sprintf("%d", msg.Decimals)),
			sdk.NewAttribute("coinType", msg.CoinType.String()),
			sdk.NewAttribute("erc20", msg.ERC20),
		),
	)

	return &types.MsgDeployFungibleCoinZRC4Response{}, nil
}
