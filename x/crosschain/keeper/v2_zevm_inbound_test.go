package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
)

func TestKeeper_GetZetaInboundDetails(t *testing.T) {
	t.Run("success with valid parameters", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		receiverChainID := big.NewInt(1)
		callOptions := gatewayzevm.CallOptions{
			GasLimit:        big.NewInt(100000),
			IsArbitraryCall: false,
		}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		expectedChain := chains.Chain{
			ChainName:  chains.ChainName_eth_mainnet,
			ChainId:    1,
			IsExternal: true,
		}
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(1)).Return(expectedChain, true)
		observerMock.On("GetChainParamsByChainID", ctx, int64(1)).Return(&observertypes.ChainParams{
			ChainId:                  receiverChainID.Int64(),
			IsSupported:              true,
			ZetaTokenContractAddress: sample.EthAddress().Hex(),
		}, true)

		// ACT
		_, err := k.GetZETAInboundDetails(ctx, receiverChainID, callOptions)

		// ASSERT
		require.NoError(t, err)
	})

	t.Run("fail when receiver chain ID is nil", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		var receiverChainID *big.Int = nil
		callOptions := gatewayzevm.CallOptions{
			GasLimit:        big.NewInt(100000),
			IsArbitraryCall: false,
		}

		// ACT
		result, err := k.GetZETAInboundDetails(ctx, receiverChainID, callOptions)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidWithdrawalEvent)
		require.Contains(t, err.Error(), "receiver chain ID is nil or zero for ZETA withdrawal")
		require.Empty(t, result)
	})

	t.Run("fail when receiver chain ID is zero", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		receiverChainID := big.NewInt(0)
		callOptions := gatewayzevm.CallOptions{
			GasLimit:        big.NewInt(100000),
			IsArbitraryCall: false,
		}

		// ACT
		_, err := k.GetZETAInboundDetails(ctx, receiverChainID, callOptions)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidWithdrawalEvent)
		require.Contains(t, err.Error(), "receiver chain ID is nil or zero for ZETA withdrawal")
	})

	t.Run("fail when chain is not supported", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		receiverChainID := big.NewInt(999)
		callOptions := gatewayzevm.CallOptions{
			GasLimit:        big.NewInt(100000),
			IsArbitraryCall: false,
		}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(999)).Return(chains.Chain{}, false)

		// ACT
		_, err := k.GetZETAInboundDetails(ctx, receiverChainID, callOptions)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		require.Contains(t, err.Error(), "chain with chainID 999 not supported")
	})

	t.Run("fail when gas limit is nil", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		receiverChainID := big.NewInt(1)
		callOptions := gatewayzevm.CallOptions{
			GasLimit:        nil,
			IsArbitraryCall: false,
		}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		expectedChain := chains.Chain{
			ChainName: chains.ChainName_eth_mainnet,
			ChainId:   1,
		}
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(1)).Return(expectedChain, true)
		observerMock.On("GetChainParamsByChainID", ctx, int64(1)).Return(&observertypes.ChainParams{
			ChainId:                  receiverChainID.Int64(),
			IsSupported:              true,
			ZetaTokenContractAddress: sample.EthAddress().Hex(),
		}, true)

		// ACT
		_, err := k.GetZETAInboundDetails(ctx, receiverChainID, callOptions)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidWithdrawalEvent)
		require.Contains(t, err.Error(), "gas limit not provided for ZETA withdrawal")
	})

	t.Run("fail when gas limit is zero", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		receiverChainID := big.NewInt(1)
		callOptions := gatewayzevm.CallOptions{
			GasLimit:        big.NewInt(0),
			IsArbitraryCall: false,
		}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		expectedChain := chains.Chain{
			ChainName: chains.ChainName_eth_mainnet,
			ChainId:   1,
		}
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(1)).Return(expectedChain, true)
		observerMock.On("GetChainParamsByChainID", ctx, int64(1)).Return(&observertypes.ChainParams{
			ChainId:                  receiverChainID.Int64(),
			IsSupported:              true,
			ZetaTokenContractAddress: sample.EthAddress().Hex(),
		}, true)

		// ACT
		_, err := k.GetZETAInboundDetails(ctx, receiverChainID, callOptions)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidWithdrawalEvent)
		require.Contains(t, err.Error(), "gas limit not provided for ZETA withdrawal")
	})
}

