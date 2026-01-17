package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/ptr"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_DeployFungibleCoinZRC20(t *testing.T) {
	t.Run("can deploy a new zrc20", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		chainID := getValidChainID(t)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		msg := types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"foo",
			"foo",
			coin.CoinType_Gas,
			1000000,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		res, err := msgServer.DeployFungibleCoinZRC20(ctx, msg)
		require.NoError(t, err)
		gasAddress := res.Address
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(gasAddress))

		// can retrieve the gas coin
		foreignCoin, found := k.GetForeignCoins(ctx, gasAddress)
		require.True(t, found)
		require.Equal(t, foreignCoin.CoinType, coin.CoinType_Gas)
		require.Contains(t, foreignCoin.Name, "foo")

		// check gas limit
		gasLimit, err := k.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress))
		require.NoError(t, err)
		require.Equal(t, uint64(1000000), gasLimit.Uint64())

		gas, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		require.Equal(t, gasAddress, gas.Hex())

		// can deploy non-gas zrc20
		msg = types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"bar",
			"bar",
			coin.CoinType_ERC20,
			2000000,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		res, err = msgServer.DeployFungibleCoinZRC20(ctx, msg)
		require.NoError(t, err)

		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.Address))

		foreignCoin, found = k.GetForeignCoins(ctx, res.Address)
		require.True(t, found)
		require.Equal(t, coin.CoinType_ERC20, foreignCoin.CoinType)
		require.Equal(t, uint64(1000), foreignCoin.LiquidityCap.Uint64())
		require.Contains(t, foreignCoin.Name, "bar")

		// check gas limit
		gasLimit, err = k.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress))
		require.NoError(t, err)
		require.Equal(t, uint64(2000000), gasLimit.Uint64())

		// gas should remain the same
		gas, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		require.NotEqual(t, res.Address, gas.Hex())
		require.Equal(t, gasAddress, gas.Hex())
	})

	t.Run("should not deploy a new zrc20 if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		msg := types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"foo",
			"foo",
			coin.CoinType_Gas,
			1000000,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should not deploy a new zrc20 with wrong decimal", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		msg := types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			78,
			"foo",
			"foo",
			coin.CoinType_Gas,
			1000000,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("should not deploy a new zrc20 with invalid chain ID", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// should not deploy a new zrc20 if not admin
		msg := types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			9999999,
			8,
			"foo",
			"foo",
			coin.CoinType_Gas,
			1000000,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should not deploy an existing gas or erc20 contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		deployMsg := types.NewMsgDeployFungibleCoinZRC20(
			admin,
			sample.EthAddress().Hex(),
			chainID,
			8,
			"foo",
			"foo",
			coin.CoinType_Gas,
			1000000,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, deployMsg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		// Attempt to deploy the same gas token twice should result in error
		_, err := keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		require.NoError(t, err)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, deployMsg, nil)

		_, err = keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinAlreadyExist)

		// Similar to above, redeploying existing erc20 should also fail
		deployMsg.CoinType = coin.CoinType_ERC20
		keepertest.MockCheckAuthorization(&authorityMock.Mock, deployMsg, nil)

		_, err = keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		require.NoError(t, err)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, deployMsg, nil)

		_, err = keeper.NewMsgServerImpl(*k).DeployFungibleCoinZRC20(ctx, deployMsg)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinAlreadyExist)
	})
}
