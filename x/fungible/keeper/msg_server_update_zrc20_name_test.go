package keeper_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
	"testing"
)

func TestMsgServer_UpdateZRC20Name(t *testing.T) {
	t.Run("should fail if not admin", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20Address := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "ZRC20", "ZRC20")

		msg := types.NewMsgUpdateZRC20Name(
			admin,
			zrc20Address.Hex(),
			"foo",
			"bar",
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		// act
		_, err := msgServer.UpdateZRC20Name(ctx, msg)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if invalid zrc20 address", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "ZRC20", "ZRC20")

		msg := types.NewMsgUpdateZRC20Name(
			admin,
			"invalid",
			"foo",
			"bar",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err := msgServer.UpdateZRC20Name(ctx, msg)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if non existent zrc20", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20Address := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "ZRC20", "ZRC20")
		k.RemoveForeignCoins(ctx, zrc20Address.Hex())

		msg := types.NewMsgUpdateZRC20Name(
			admin,
			zrc20Address.Hex(),
			"foo",
			"bar",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err := msgServer.UpdateZRC20Name(ctx, msg)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("can update name and symbol", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20Address := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "ZRC20", "ZRC20")

		msg := types.NewMsgUpdateZRC20Name(
			admin,
			zrc20Address.Hex(),
			"foo",
			"bar",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err := msgServer.UpdateZRC20Name(ctx, msg)

		// assert
		require.NoError(t, err)

		// check the name and symbol
		name, err := k.ZRC20Name(ctx, zrc20Address)
		require.NoError(t, err)
		require.Equal(t, "foo", name)

		symbol, err := k.ZRC20Symbol(ctx, zrc20Address)
		require.NoError(t, err)
		require.Equal(t, "bar", symbol)

		// check object
		fc, found := k.GetForeignCoins(ctx, zrc20Address.Hex())
		require.True(t, found)
		require.Equal(t, "foo", fc.Name)
		require.Equal(t, "bar", fc.Symbol)

		// can update name only
		// arrange
		msg = types.NewMsgUpdateZRC20Name(
			admin,
			zrc20Address.Hex(),
			"foo2",
			"",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err = msgServer.UpdateZRC20Name(ctx, msg)

		// assert
		require.NoError(t, err)

		name, err = k.ZRC20Name(ctx, zrc20Address)
		require.NoError(t, err)
		require.Equal(t, "foo2", name)

		symbol, err = k.ZRC20Symbol(ctx, zrc20Address)
		require.NoError(t, err)
		require.Equal(t, "bar", symbol)

		// check object
		fc, found = k.GetForeignCoins(ctx, zrc20Address.Hex())
		require.True(t, found)
		require.Equal(t, "foo2", fc.Name)
		require.Equal(t, "bar", fc.Symbol)

		// can update symbol only
		// arrange
		msg = types.NewMsgUpdateZRC20Name(
			admin,
			zrc20Address.Hex(),
			"",
			"bar2",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err = msgServer.UpdateZRC20Name(ctx, msg)

		// assert
		require.NoError(t, err)

		name, err = k.ZRC20Name(ctx, zrc20Address)
		require.NoError(t, err)
		require.Equal(t, "foo2", name)

		symbol, err = k.ZRC20Symbol(ctx, zrc20Address)
		require.NoError(t, err)
		require.Equal(t, "bar2", symbol)

		// check object
		fc, found = k.GetForeignCoins(ctx, zrc20Address.Hex())
		require.True(t, found)
		require.Equal(t, "foo2", fc.Name)
		require.Equal(t, "bar2", fc.Symbol)
	})
}
