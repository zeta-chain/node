package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/contracts"
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
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.ZRC20DepositAndCallContract(
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
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("can deposit non-gas coin for transfers", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.ZRC20DepositAndCallContract(
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
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should fail if trying to call a contract with data to a EOC", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		deployZRC20(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte("DEADBEEF"),
			common.CoinType_ERC20,
			assetAddress,
		)
		require.ErrorIs(t, err, types.ErrCallNonContract)
	})

	t.Run("can deposit coin for transfers with liquidity cap not reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

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
		_, contractCall, err := k.ZRC20DepositAndCallContract(
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
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(500), balance)
	})

	t.Run("should fail if coin paused", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// pause the coin
		coin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		coin.Paused = true
		k.SetForeignCoins(ctx, coin)

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
		require.ErrorIs(t, err, types.ErrPausedZRC20)
	})

	t.Run("should fail if liquidity cap reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

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
		chain := chainList[0].ChainId

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
		chain := chainList[0].ChainId
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

	t.Run("should return contract call if receiver is a contract", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		example, err := k.DeployContract(ctx, contracts.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, example)

		// deposit
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			example,
			big.NewInt(42),
			chain,
			[]byte{},
			common.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.NoError(t, err)
		require.True(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, example)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)

		// check onCrossChainCall() hook was called
		assertExampleBarValue(t, ctx, k, example, 42)
	})

	t.Run("should fail if call contract fails", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := common.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		reverter, err := k.DeployContract(ctx, contracts.ReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, reverter)

		// deposit
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			reverter,
			big.NewInt(42),
			chain,
			[]byte{},
			common.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.Error(t, err)
		require.True(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, reverter)
		require.NoError(t, err)
		require.EqualValues(t, int64(0), balance.Int64())
	})
}
