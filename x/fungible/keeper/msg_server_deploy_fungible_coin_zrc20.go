package keeper

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	zetacommon "github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// DeployFungibleCoinZRC20 deploys a fungible coin from a connected chains as a ZRC20 on ZetaChain.
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
		// #nosec G701 always in range
		_, err := k.SetupChainGasCoinAndPool(ctx, msg.ForeignChainId, msg.Name, msg.Symbol, uint8(msg.Decimals))
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
		}
	} else {
		// #nosec G701 always in range
		addr, err := k.DeployZRC20Contract(ctx, msg.Name, msg.Symbol, uint8(msg.Decimals), msg.ForeignChainId, msg.CoinType, msg.ERC20, big.NewInt(msg.GasLimit))
		if err != nil {
			return nil, err
		}

		err = ctx.EventManager().EmitTypedEvent(
			&types.EventZRC20Deployed{
				MsgTypeUrl: sdk.MsgTypeURL(&types.MsgDeployFungibleCoinZRC20{}),
				ChainId:    msg.ForeignChainId,
				Contract:   addr.String(),
				Name:       msg.Name,
				Symbol:     msg.Symbol,
				// #nosec G701 always in range
				Decimals: int64(msg.Decimals),
				CoinType: msg.CoinType,
				Erc20:    msg.ERC20,
				GasLimit: msg.GasLimit,
			},
		)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to emit event")
		}

	}

	return &types.MsgDeployFungibleCoinZRC20Response{}, nil
}
