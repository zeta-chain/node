package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
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

		assertUnpaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			assert.True(t, found)
			assert.False(t, fc.Paused)
		}
		assertPaused := func(zrc20 string) {
			fc, found := k.GetForeignCoins(ctx, zrc20)
			assert.True(t, found)
			assert.True(t, fc.Paused)
		}

		// setup zrc20
		zrc20A, zrc20B, zrc20C := sample.EthAddress().String(), sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20B))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20C))
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)

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
		assert.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// can unpause zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20A,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		assert.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		// can pause already paused zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20B,
			},
			types.UpdatePausedStatusAction_PAUSE,
		))
		assert.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// can unpause already unpaused zrc20
		_, err = msgServer.UpdateZRC20PausedStatus(ctx, types.NewMsgUpdateZRC20PausedStatus(
			admin,
			[]string{
				zrc20C,
			},
			types.UpdatePausedStatusAction_UNPAUSE,
		))
		assert.NoError(t, err)
		assertUnpaused(zrc20A)
		assertPaused(zrc20B)
		assertUnpaused(zrc20C)

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
		assert.NoError(t, err)
		assertPaused(zrc20A)
		assertPaused(zrc20B)
		assertPaused(zrc20C)

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
		assert.NoError(t, err)
		assertUnpaused(zrc20A)
		assertUnpaused(zrc20B)
		assertUnpaused(zrc20C)
	})

	t.Run("should fail if invalid message", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group1)

		invalidMsg := types.NewMsgUpdateZRC20PausedStatus(admin, []string{}, types.UpdatePausedStatusAction_PAUSE)
		assert.ErrorIs(t, invalidMsg.ValidateBasic(), sdkerrors.ErrInvalidRequest)

		_, err := msgServer.UpdateZRC20PausedStatus(ctx, invalidMsg)
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
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

		assert.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
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
		assert.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
