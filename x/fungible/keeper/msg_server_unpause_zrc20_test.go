package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_UnpauseZRC20(t *testing.T) {
	t.Run("can unpause status of zrc20", func(t *testing.T) {
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
		zrc20A, zrc20B, zrc20C := sample.EthAddress().String(), sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		fcB := sample.ForeignCoins(t, zrc20B)
		fcB.Paused = true
		k.SetForeignCoins(ctx, fcB)
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20C))
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// can unpause zrc20
		_, err := msgServer.UnpauseZRC20(ctx, types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20A,
			},
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// can unpause already unpaused zrc20
		_, err = msgServer.UnpauseZRC20(ctx, types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20C,
			},
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// can unpause all zrc20
		_, err = msgServer.UnpauseZRC20(ctx, types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
		))
		require.NoError(t, err)
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)
	})

	t.Run("should fail if invalid message", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		invalidMsg := types.NewMsgUnpauseZRC20(admin, []string{})
		require.ErrorIs(t, invalidMsg.ValidateBasic(), sdkerrors.ErrInvalidRequest)

		_, err := msgServer.UnpauseZRC20(ctx, invalidMsg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)

		_, err := msgServer.UnpauseZRC20(ctx, types.NewMsgUnpauseZRC20(
			admin,
			[]string{sample.EthAddress().String()},
		))

		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		zrc20A, zrc20B := sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))

		_, err := msgServer.UnpauseZRC20(ctx, types.NewMsgUnpauseZRC20(
			admin,
			[]string{
				zrc20A,
				sample.EthAddress().String(),
				zrc20B,
			},
		))
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
