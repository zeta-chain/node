package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgServer_UpdateZRC20LiquidityCap(t *testing.T) {
	t.Run("can update the liquidity cap of zrc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		foreignCoin := sample.ForeignCoins(t, coinAddress)
		foreignCoin.LiquidityCap = math.Uint{}
		k.SetForeignCoins(ctx, foreignCoin)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// can update liquidity cap
		msg := types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found := k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.NewUint(42)), "invalid liquidity cap", coin.LiquidityCap.String())

		// can update liquidity cap again
		msg = types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(4200000),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(
			t,
			coin.LiquidityCap.Equal(math.NewUint(4200000)),
			"invalid liquidity cap",
			coin.LiquidityCap.String(),
		)

		// can set liquidity cap to 0
		msg = types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(0),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.ZeroUint()), "invalid liquidity cap", coin.LiquidityCap.String())

		// can set liquidity cap to nil
		msg = types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.ZeroUint()), "invalid liquidity cap", coin.LiquidityCap.String())
	})

	t.Run("should fail if not admin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		foreignCoin := sample.ForeignCoins(t, coinAddress)
		foreignCoin.LiquidityCap = math.Uint{}
		k.SetForeignCoins(ctx, foreignCoin)
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
