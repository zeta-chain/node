package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_DeployFungibleCoinZRC20(t *testing.T) {
	t.Run("can deploy a new zrc20", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		res, err := msgServer.DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"foo",
			"foo",
			common.CoinType_Gas,
			1000000,
		))
		require.NoError(t, err)
		gasAddress := res.Address
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(gasAddress))

		// can retrieve the gas coin
		foreignCoin, found := k.GetForeignCoins(ctx, gasAddress)
		require.True(t, found)
		require.Equal(t, foreignCoin.CoinType, common.CoinType_Gas)
		require.Contains(t, foreignCoin.Name, "foo")

		gas, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		require.Equal(t, gasAddress, gas.Hex())

		// can deploy non-gas zrc20
		res, err = msgServer.DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"bar",
			"bar",
			common.CoinType_ERC20,
			1000000,
		))
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.Address))

		foreignCoin, found = k.GetForeignCoins(ctx, res.Address)
		require.True(t, found)
		require.Equal(t, foreignCoin.CoinType, common.CoinType_ERC20)
		require.Contains(t, foreignCoin.Name, "bar")

		// gas should remain the same
		gas, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		require.NotEqual(t, res.Address, gas.Hex())
		require.Equal(t, gasAddress, gas.Hex())
	})

	t.Run("should not deploy a new zrc20 if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			sample.AccAddress(),
			sample.EthAddress().Hex(),
			chainID,
			8,
			"foo",
			"foo",
			common.CoinType_Gas,
			1000000,
		))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should not deploy a new zrc20 with wrong decimal", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			256,
			"foo",
			"foo",
			common.CoinType_Gas,
			1000000,
		))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("should not deploy a new zrc20 with invalid chain ID", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			9999999,
			8,
			"foo",
			"foo",
			common.CoinType_Gas,
			1000000,
		))
		require.Error(t, err)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})
}
