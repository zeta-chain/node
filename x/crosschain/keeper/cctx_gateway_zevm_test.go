package keeper_test

import (
	"math/big"
	"strings"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestKeeper_InitiateOutboundZEVM(t *testing.T) {
	t.Run("initiate outbound from gateway ZEVM successfully", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		gatewayZEVM := keeper.NewCCTXGatewayZEVM(*k)

		// setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// mock up LegacyZETADepositAndCallContract
		fungibleMock.On("LegacyZETADepositAndCallContract", mock.Anything,
			mock.Anything,
			receiver, int64(1), amount, mock.Anything, mock.Anything).Return(nil, nil)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
		cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.GetInboundParams().SenderChainId = 1

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayZEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_OutboundMined, newStatus)
	})

	t.Run("should return aborted status on insufficient depositor fee", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		gatewayZEVM := keeper.NewCCTXGatewayZEVM(*k)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.Status = types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE
		cctx.InboundParams.ErrorMessage = "deposited amount is less than depositor fee"

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayZEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, "deposited amount is less than depositor fee")
	})

	t.Run("should return aborted status on unknown inbound status", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		gatewayZEVM := keeper.NewCCTXGatewayZEVM(*k)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.Status = types.InboundStatus(1000)

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayZEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
	})

	t.Run("should return reverted status on invalid memo inbound status", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})
		gatewayZEVM := keeper.NewCCTXGatewayZEVM(*k)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     1,
			MedianIndex: 0,
			Prices:      []uint64{1},
		})

		//mock necessary calls made during creation of the revert outbound
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(sample.Chain(1), true)
		fungibleMock.On("GetGasCoinForForeignCoin", mock.Anything, mock.Anything).
			Return(fungibletypes.ForeignCoins{Zrc20ContractAddress: "0x"}, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(ethcommon.Address{}, nil)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(ethcommon.Address{}, nil)
		fungibleMock.On("QueryProtocolFlatFee", mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).
			Return(observertypes.ChainNonces{Nonce: 1}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).
			Return(observertypes.PendingNonces{NonceHigh: int64(1)}, true)
		observerMock.On("SetChainNonces", mock.Anything, mock.Anything)
		observerMock.On("SetPendingNonces", mock.Anything, mock.Anything)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.SenderChainId = 1
		cctx.InboundParams.Status = types.InboundStatus_INVALID_MEMO
		cctx.InboundParams.ErrorMessage = "invalid revert address in memo"
		cctx.OutboundParams = []*types.OutboundParams{cctx.OutboundParams[0]}

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayZEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_PendingRevert, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, "invalid revert address in memo")
	})

	t.Run("should return reverted status on excessive NoAssetCall funds status", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})
		gatewayZEVM := keeper.NewCCTXGatewayZEVM(*k)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     1,
			MedianIndex: 0,
			Prices:      []uint64{1},
		})

		//mock necessary calls made during creation of the revert outbound
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(sample.Chain(1), true)
		fungibleMock.On("GetGasCoinForForeignCoin", mock.Anything, mock.Anything).
			Return(fungibletypes.ForeignCoins{Zrc20ContractAddress: "0x"}, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(ethcommon.Address{}, nil)
		fungibleMock.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
			Return(ethcommon.Address{}, nil)
		fungibleMock.On("QueryProtocolFlatFee", mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).
			Return(observertypes.ChainNonces{Nonce: 1}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).
			Return(observertypes.PendingNonces{NonceHigh: int64(1)}, true)
		observerMock.On("SetChainNonces", mock.Anything, mock.Anything)
		observerMock.On("SetPendingNonces", mock.Anything, mock.Anything)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.SenderChainId = 1
		cctx.InboundParams.Status = types.InboundStatus_EXCESSIVE_NOASSETCALL_FUNDS
		cctx.InboundParams.ErrorMessage = "remaining funds of 100001 satoshis exceed 100000 satoshis"
		cctx.OutboundParams = []*types.OutboundParams{cctx.OutboundParams[0]}

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayZEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_PendingRevert, newStatus)
		require.Contains(t, cctx.CctxStatus.ErrorMessage, "remaining funds of 100001 satoshis exceed 100000 satoshis")
	})

	t.Run(
		"should return aborted status on 'error during deposit that is not smart contract revert'",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})
			gatewayZEVM := keeper.NewCCTXGatewayZEVM(*k)

			// setup mock data
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)

			// mock up LegacyZETADepositAndCallContract
			fungibleMock.On("LegacyZETADepositAndCallContract", mock.Anything,
				mock.Anything,
				receiver, int64(1), amount, mock.Anything, mock.Anything).Return(nil, sample.ErrSample)

			// mock up CCTX data
			cctx := sample.CrossChainTx(t, "test")
			cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
			cctx.GetCurrentOutboundParam().Receiver = receiver.String()
			cctx.GetCurrentOutboundParam().ReceiverChainId = chains.ZetaChainPrivnet.ChainId
			cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
			cctx.InboundParams.CoinType = coin.CoinType_Zeta
			cctx.GetInboundParams().SenderChainId = 1

			// ACT
			// call InitiateOutbound
			newStatus, err := gatewayZEVM.InitiateOutbound(
				ctx,
				keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
			)

			// ASSERT
			require.Error(t, err)
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Equal(t, types.CctxStatus_Aborted, newStatus)
			require.Contains(
				t,
				cctx.CctxStatus.StatusMessage,
				"outbound failed but the universal contract did not revert",
			)
			require.True(t, strings.Contains(cctx.CctxStatus.ErrorMessage, sample.ErrSample.Error()))
		},
	)
}
