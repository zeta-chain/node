package keeper_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"testing"
)

func TestMsgServer_RemoveForeignCoin(t *testing.T) {
	t.Run("can remove a foreign coin", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", "foo")

		_, found := k.GetForeignCoins(ctx, zrc20.Hex())
		require.True(t, found)

		_, err := msgServer.RemoveForeignCoin(ctx, types.NewMsgRemoveForeignCoin(admin, zrc20.Hex()))
		require.NoError(t, err)
		_, found = k.GetForeignCoins(ctx, zrc20.Hex())
		require.False(t, found)
	})

	t.Run("should fail if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", "foo")

		_, err := msgServer.RemoveForeignCoin(ctx, types.NewMsgRemoveForeignCoin(sample.AccAddress(), zrc20.Hex()))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)

		_, err := msgServer.RemoveForeignCoin(ctx, types.NewMsgRemoveForeignCoin(admin, sample.EthAddress().Hex()))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})
}
