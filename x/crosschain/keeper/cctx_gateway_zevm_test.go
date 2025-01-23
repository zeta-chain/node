package keeper_test

import (
	"math/big"
	"strings"
	"testing"

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

		// mock up ZETADepositAndCallContract
		fungibleMock.On("ZETADepositAndCallContract", mock.Anything,
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
		cctx.InboundParams.Status = types.InboundStatus_insufficient_depositor_fee

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

			// mock up ZETADepositAndCallContract
			fungibleMock.On("ZETADepositAndCallContract", mock.Anything,
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
			require.True(
				t,
				strings.Contains(
					cctx.CctxStatus.StatusMessage,
					"error during deposit that is not smart contract revert",
				),
			)
			require.True(t, strings.Contains(cctx.CctxStatus.ErrorMessage, sample.ErrSample.Error()))
		},
	)
}
