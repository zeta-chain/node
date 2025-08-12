package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_ProcessRevert(t *testing.T) {
	t.Run("should process NoAssetCall revert with callOnRevert true", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(0)
		revertMessage := []byte("revert message")
		callOnRevert := true

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_NoAssetCall,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		assertTestDAppV2MessageAndAmount(t, ctx, k, revertAddress, string(revertMessage), 0)
		balance := sdkk.BankKeeper.GetBalance(ctx, revertAddress.Bytes(), "azeta")
		require.Equal(t, amount.Int64(), balance.Amount.Int64())
	})

	t.Run("should process NoAssetCall revert with callOnRevert false", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertAddress := sample.EthAddress()
		revertMessage := []byte("revert message")
		callOnRevert := false

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_NoAssetCall,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, resp)
	})

	t.Run("should fail if system contracts not deployed", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertAddress := sample.EthAddress()
		revertMessage := []byte("revert message")
		callOnRevert := true

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Zeta,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, types.ErrSystemContractNotFound)
	})

	t.Run("should fail for unsupported coin type", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertAddress := sample.EthAddress()
		revertMessage := []byte("revert message")
		callOnRevert := true

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Cmd,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "unsupported coin type for revert")
	})

	t.Run("should process Zeta revert with callOnRevert true", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := true

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Zeta,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		assertTestDAppV2MessageAndAmount(t, ctx, k, revertAddress, string(revertMessage), 0)
		balance := sdkk.BankKeeper.GetBalance(ctx, revertAddress.Bytes(), "azeta")
		require.Equal(t, amount.Int64(), balance.Amount.Int64())
	})

	t.Run("should process Zeta revert with callOnRevert false", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := false

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Zeta,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)

		balance := sdkk.BankKeeper.GetBalance(ctx, revertAddress.Bytes(), "azeta")
		require.Equal(t, amount.Int64(), balance.Amount.Int64())
	})

	t.Run("should process ERC20 revert with callOnRevert true", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := true
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		assertTestDAppV2MessageAndAmount(t, ctx, k, revertAddress, string(revertMessage), 0)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, revertAddress)
		require.NoError(t, err)
		require.Equal(t, amount.Int64(), balance.Int64())
	})

	t.Run("should process ERC20 revert with callOnRevert false", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := false
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, revertAddress)
		require.NoError(t, err)
		require.Equal(t, amount.Int64(), balance.Int64())
	})

	t.Run("should process Gas revert with callOnRevert true", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := true

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Gas,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		assertTestDAppV2MessageAndAmount(t, ctx, k, revertAddress, string(revertMessage), 0)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, revertAddress)
		require.NoError(t, err)
		require.Equal(t, amount.Int64(), balance.Int64())
	})

	t.Run("should process Gas revert with callOnRevert false", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := false

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress(), sample.EthAddress())
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testDAppV2)
		revertAddress := testDAppV2

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Gas,
			"",
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, revertAddress)
		require.NoError(t, err)
		require.Equal(t, amount.Int64(), balance.Int64())
	})

	t.Run("should fail if liquidity cap reached for ERC20", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := true
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		foreignCoin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		foreignCoin.LiquidityCap = math.NewUint(50)
		k.SetForeignCoins(ctx, foreignCoin)

		revertAddress := sample.EthAddress()

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			revertAddress,
			callOnRevert,
			revertMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, types.ErrForeignCoinCapReached)
	})
}
