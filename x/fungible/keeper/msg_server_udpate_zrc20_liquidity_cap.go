package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateZRC20LiquidityCap updates the liquidity cap for a ZRC20 token.
//
// Authorized: admin policy group 2.
func (k msgServer) UpdateZRC20LiquidityCap(goCtx context.Context, msg *types.MsgUpdateZRC20LiquidityCap) (*types.MsgUpdateZRC20LiquidityCapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check authorization
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_group2) {
		return nil, cosmoserrors.Wrap(sdkerrors.ErrUnauthorized, "update can only be executed by group 2 policy group")
	}

	// fetch the foreign coin
	coin, found := k.GetForeignCoins(ctx, msg.Zrc20Address)
	if !found {
		return nil, types.ErrForeignCoinNotFound
	}

	// update the liquidity cap
	coin.LiquidityCap = msg.LiquidityCap
	k.SetForeignCoins(ctx, coin)

	return &types.MsgUpdateZRC20LiquidityCapResponse{}, nil
}
