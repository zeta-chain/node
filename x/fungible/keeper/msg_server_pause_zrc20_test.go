package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_PauseZRC20(t *testing.T) {
	t.Run("can pause zrc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		assertUnpaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.False(t, fc.Paused)
		}
		assertPaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.True(t, fc.Paused)
		}

		// setup zrc20
		zrc20A, zrc20B, zrc20C := sample.EthAddress().
			String(),
			sample.EthAddress().
				String(),
			sample.EthAddress().
				String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20C))
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)

		// can pause zrc20
		msg := types.NewMsgPauseZRC20(
			admin,
			[]string{
				zrc20A,
				zrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.PauseZRC20(ctx, msg)
		require.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		// can pause already paused zrc20
		msg = types.NewMsgPauseZRC20(
			admin,
			[]string{
				zrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.PauseZRC20(ctx, msg)
		require.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		// can pause all zrc20
		msg = types.NewMsgPauseZRC20(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.PauseZRC20(ctx, msg)
		require.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertPaused(zrc20C)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgPauseZRC20(
			admin,
			[]string{sample.EthAddress().String()},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.PauseZRC20(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		zrc20A, zrc20B := sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))

		msg := types.NewMsgPauseZRC20(
			admin,
			[]string{
				zrc20A,
				sample.EthAddress().String(),
				zrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.PauseZRC20(ctx, msg)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
