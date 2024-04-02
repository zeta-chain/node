package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_ProcessSuccessfulOutbound(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	cctx := sample.CrossChainTx(t, "test")
	// transition to reverted if pending revert
	cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
	k.ProcessSuccessfulOutbound(ctx, cctx, sample.String())
	require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Reverted)
	// transition to outbound mined if pending outbound
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	k.ProcessSuccessfulOutbound(ctx, cctx, sample.String())
	require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
	// do nothing if it's in any other state
	k.ProcessSuccessfulOutbound(ctx, cctx, sample.String())
	require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
}

func TestKeeper_ProcessFailedOutbound(t *testing.T) {
	t.Run("successfully process failed outbound set to aborted for type cmd", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundTxParams.CoinType = coin.CoinType_Cmd
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("successfully process failed outbound set to aborted for withdraw tx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundTxParams.SenderChainId = chains.ZetaChainMainnet().ChainId
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("successfully process failed outbound set to pending revert", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, *senderChain)

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingRevert)
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, cctx.GetCurrentOutTxParam().TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundTxParams[0].TxFinalizationStatus)
	})

	t.Run("successfully process failed outbound set to pending revert if gas limit is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 0)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, *senderChain)

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingRevert)
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, cctx.GetCurrentOutTxParam().TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundTxParams[0].TxFinalizationStatus)
	})

	t.Run("unable to process revert when update nonce fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// mock failed UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.ErrorIs(t, err, types.ErrCannotFindReceiverNonce)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingOutbound)
	})

	t.Run("unable to process revert when PayGasAndUpdateCctx fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingOutbound)
	})

	t.Run("unable to process revert when GetRevertGasLimit fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// mock failed GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, false).Once()

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingOutbound)
	})
}

func TestKeeper_ProcessOutbound(t *testing.T) {
	t.Run("successfully process outbound with ballot finalized to success", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.GoerliChain(), "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_SuccessObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
	})

	t.Run("successfully process outbound with ballot finalized to failed and old status is Pending Revert", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.GoerliChain(), "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("successfully process outbound with ballot finalized to failed and coin-type is CMD", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.GoerliChain(), "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundTxParams.CoinType = coin.CoinType_Cmd
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("do not process if cctx invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.GoerliChain(), "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundTxParams = nil
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.Error(t, err)
	})

	t.Run("do not process outbound on error, no new outbound created", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		oldOutTxParamsLen := len(cctx.OutboundTxParams)
		// mock failed GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, false).Once()

		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingOutbound)
		// New outbound not added and the old outbound is not finalized
		require.Len(t, cctx.OutboundTxParams, oldOutTxParamsLen)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_NotFinalized)
	})

	t.Run("do not process outbound if the cctx has already been reverted once", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.OutboundTxParams = append(cctx.OutboundTxParams, sample.OutboundTxParams(sample.Rand()))
		cctx.OutboundTxParams[1].ReceiverChainId = 5
		cctx.OutboundTxParams[1].OutboundTxBallotIndex = ""
		cctx.OutboundTxParams[1].OutboundTxHash = ""

		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.Error(t, err)
	})

	t.Run("successfully revert a outbound and create a new revert tx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		oldOutTxParamsLen := len(cctx.OutboundTxParams)
		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, *senderChain)

		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingRevert)
		// New outbound added for revert and the old outbound is finalized
		require.Len(t, cctx.OutboundTxParams, oldOutTxParamsLen+1)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_NotFinalized)
		require.Equal(t, cctx.OutboundTxParams[oldOutTxParamsLen-1].TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})
}
