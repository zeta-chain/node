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

// Deploys a fungible coin from a connected chains as a ZRC20 on ZetaChain.
//
// If this is a gas coin, the following happens:
//
// * ZRC20 contract for the coin is deployed
// * contract address of ZRC20 is set as a token address in the system
// contract
// * ZETA tokens are minted and deposited into the module account
// * setGasZetaPool is called on the system contract to add the information
// about the pool to the system contract
// * addLiquidityETH is called to add liquidity to the pool
//
// If this is a non-gas coin, the following happens:
//
// * ZRC20 contract for the coin is deployed
// * The coin is added to the list of foreign coins in the module's state
//
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) DeployFungibleCoinZRC20(goCtx context.Context, msg *types.MsgDeployFungibleCoinZRC20) (*types.MsgDeployFungibleCoinZRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
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
