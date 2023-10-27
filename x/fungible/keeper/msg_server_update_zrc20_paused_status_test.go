package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_UpdateZRC20PausedStatus(t *testing.T) {
	t.Run("can update the paused status of zrc20", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		requireUnpaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.False(t, fc.Paused)
		}
		requirePaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			require.True(t, found)
			require.True(t, fc.Paused)
		}

		// setup zrc20
		zrc20A, zrc20B, zrc20C := sample.EthAddress().String(), sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20C))
		requireUnpaused(zrc20A)
		requireUnpaused(zrc20B)
		requireUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		// can pause zrc20
		_, err := msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.NoError(t, err)
		requirePaused(zrc20A)
		requirePaused(zrc20B)
		requireUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// can unpause zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		require.NoError(t, err)
		requireUnpaused(zrc20A)
		requirePaused(zrc20B)
		requireUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		// can pause already paused zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.NoError(t, err)
		requireUnpaused(zrc20A)
		requirePaused(zrc20B)
		requireUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// can unpause already unpaused zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20C,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		require.NoError(t, err)
		requireUnpaused(zrc20A)
		requirePaused(zrc20B)
		requireUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		// can pause all zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.NoError(t, err)
		requirePaused(zrc20A)
		requirePaused(zrc20B)
		requirePaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// can unpause all zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				zrc20B,
				zrc20C,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		require.NoError(t, err)
		requireUnpaused(zrc20A)
		requireUnpaused(zrc20B)
		requireUnpaused(zrc20C)
	})

	t.Run("should fail if invalid message", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		invalidMsg := types.NewMsgUpdateZRC20PausedStatus(admin, []string{}, types.UpdatePausedStatusAction_PAUSE)
		require.ErrorIs(t, invalidMsg.ValidateBasic(), sdkerrors.ErrInvalidRequest)

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, invalidMsg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			sample.AccAddress(),
			[]string{sample.EthAddress().String()},
			types.UpdatePausedStatusAction_PAUSE,
		))

		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			sample.AccAddress(),
			[]string{sample.EthAddress().String()},
			types.UpdatePausedStatusAction_UNPAUSE,
		))

		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if zrc20 does not exist", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		zrc20A, zrc20B := sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
				sample.EthAddress().String(),
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
