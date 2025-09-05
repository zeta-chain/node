package v4_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/ptr"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	v4 "github.com/zeta-chain/node/x/fungible/migrations/v4"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestGetSuiChain(t *testing.T) {
	t.Run("returns SuiMainnet for ZetaChainMainnet", func(t *testing.T) {
		// ARRANGE
		chainID := chains.ZetaChainMainnet.ChainId

		// ACT
		chain, err := v4.GetSuiChain(chainID)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, chains.SuiMainnet, chain)
	})

	t.Run("returns SuiTestnet for ZetaChainTestnet", func(t *testing.T) {
		// ARRANGE
		chainID := chains.ZetaChainTestnet.ChainId

		// ACT
		chain, err := v4.GetSuiChain(chainID)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, chains.SuiTestnet, chain)
	})

	t.Run("returns SuiLocalnet for ZetaChainPrivnet", func(t *testing.T) {
		// ARRANGE
		chainID := chains.ZetaChainPrivnet.ChainId

		// ACT
		chain, err := v4.GetSuiChain(chainID)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, chains.SuiLocalnet, chain)
	})

	t.Run("returns error for unsupported chain ID", func(t *testing.T) {
		// ARRANGE
		unsupportedChainID := int64(999999)

		// ACT
		chain, err := v4.GetSuiChain(unsupportedChainID)

		// ASSERT
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported chain ID: 999999")
		require.Equal(t, chains.Chain{}, chain)
	})
}

func TestMigrateStore(t *testing.T) {
	t.Run("successful migration burns SUI balance from stability pool", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		ctx = ctx.WithChainID("zetachain_7000-1")
		testLogger := sample.NewTestLogger()
		ctx = ctx.WithLogger(testLogger)
		chainID := chains.SuiMainnet.ChainId
		deploySystemContracts(t, ctx, k)
		_ = setupGasCoin(t, ctx, k, chainID, "SUI", "SUI")
		ethGasZRC20 := setupGasCoin(t, ctx, k, chains.Ethereum.ChainId, "ETH", "ETH")
		k.EnsureGasStabilityPoolAccountCreated(ctx)

		suiGasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		stabilityPoolAddress := types.GasStabilityPoolAddressEVM()
		suiBalance := big.NewInt(1000000)

		_, err = k.DepositZRC20(ctx, suiGasZRC20, stabilityPoolAddress, suiBalance)
		require.NoError(t, err)
		_, err = k.DepositZRC20(ctx, ethGasZRC20, stabilityPoolAddress, suiBalance)
		require.NoError(t, err)
		fetchedBalance, err := k.ZRC20BalanceOf(ctx, suiGasZRC20, stabilityPoolAddress)
		require.NoError(t, err)
		require.Equal(t, suiBalance, fetchedBalance)

		// ACT
		err = v4.MigrateStore(ctx, k)

		// ASSERT
		require.NoError(t, err)
		logOutput := testLogger.String()
		require.Contains(t, logOutput, "SUI gas ZRC20 burned from stability pool", "Expected log for successful burn of SUI gas ZRC20")
		fetchedBalance, err = k.ZRC20BalanceOf(ctx, suiGasZRC20, stabilityPoolAddress)
		require.NoError(t, err)
		require.Equal(t, int64(0), fetchedBalance.Int64(), "SUI balance should be zero after migration")

		fetchedBalance, err = k.ZRC20BalanceOf(ctx, ethGasZRC20, stabilityPoolAddress)
		require.NoError(t, err)
		require.Equal(t, int64(1000000), fetchedBalance.Int64(), "Eth balance should remain after migration")
	})

	t.Run("logs error and returns nil if gas zrc20 not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		ctx = ctx.WithChainID("zetachain_7000-1")
		testLogger := sample.NewTestLogger()
		ctx = ctx.WithLogger(testLogger)
		deploySystemContracts(t, ctx, k)
		ethGasZRC20 := setupGasCoin(t, ctx, k, chains.Ethereum.ChainId, "ETH", "ETH")
		k.EnsureGasStabilityPoolAccountCreated(ctx)

		stabilityPoolAddress := types.GasStabilityPoolAddressEVM()
		suiBalance := big.NewInt(1000000)

		_, err := k.DepositZRC20(ctx, ethGasZRC20, stabilityPoolAddress, suiBalance)
		require.NoError(t, err)

		// ACT
		err = v4.MigrateStore(ctx, k)

		// ASSERT
		require.NoError(t, err)
		logOutput := testLogger.String()
		require.Contains(t, logOutput, "failed to query SUI gas coin ZRC20", "Expected error log for missing SUI gas coin")
	})

	t.Run("logs error and returns nil if unable to burn tokens", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		testLogger := sample.NewTestLogger()
		ctx = ctx.WithLogger(testLogger)
		ctx = ctx.WithChainID("zetachain_7000-1")
		deploySystemContracts(t, ctx, k)
		chainID := chains.SuiMainnet.ChainId
		_ = setupGasCoin(t, ctx, k, chainID, "SUI", "SUI")

		// ACT
		err := v4.MigrateStore(ctx, k)

		// ASSERT
		require.NoError(t, err)
		logOutput := testLogger.String()
		require.Contains(t, logOutput, "failed to burn SUI gas ZRC20 from stability pool", "Expected error log for burn failure")
	})
}

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	var err error

	wzeta, err = k.DeployWZETA(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wzeta)

	uniswapV2Factory, err = k.DeployUniswapV2Factory(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Factory)

	uniswapV2Router, err = k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Router)

	connector, err = k.DeployConnectorZEVM(ctx, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, connector)

	systemContract, err = k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, uniswapV2Router)
	require.NoError(t, err)
	require.NotEmpty(t, systemContract)

	return
}

func setupGasCoin(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	chainID int64,
	assetName string,
	symbol string,
) (zrc20 common.Address) {
	addr, err := k.SetupChainGasCoinAndPool(
		ctx,
		chainID,
		assetName,
		symbol,
		8,
		nil,
		ptr.Ptr(sdkmath.NewUint(1000)),
	)
	require.NoError(t, err)

	// increase the default liquidity cap
	foreignCoin, found := k.GetForeignCoins(ctx, addr.Hex())
	require.True(t, found)
	foreignCoin.LiquidityCap = sdkmath.NewUint(1e18).MulUint64(1e12)
	k.SetForeignCoins(ctx, foreignCoin)

	return addr
}
