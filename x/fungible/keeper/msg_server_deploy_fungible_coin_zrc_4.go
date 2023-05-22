package keeper

import (
	"context"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	zetacommon "github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) DeployFungibleCoinZRC20(goCtx context.Context, msg *types.MsgDeployFungibleCoinZRC20) (*types.MsgDeployFungibleCoinZRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaobserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Deploy can only be executed by the correct policy account")
	}
	if msg.Decimals > 255 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "decimals must be less than 256")
	}
	if msg.CoinType == zetacommon.CoinType_Gas {
		_, err := k.setupChainGasCoinAndPool(ctx, msg.ForeignChain, msg.Name, msg.Symbol, uint8(msg.Decimals))
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
		}
	} else {
		addr, err := k.DeployZRC20Contract(ctx, msg.Name, msg.Symbol, uint8(msg.Decimals), msg.ForeignChain, msg.CoinType, msg.ERC20, big.NewInt(msg.GasLimit))
		if err != nil {
			return nil, err
		}

		//FIXME : declare the attributes as constants , in x/fungible/types
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute("action", "DeployFungibleCoinZRC20"),
				sdk.NewAttribute("chain", msg.ForeignChain),
				sdk.NewAttribute("contract", addr.String()),
				sdk.NewAttribute("name", msg.Name),
				sdk.NewAttribute("symbol", msg.Symbol),
				sdk.NewAttribute("decimals", fmt.Sprintf("%d", msg.Decimals)),
				sdk.NewAttribute("coinType", msg.CoinType.String()),
				sdk.NewAttribute("erc20", msg.ERC20),
				sdk.NewAttribute("gasLimit", fmt.Sprintf("%d", msg.GasLimit)),
			),
		)
	}

	return &types.MsgDeployFungibleCoinZRC20Response{}, nil
}
