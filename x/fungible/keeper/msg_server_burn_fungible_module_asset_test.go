package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_BurnFungibleModuleAsset(t *testing.T) {
	t.Run("can burn the asset on the fungible module", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// set coin admin
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		// deploy the system contract and a ZRC20 contract
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20Addr := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "alpha", "alpha")

		// send tokens to the fungible module
		amount := big.NewInt(1000)
		_, err := k.DepositZRC20(
			ctx,
			zrc20Addr,
			types.ModuleAddressEVM,
			amount,
		)
		require.NoError(t, err)

		// check the balance of the fungible module
		balance, err := k.ZRC20BalanceOf(ctx, zrc20Addr, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.Equal(t, amount.Uint64(), balance.Uint64())

		// can update the protocol fee and gas limit
		msg := types.NewMsgBurnFungibleModuleAsset(
			admin,
			zrc20Addr.String(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.BurnFungibleModuleAsset(ctx, msg)
		require.NoError(t, err)

		// check the balance of the fungible module after burn
		balance, err = k.ZRC20BalanceOf(ctx, zrc20Addr, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.Zero(t, balance.Uint64(), "balance should be zero after burn")
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgBurnFungibleModuleAsset(
			admin,
			sample.EthAddress().String(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.BurnFungibleModuleAsset(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if invalid zrc20 address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgBurnFungibleModuleAsset(
			admin,
			"invalid_address",
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.BurnFungibleModuleAsset(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if can't retrieve the foreign coin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgBurnFungibleModuleAsset(
			admin,
			sample.EthAddress().String(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.BurnFungibleModuleAsset(ctx, msg)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
