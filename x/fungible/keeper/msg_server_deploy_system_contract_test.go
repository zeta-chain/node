package keeper_test

import (
	"context"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgServer_DeploySystemContracts(t *testing.T) {
	t.Run("admin can deploy system contracts", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		res, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.NoError(t, err)
		require.NotNil(t, res)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.UniswapV2Factory))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.Wzeta))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.UniswapV2Router))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.ConnectorZEVM))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.SystemContract))
	})

	t.Run("non-admin cannot deploy system contracts", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		nonadmin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, nonadmin, authoritytypes.PolicyType_groupOperational, false)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(nonadmin))
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if uniswapv2factory contract deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// mock failed uniswapv2factory deployment
		mockFailedContractDeployment(ctx, t, k)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if wzeta contract deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// mock successful uniswapv2factory deployment
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed wzeta deployment deployment
		mockFailedContractDeployment(ctx, t, k)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if uniswapv2router deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// mock successful uniswapv2factory and wzeta deployments
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed uniswapv2router deployment
		mockFailedContractDeployment(ctx, t, k)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if connectorzevm deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// mock successful uniswapv2factory, wzeta and uniswapv2router deployments
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed connectorzevm deployment
		mockFailedContractDeployment(ctx, t, k)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if system contract deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		// mock successful uniswapv2factory, wzeta, uniswapv2router and connectorzevm deployments
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed system contract deployment
		mockFailedContractDeployment(ctx, t, k)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})
}

func mockSuccessfulContractDeployment(ctx context.Context, t *testing.T, k *keeper.Keeper) {
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
	mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
	mockEVMKeeper.On(
		"EstimateGas",
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.EstimateGasResponse{Gas: 5}, nil)
	mockEVMKeeper.MockEVMSuccessCallOnce()
}

func mockFailedContractDeployment(ctx context.Context, t *testing.T, k *keeper.Keeper) {
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
	mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
	mockEVMKeeper.On(
		"EstimateGas",
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.EstimateGasResponse{Gas: 5}, nil)
	mockEVMKeeper.MockEVMFailCallOnce()
}
