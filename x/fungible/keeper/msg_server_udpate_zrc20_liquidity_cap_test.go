package keeper_test

import (
	"testing"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
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
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

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

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

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

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

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

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

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
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, false)

		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		_, err := msgServer.UpdateZRC20LiquidityCap(ctx, types.NewMsgUpdateZRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
