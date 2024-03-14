package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetOutbound(t *testing.T) {
	t.Run("successfully get outbound tx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		tss := sample.Tss()
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		zk.ObserverKeeper.SetTSS(ctx, tss)
		err := k.GetOutbound(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxHash:              hash,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxGasUsed:           100,
			ObservedOutTxEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutTxEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxHash, hash)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxGasUsed, uint64(100))
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasPrice, sdkmath.NewInt(100))
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasLimit, uint64(20))
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight, uint64(10))
		require.Equal(t, cctx.CctxStatus.LastUpdateTimestamp, ctx.BlockHeader().Time.Unix())
	})

	t.Run("successfully get outbound tx for failed ballot without amount check", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		tss := sample.Tss()
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		zk.ObserverKeeper.SetTSS(ctx, tss)
		err := k.GetOutbound(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
			ObservedOutTxHash:              hash,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxGasUsed:           100,
			ObservedOutTxEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutTxEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_FailureObservation)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxHash, hash)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxGasUsed, uint64(100))
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasPrice, sdkmath.NewInt(100))
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasLimit, uint64(20))
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight, uint64(10))
		require.Equal(t, cctx.CctxStatus.LastUpdateTimestamp, ctx.BlockHeader().Time.Unix())
	})

	t.Run("failed to get outbound tx if amount does not match value received", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		tss := sample.Tss()
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		zk.ObserverKeeper.SetTSS(ctx, tss)
		err := k.GetOutbound(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
			ValueReceived:                  sdkmath.NewUint(100),
			ObservedOutTxHash:              hash,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxGasUsed:           100,
			ObservedOutTxEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutTxEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("failed to get outbound tx if tss mismatch", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		tss := sample.Tss()
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		zk.ObserverKeeper.SetTSS(ctx, tss)
		err := k.GetOutbound(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxHash:              hash,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxGasUsed:           100,
			ObservedOutTxEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutTxEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.ErrorIs(t, err, types.ErrTssMismatch)
	})

	t.Run("failed to get outbound tx if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		tss := sample.Tss()
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		err := k.GetOutbound(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxHash:              hash,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxGasUsed:           100,
			ObservedOutTxEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutTxEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})
}

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
		cctx.CoinType = common.CoinType_Cmd
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
	})

	t.Run("successfully process failed outbound set to aborted for withdraw tx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundTxParams.SenderChainId = common.ZetaChainMainnet().ChainId
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, *senderChain)

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessFailedOutbound(ctx, cctx, sample.String())
		require.NoError(t, err)
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

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
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_SuccessObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
	})

	t.Run("successfully process outbound with ballot finalized to failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.CoinType = common.CoinType_Cmd
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
	})

	t.Run("do not process outbound on error", func(t *testing.T) {
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

		// mock failed GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, false).Once()

		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingOutbound)
	})
}

func TestKeeper_SaveFailedOutBound(t *testing.T) {
	t.Run("successfully save failed outbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    "",
			ChainId:  cctx.GetCurrentOutTxParam().ReceiverChainId,
			Nonce:    cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			HashList: nil,
		})
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SaveFailedOutBound(ctx, cctx, sample.String())
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
		_, found := k.GetOutTxTracker(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
		require.False(t, found)
	})
}

func TestKeeper_SaveSuccessfulOutBound(t *testing.T) {
	t.Run("successfully save successful outbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    "",
			ChainId:  cctx.GetCurrentOutTxParam().ReceiverChainId,
			Nonce:    cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			HashList: nil,
		})
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SaveSuccessfulOutBound(ctx, cctx, sample.String())
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxBallotIndex, sample.String())
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
		_, found := k.GetOutTxTracker(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
		require.False(t, found)
	})
}
