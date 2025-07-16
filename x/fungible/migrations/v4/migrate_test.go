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
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	v4 "github.com/zeta-chain/node/x/fungible/migrations/v4"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("successful migration burns SUI balance from stability pool", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		ctx = ctx.WithChainID("zetachain_7000-1")
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
		fetchedBalance, err = k.ZRC20BalanceOf(ctx, suiGasZRC20, stabilityPoolAddress)
		require.NoError(t, err)
		require.Equal(t, int64(0), fetchedBalance.Int64(), "SUI balance should be zero after migration")

		fetchedBalance, err = k.ZRC20BalanceOf(ctx, ethGasZRC20, stabilityPoolAddress)
		require.NoError(t, err)
		require.Equal(t, int64(1000000), fetchedBalance.Int64(), "Eth balance should remain after migration")
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