func TestKeeper_GetErc20InboundDetails(t *testing.T) {
	t.Run("success with valid foreign coin", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		zrc20 := sample.EthAddress()
		callEvent := false

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		foreignCoin := fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().Hex(),
			Asset:                "USDT",
			ForeignChainId:       1,
			CoinType:             coin.CoinType_ERC20,
		}
		fungibleMock.On("GetForeignCoins", ctx, zrc20.Hex()).Return(foreignCoin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		expectedChain := chains.Chain{
			ChainName: chains.ChainName_eth_mainnet,
			ChainId:   1,
		}
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(1)).Return(expectedChain, true)

		gasLimit := big.NewInt(100000)
		fungibleMock.On("QueryGasLimit", ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress)).Return(gasLimit, nil)

		// ACT
		_, err := k.GetERC20InboundDetails(ctx, zrc20, callEvent)

		// ASSERT
		require.NoError(t, err)
	})

	t.Run("return empty result when foreign coin not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		zrc20 := sample.EthAddress()
		callEvent := false

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("GetForeignCoins", ctx, zrc20.Hex()).Return(fungibletypes.ForeignCoins{}, false)

		// ACT
		result, err := k.GetERC20InboundDetails(ctx, zrc20, callEvent)

		// ASSERT
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("fail when chain is not supported", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		zrc20 := sample.EthAddress()
		callEvent := false

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		foreignCoin := fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().Hex(),
			Asset:                "USDT",
			ForeignChainId:       999,
			CoinType:             coin.CoinType_ERC20,
		}
		fungibleMock.On("GetForeignCoins", ctx, zrc20.Hex()).Return(foreignCoin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(999)).Return(chains.Chain{}, false)

		// ACT
		_, err := k.GetERC20InboundDetails(ctx, zrc20, callEvent)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		require.Contains(t, err.Error(), "chain with chainID 999 not supported")
	})

	t.Run("fail when gas limit query fails", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		zrc20 := sample.EthAddress()
		callEvent := false

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		foreignCoin := fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().Hex(),
			Asset:                "USDT",
			ForeignChainId:       1,
			CoinType:             coin.CoinType_ERC20,
		}
		fungibleMock.On("GetForeignCoins", ctx, zrc20.Hex()).Return(foreignCoin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		expectedChain := chains.Chain{
			ChainName: chains.ChainName_eth_mainnet,
			ChainId:   1,
		}
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(1)).Return(expectedChain, true)

		fungibleMock.On("QueryGasLimit", ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress)).Return(nil, errors.New("gas limit query failed"))

		// ACT
		_, err := k.GetERC20InboundDetails(ctx, zrc20, callEvent)

		// ASSERT
		require.Error(t, err)
		require.Contains(t, err.Error(), "gas limit query failed")
	})

	t.Run("success with call event (NoAssetCall)", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		zrc20 := sample.EthAddress()
		callEvent := true

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		foreignCoin := fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().Hex(),
			Asset:                "USDT",
			ForeignChainId:       1,
			CoinType:             coin.CoinType_ERC20,
		}
		fungibleMock.On("GetForeignCoins", ctx, zrc20.Hex()).Return(foreignCoin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		expectedChain := chains.Chain{
			ChainName: chains.ChainName_eth_mainnet,
			ChainId:   1,
		}
		observerMock.On("GetSupportedChainFromChainID", ctx, int64(1)).Return(expectedChain, true)

		gasLimit := big.NewInt(100000)
		fungibleMock.On("QueryGasLimit", ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress)).Return(gasLimit, nil)

		// ACT
		_, err := k.GetERC20InboundDetails(ctx, zrc20, callEvent)

		// ASSERT
		require.NoError(t, err)
	})
}
