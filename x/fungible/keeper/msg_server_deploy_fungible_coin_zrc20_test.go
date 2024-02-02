package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
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
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)
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
		assert.NoError(t, err)
		gasAddress := res.Address
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(gasAddress))

		// can retrieve the gas coin
		foreignCoin, found := k.GetForeignCoins(ctx, gasAddress)
		assert.True(t, found)
		assert.Equal(t, foreignCoin.CoinType, common.CoinType_Gas)
		assert.Contains(t, foreignCoin.Name, "foo")

		// check gas limit
		gasLimit, err := k.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress))
		assert.NoError(t, err)
		assert.Equal(t, uint64(1000000), gasLimit.Uint64())

		gas, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		assert.NoError(t, err)
		assert.Equal(t, gasAddress, gas.Hex())

		// can deploy non-gas zrc20
		res, err = msgServer.DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"bar",
			"bar",
			common.CoinType_ERC20,
			2000000,
		))
		assert.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.Address))

		foreignCoin, found = k.GetForeignCoins(ctx, res.Address)
		assert.True(t, found)
		assert.Equal(t, foreignCoin.CoinType, common.CoinType_ERC20)
		assert.Contains(t, foreignCoin.Name, "bar")

		// check gas limit
		gasLimit, err = k.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress))
		assert.NoError(t, err)
		assert.Equal(t, uint64(2000000), gasLimit.Uint64())

		// gas should remain the same
		gas, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		assert.NoError(t, err)
		assert.NotEqual(t, res.Address, gas.Hex())
		assert.Equal(t, gasAddress, gas.Hex())
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
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should not deploy a new zrc20 with wrong decimal", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			78,
			"foo",
			"foo",
			common.CoinType_Gas,
			1000000,
		))
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("should not deploy a new zrc20 with invalid chain ID", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

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
		assert.Error(t, err)
		assert.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should not deploy an existing gas or erc20 contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)
		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		deployMsg := types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"foo",
			"foo",
			common.CoinType_Gas,
			1000000,
		)

		// Attempt to deploy the same gas token twice should result in error
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		assert.NoError(t, err)
		_, err = keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		assert.Error(t, err)
		assert.ErrorIs(t, err, types.ErrForeignCoinAlreadyExist)

		// Similar to above, redeploying existing erc20 should also fail
		deployMsg.CoinType = common.CoinType_ERC20
		_, err = keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		assert.NoError(t, err)
		_, err = keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		assert.Error(t, err)
		assert.ErrorIs(t, err, types.ErrForeignCoinAlreadyExist)
	})
}
