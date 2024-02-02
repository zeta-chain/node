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

func TestMsgServer_RemoveForeignCoin(t *testing.T) {
	t.Run("can remove a foreign coin", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", "foo")

		_, found := k.GetForeignCoins(ctx, zrc20.Hex())
		assert.True(t, found)

		_, err := msgServer.RemoveForeignCoin(ctx, types.NewMsgRemoveForeignCoin(admin, zrc20.Hex()))
		assert.NoError(t, err)
		_, found = k.GetForeignCoins(ctx, zrc20.Hex())
		assert.False(t, found)
	})

	t.Run("should fail if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", "foo")

		_, err := msgServer.RemoveForeignCoin(ctx, types.NewMsgRemoveForeignCoin(sample.AccAddress(), zrc20.Hex()))
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		_, err := msgServer.RemoveForeignCoin(ctx, types.NewMsgRemoveForeignCoin(admin, sample.EthAddress().Hex()))
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})
}
