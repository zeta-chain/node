package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_ProcessAbort(t *testing.T) {
	t.Run("should return a onAbortFailError if onAbortFailed", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// onAbort will fail because the testAbort contract is not a valid contract
		abortAddress := sample.EthAddress()

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			sample.EthAddress().String(),
			big.NewInt(82),
			false,
			chainID,
			coin.CoinType_Gas,
			"",
			abortAddress,
			[]byte("foo"),
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrOnAbortFailed)

		// account still founded
		balance, err := k.BalanceOfZRC4(ctx, zrc20, abortAddress)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(82), balance)
	})

	t.Run("can't process abort for invalid chain ID", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// deploy test dapp
		testAbort := deployTestAbort(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			sample.EthAddress().String(),
			big.NewInt(82),
			false,
			919191,
			coin.CoinType_Gas,
			"",
			testAbort,
			[]byte("foo"),
		)

		// ASSERT
		require.Error(t, err)
		require.NotErrorIs(t, err, fungibletypes.ErrOnAbortFailed)
	})

	t.Run("should process NoAssetCall abort", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortMessage := []byte("abort message")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testAbort, err := k.DeployContract(ctx, testabort.TestAbortMetaData)
		require.NoError(t, err)
		require.NotEmpty(t, testAbort)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testAbort)
		abortAddress := testAbort

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_NoAssetCall,
			"",
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, abortAddress)
		require.NoError(t, err)
		require.Equal(t, int64(0), balance.Int64())
	})

	t.Run("should process Zeta abort", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := false
		abortMessage := []byte("abort message")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testAbort, err := k.DeployContract(ctx, testabort.TestAbortMetaData)
		require.NoError(t, err)
		require.NotEmpty(t, testAbort)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testAbort)
		abortAddress := testAbort

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Zeta,
			"",
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		balance := sdkk.BankKeeper.GetBalance(ctx, abortAddress.Bytes(), "azeta")
		require.Equal(t, amount.Int64(), balance.Amount.Int64())
	})

	t.Run("should process ERC20 abort", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortMessage := []byte("abort message")
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		testAbort, err := k.DeployContract(ctx, testabort.TestAbortMetaData)
		require.NoError(t, err)
		require.NotEmpty(t, testAbort)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testAbort)
		abortAddress := testAbort

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, abortAddress)
		require.NoError(t, err)
		require.Equal(t, amount.Int64(), balance.Int64())
	})

	t.Run("unable to process abort if the the universal contract is not abortable for coin-type ERC20", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortMessage := []byte("abort message")
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		abortAddress := sample.EthAddress()

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("unable to process abort if the the universal contract is not abortable for coin-type Gas", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortMessage := []byte("abort message")
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		abortAddress := sample.EthAddress()

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Gas,
			assetAddress,
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("unable to process abort if the the universal contract is not abortable for coin-type Zeta", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortMessage := []byte("abort message")
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		abortAddress := sample.EthAddress()

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Zeta,
			assetAddress,
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("should process Gas abort", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := false
		abortMessage := []byte("abort message")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		testAbort, err := k.DeployContract(ctx, testabort.TestAbortMetaData)
		require.NoError(t, err)
		require.NotEmpty(t, testAbort)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, testAbort)
		abortAddress := testAbort

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Gas,
			"",
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "", resp.VmError)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, abortAddress)
		require.NoError(t, err)
		require.Equal(t, amount.Int64(), balance.Int64())
	})

	t.Run("should fail if system contracts not deployed for Zeta", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortAddress := sample.EthAddress()
		abortMessage := []byte("abort message")

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Zeta,
			"",
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, fungibletypes.ErrSystemContractNotFound)
	})

	t.Run("should fail if gas coin not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortAddress := sample.EthAddress()
		abortMessage := []byte("abort message")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Gas,
			"",
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, crosschaintypes.ErrGasCoinNotFound)
	})

	t.Run("should fail if foreign coin not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortAddress := sample.EthAddress()
		abortMessage := []byte("abort message")
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, crosschaintypes.ErrForeignCoinNotFound)
	})

	t.Run("should fail if ZRC20 is paused", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortAddress := sample.EthAddress()
		abortMessage := []byte("abort message")

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		foreignCoin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		foreignCoin.Paused = true
		k.SetForeignCoins(ctx, foreignCoin)

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_Gas,
			"",
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, fungibletypes.ErrPausedZRC20)
	})

	t.Run("should fail if liquidity cap reached", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		inboundSender := sample.EthAddress().String()
		amount := big.NewInt(100)
		outgoing := true
		abortAddress := sample.EthAddress()
		abortMessage := []byte("abort message")
		assetAddress := sample.EthAddress().String()

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", assetAddress, "foobar")

		foreignCoin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		foreignCoin.LiquidityCap = math.NewUint(50)
		k.SetForeignCoins(ctx, foreignCoin)

		// ACT
		resp, err := k.ProcessAbort(
			ctx,
			inboundSender,
			amount,
			outgoing,
			chainID,
			coin.CoinType_ERC20,
			assetAddress,
			abortAddress,
			abortMessage,
		)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorIs(t, err, fungibletypes.ErrForeignCoinCapReached)
	})
}
