package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_GetOutbound(t *testing.T) {
	t.Run("successfully get outbound tx", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := keeper.SetOutboundValues(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
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
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := keeper.SetOutboundValues(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
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
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)

		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := keeper.SetOutboundValues(ctx, cctx, types.MsgVoteOnObservedOutboundTx{
			ValueReceived:                  sdkmath.NewUint(100),
			ObservedOutTxHash:              hash,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxGasUsed:           100,
			ObservedOutTxEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutTxEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
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
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("successfully process failed outbound set to aborted for withdraw tx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundTxParams.SenderChainId = common.ZetaChainMainnet().ChainId
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, *senderChain)

		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		// Remove the first outbound tx param to make the scenario realistic
		oldParams := cctx.OutboundTxParams
		cctx.OutboundTxParams = make([]*types.OutboundTxParams, 1)
		cctx.OutboundTxParams[0] = oldParams[1]
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

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

	t.Run("successfully process outbound with ballot finalized to failed and old status is Pending Revert", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("successfully process outbound with ballot finalized to failed and coin-type is CMD", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.CoinType = common.CoinType_Cmd
		err := k.ProcessOutbound(ctx, cctx, observertypes.BallotStatus_BallotFinalized_FailureObservation, sample.String())
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutTxParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
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
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

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
		k.SaveFailedOutBound(ctx, cctx, sample.String(), sample.ZetaIndex(t))
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
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
		_, found := k.GetOutTxTracker(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
		require.False(t, found)
	})
}

func TestKeeper_SaveOutbound(t *testing.T) {
	t.Run("successfully save outbound", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// setup state for crosschain and observer modules
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		ballotIndex := sample.String()
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    "",
			ChainId:  cctx.GetCurrentOutTxParam().ReceiverChainId,
			Nonce:    cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			HashList: nil,
		})

		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce) - 1,
			NonceHigh: int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce) + 1,
			ChainId:   cctx.GetCurrentOutTxParam().ReceiverChainId,
			Tss:       cctx.GetCurrentOutTxParam().TssPubkey,
		})
		zk.ObserverKeeper.SetTSS(ctx, observertypes.TSS{
			TssPubkey: cctx.GetCurrentOutTxParam().TssPubkey,
		})

		// Save outbound and assert all values are successfully saved
		k.SaveOutbound(ctx, cctx, ballotIndex)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxBallotIndex, ballotIndex)
		_, found := k.GetOutTxTracker(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
		require.False(t, found)
		pn, found := zk.ObserverKeeper.GetPendingNonces(ctx, cctx.GetCurrentOutTxParam().TssPubkey, cctx.GetCurrentOutTxParam().ReceiverChainId)
		require.True(t, found)
		require.Equal(t, pn.NonceLow, int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce)+1)
		require.Equal(t, pn.NonceHigh, int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce)+1)
		_, found = k.GetInTxHashToCctx(ctx, cctx.InboundTxParams.InboundTxObservedHash)
		require.True(t, found)
		_, found = zk.ObserverKeeper.GetNonceToCctx(ctx, cctx.GetCurrentOutTxParam().TssPubkey, cctx.GetCurrentOutTxParam().ReceiverChainId, int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce))
		require.True(t, found)
	})
}

func Test_SetRevertOutboundValues(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.OutboundTxParams = cctx.OutboundTxParams[:1]
	keeper.SetRevertOutboundValues(cctx, 100)
	require.Len(t, cctx.OutboundTxParams, 2)
	require.Equal(t, cctx.GetCurrentOutTxParam().Receiver, cctx.InboundTxParams.Sender)
	require.Equal(t, cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.InboundTxParams.SenderChainId)
	require.Equal(t, cctx.GetCurrentOutTxParam().Amount, cctx.InboundTxParams.Amount)
	require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxGasLimit, uint64(100))
	require.Equal(t, cctx.GetCurrentOutTxParam().TssPubkey, cctx.OutboundTxParams[0].TssPubkey)
	require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundTxParams[0].TxFinalizationStatus)
}

func TestKeeper_ValidateOutboundMessage(t *testing.T) {
	t.Run("successfully validate outbound message", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
		_, err := k.ValidateOutboundMessage(ctx, types.MsgVoteOnObservedOutboundTx{
			CctxHash:      cctx.Index,
			OutTxTssNonce: cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			OutTxChain:    cctx.GetCurrentOutTxParam().ReceiverChainId,
		})
		require.NoError(t, err)
	})

	t.Run("failed to validate outbound message if cctx not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		msg := types.MsgVoteOnObservedOutboundTx{
			CctxHash:      sample.String(),
			OutTxTssNonce: 1,
		}
		_, err := k.ValidateOutboundMessage(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		require.ErrorContains(t, err, fmt.Sprintf("CCTX %s does not exist", msg.CctxHash))
	})

	t.Run("failed to validate outbound message if nonce does not match", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		msg := types.MsgVoteOnObservedOutboundTx{
			CctxHash:      cctx.Index,
			OutTxTssNonce: 2,
		}
		_, err := k.ValidateOutboundMessage(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		require.ErrorContains(t, err, fmt.Sprintf("OutTxTssNonce %d does not match CCTX OutTxTssNonce %d", msg.OutTxTssNonce, cctx.GetCurrentOutTxParam().OutboundTxTssNonce))
	})

	t.Run("failed to validate outbound message if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		_, err := k.ValidateOutboundMessage(ctx, types.MsgVoteOnObservedOutboundTx{
			CctxHash:      cctx.Index,
			OutTxTssNonce: cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
		})
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})

	t.Run("failed to validate outbound message if chain does not match", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
		_, err := k.ValidateOutboundMessage(ctx, types.MsgVoteOnObservedOutboundTx{
			CctxHash:      cctx.Index,
			OutTxTssNonce: cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			OutTxChain:    2,
		})
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		require.ErrorContains(t, err, fmt.Sprintf("OutTxChain %d does not match CCTX OutTxChain %d", 2, cctx.GetCurrentOutTxParam().ReceiverChainId))
	})
}
