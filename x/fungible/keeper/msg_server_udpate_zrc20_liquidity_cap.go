package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// UpdateZRC20LiquidityCap updates the liquidity cap for a ZRC20 token.
//
// Authorized: admin policy group 2.
func (k msgServer) UpdateZRC20LiquidityCap(goCtx context.Context, msg *types.MsgUpdateZRC20LiquidityCap) (*types.MsgUpdateZRC20LiquidityCapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check authorization
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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
