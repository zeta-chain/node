package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_UpdateZRC20LiquidityCap(t *testing.T) {
	t.Run("can update the liquidity cap of zrc20", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		foreignCoin := sample.ForeignCoins(t, coinAddress)
		foreignCoin.LiquidityCap = math.Uint{}
		k.SetForeignCoins(ctx, foreignCoin)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// can update liquidity cap
		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		))
		require.NoError(t, err)

		coin, found := k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.NewUint(42)), "invalid liquidity cap", coin.LiquidityCap.String())

		// can update liquidity cap again
		_, err = msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(4200000),
		))
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.NewUint(4200000)), "invalid liquidity cap", coin.LiquidityCap.String())

		// can set liquidity cap to 0
		_, err = msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(0),
		))
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.ZeroUint()), "invalid liquidity cap", coin.LiquidityCap.String())

		// can set liquidity cap to nil
		_, err = msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.Uint{},
		))
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.ZeroUint()), "invalid liquidity cap", coin.LiquidityCap.String())
	})

	t.Run("should fail if not admin", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		foreignCoin := sample.ForeignCoins(t, coinAddress)
		foreignCoin.LiquidityCap = math.Uint{}
		k.SetForeignCoins(ctx, foreignCoin)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			sample.AccAddress(),
			coinAddress,
			math.NewUint(42),
		))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
