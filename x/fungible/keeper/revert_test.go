package keeper_test

import (
	"math/big"
	"testing"

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
		amount := big.NewInt(100)
		revertMessage := []byte("revert message")
		callOnRevert := true

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// deploy test dapp as revert address
		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress())
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
		require.Equal(t, "", resp.VmError, "VmError should be empty for successful execution")

		// verify the revert was processed by checking the test dapp state
		assertTestDAppV2MessageAndAmount(t, ctx, k, revertAddress, string(revertMessage), 0)
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

		// deploy the system contracts
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
			coin.CoinType_NoAssetCall,
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

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		resp, err := k.ProcessRevert(
			ctx,
			inboundSender,
			amount,
			chainID,
			coin.CoinType_Cmd, // Unsupported coin type
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
}
