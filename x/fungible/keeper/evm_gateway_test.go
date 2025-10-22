package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"
)

func TestKeeper_DepositAndCallZeta(t *testing.T) {
	t.Run("DepositAndCallZeta successfully", func(t *testing.T) {
		// ZETA v2 not enabled
		// TODO: enable back
		// https://github.com/zeta-chain/node/issues/4373
		t.Skip()
		
		// Arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		chainID := chains.Ethereum.ChainId
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		testMessage := "test message"
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		err := k.MintZetaToFungibleModule(ctx, amount)
		require.NoError(t, err)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		_ = deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)

		// Act
		res, err := k.DepositAndCallZeta(
			ctx,
			gatewayzevm.MessageContext{
				Sender:    inboundSender.Bytes(),
				SenderEVM: inboundSender,
				ChainID:   big.NewInt(chainID),
			},
			amount,
			testDAppV2,
			[]byte(testMessage),
		)
		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, res)
		assertTestDAppV2MessageAndAmount(t, ctx, k, testDAppV2, testMessage, amount.Int64())
	})

	t.Run("DepositAndCallZeta fails if system contract is not present", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		message := []byte("test message")
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)

		// Act
		_, err = k.DepositAndCallZeta(
			ctx,
			gatewayzevm.MessageContext{
				Sender:    inboundSender.Bytes(),
				SenderEVM: inboundSender,
				ChainID:   big.NewInt(chains.Ethereum.ChainId),
			},
			amount,
			testDAppV2,
			message,
		)

		// Assert
		require.ErrorIs(t, types.ErrSystemContractNotFound, err)
	})

	t.Run("DepositAndCallZeta fails if gateway is not set", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		message := []byte("test message")
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		k.SetSystemContract(ctx, types.SystemContract{
			ConnectorZevm:  sample.EthAddress().Hex(),
			SystemContract: sample.EthAddress().Hex(),
		})

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)

		// Act
		_, err = k.DepositAndCallZeta(
			ctx,
			gatewayzevm.MessageContext{
				Sender:    inboundSender.Bytes(),
				SenderEVM: inboundSender,
				ChainID:   big.NewInt(chains.Ethereum.ChainId),
			},
			amount,
			testDAppV2,
			message,
		)

		// Assert
		require.ErrorIs(t, types.ErrGatewayContractNotSet, err)
	})
}

func TestKeeper_CallDepositAndRevert(t *testing.T) {
	t.Run("CallDepositAndRevert successfully", func(t *testing.T) {
		// Arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		chainID := chains.Ethereum.ChainId
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		message := []byte("test message")
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)

		// Act
		res, err := k.CallDepositAndRevert(
			ctx,
			inboundSender.Hex(),
			zrc20,
			amount,
			testDAppV2,
			message,
		)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, res)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, testDAppV2)
		require.NoError(t, err)
		require.Equal(t, amount, balance)

		// Check message was set in TestDAppV2
		testDAppABI, err := testdappv2.TestDAppV2MetaData.GetAbi()
		require.NoError(t, err)

		messageCalledData, err := testDAppABI.Pack("getCalledWithMessage", string(message))
		require.NoError(t, err)

		messageResponse, err := k.CallEVMWithData(ctx, types.ModuleAddressEVM, &testDAppV2, messageCalledData, true, true, big.NewInt(0), nil)
		require.NoError(t, err)
		require.False(t, messageResponse.Failed())

		messageResult, err := testDAppABI.Unpack("getCalledWithMessage", messageResponse.Ret)
		require.NoError(t, err)
		wasCalled := messageResult[0].(bool)

		require.True(t, wasCalled, "TestDAppV2 should have been called with message: %s", string(message))
	})

	t.Run("CallDepositAndRevert fails if system contract is not present", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		message := []byte("test message")
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)

		// Act
		_, err = k.CallDepositAndRevert(
			ctx,
			inboundSender.Hex(),
			sample.EthAddress(),
			amount,
			testDAppV2,
			message,
		)

		// Assert
		require.ErrorIs(t, types.ErrSystemContractNotFound, err)
	})

	t.Run("CallDepositAndRevert fails if gateway is not set", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		message := []byte("test message")
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		k.SetSystemContract(ctx, types.SystemContract{
			ConnectorZevm:  sample.EthAddress().Hex(),
			SystemContract: sample.EthAddress().Hex(),
		})

		testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, types.ModuleAddressEVM)
		require.NoError(t, err)
		require.NotEmpty(t, testDAppV2)

		// Act
		_, err = k.CallDepositAndRevert(
			ctx,
			inboundSender.Hex(),
			sample.EthAddress(),
			amount,
			testDAppV2,
			message,
		)

		// Assert
		require.ErrorIs(t, types.ErrGatewayContractNotSet, err)
	})

	t.Run("CallDepositAndRevert fails if target contract is not present", func(t *testing.T) {
		// Arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		chainID := chains.Ethereum.ChainId
		inboundSender := sample.EthAddress()
		amount := big.NewInt(1000)
		message := []byte("test message")
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")

		// Act
		_, err := k.CallDepositAndRevert(
			ctx,
			inboundSender.Hex(),
			zrc20,
			amount,
			sample.EthAddress(),
			message,
		)

		// Assert
		require.ErrorContains(t, err, "contract_call_error")
	})
}
