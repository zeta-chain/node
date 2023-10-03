package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_ZRC20DepositAndCallContract(t *testing.T) {
	t.Run("can deposit gas coin for transfers", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0]

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain.ChainId, "foobar", "foobar")

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			common.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.NoError(t, err)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("can deposit non-gas coin for transfers", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0]
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chain.ChainId, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			common.CoinType_ERC20,
			assetAddress,
		)
		require.NoError(t, err)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("can deposit coin for transfers with liquidity cap not reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0]

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain.ChainId, "foobar", "foobar")

		// there is an initial total supply minted during gas pool setup
		initialTotalSupply, err := k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)

		// set a liquidity cap
		coin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		coin.LiquidityCap = math.NewUint(initialTotalSupply.Uint64() + 1000)
		k.SetForeignCoins(ctx, coin)

		// increase total supply
		_, err = k.DepositZRC20(ctx, zrc20, sample.EthAddress(), big.NewInt(500))
		require.NoError(t, err)

		// deposit
		to := sample.EthAddress()
		_, _, err = k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(500),
			chain,
			[]byte{},
			common.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.NoError(t, err)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(500), balance)
	})

	t.Run("should fail if liquidity cap reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0]

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain.ChainId, "foobar", "foobar")

		// there is an initial total supply minted during gas pool setup
		initialTotalSupply, err := k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)

		// set a liquidity cap
		coin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		coin.LiquidityCap = math.NewUint(initialTotalSupply.Uint64() + 1000)
		k.SetForeignCoins(ctx, coin)

		// increase total supply
		_, err = k.DepositZRC20(ctx, zrc20, sample.EthAddress(), big.NewInt(500))
		require.NoError(t, err)

		// deposit (500 + 501 > 1000)
		to := sample.EthAddress()
		_, _, err = k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(501),
			chain,
			[]byte{},
			common.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.ErrorIs(t, err, types.ErrForeignCoinCapReached)
	})

	t.Run("should fail if gas coin not found", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0]

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			common.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.ErrorIs(t, err, crosschaintypes.ErrGasCoinNotFound)
	})

	t.Run("should fail if zrc20 not found", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0]
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			common.CoinType_ERC20,
			assetAddress,
		)
		require.ErrorIs(t, err, crosschaintypes.ErrForeignCoinNotFound)
	})

	// TODO: add test cases checking DepositZRC20AndCallContract
	// https://github.com/zeta-chain/node/issues/1206
}
