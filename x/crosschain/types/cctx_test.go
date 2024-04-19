package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestCrossChainTx_GetCCTXIndexBytes(t *testing.T) {
	cctx := sample.CrossChainTx(t, "sample")
	indexBytes, err := cctx.GetCCTXIndexBytes()
	require.NoError(t, err)
	require.Equal(t, cctx.Index, types.GetCctxIndexFromBytes(indexBytes))
}

func Test_InitializeCCTX(t *testing.T) {
	t.Run("should return a cctx with correct values", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		senderChain := chains.GoerliChain()
		sender := sample.EthAddress()
		receiverChain := chains.GoerliChain()
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		intxBlockHeight := uint64(420)
		intxHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteOnObservedInboundTx{
			Creator:       creator,
			Sender:        sender.String(),
			SenderChainId: senderChain.ChainId,
			Receiver:      receiver.String(),
			ReceiverChain: receiverChain.ChainId,
			Amount:        amount,
			Message:       message,
			InTxHash:      intxHash.String(),
			InBlockHeight: intxBlockHeight,
			GasLimit:      gasLimit,
			CoinType:      cointType,
			TxOrigin:      sender.String(),
			Asset:         asset,
			EventIndex:    eventIndex,
		}
		cctx, err := types.NewCCTX(ctx, msg, tss.TssPubkey)
		require.NoError(t, err)
		require.Equal(t, receiver.String(), cctx.GetCurrentOutTxParam().Receiver)
		require.Equal(t, receiverChain.ChainId, cctx.GetCurrentOutTxParam().ReceiverChainId)
		require.Equal(t, sender.String(), cctx.GetInboundTxParams().Sender)
		require.Equal(t, senderChain.ChainId, cctx.GetInboundTxParams().SenderChainId)
		require.Equal(t, amount, cctx.GetInboundTxParams().Amount)
		require.Equal(t, message, cctx.RelayedMessage)
		require.Equal(t, intxHash.String(), cctx.GetInboundTxParams().InboundTxObservedHash)
		require.Equal(t, intxBlockHeight, cctx.GetInboundTxParams().InboundTxObservedExternalHeight)
		require.Equal(t, gasLimit, cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
		require.Equal(t, asset, cctx.GetInboundTxParams().Asset)
		require.Equal(t, cointType, cctx.InboundTxParams.CoinType)
		require.Equal(t, uint64(0), cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
		require.Equal(t, sdkmath.ZeroUint(), cctx.GetCurrentOutTxParam().Amount)
		require.Equal(t, types.CctxStatus_PendingInbound, cctx.CctxStatus.Status)
		require.Equal(t, false, cctx.CctxStatus.IsAbortRefunded)
	})
	t.Run("should return an error if the cctx is invalid", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		senderChain := chains.GoerliChain()
		sender := sample.EthAddress()
		receiverChain := chains.GoerliChain()
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		intxBlockHeight := uint64(420)
		intxHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteOnObservedInboundTx{
			Creator:       creator,
			Sender:        "invalid",
			SenderChainId: senderChain.ChainId,
			Receiver:      receiver.String(),
			ReceiverChain: receiverChain.ChainId,
			Amount:        amount,
			Message:       message,
			InTxHash:      intxHash.String(),
			InBlockHeight: intxBlockHeight,
			GasLimit:      gasLimit,
			CoinType:      cointType,
			TxOrigin:      sender.String(),
			Asset:         asset,
			EventIndex:    eventIndex,
		}
		_, err := types.NewCCTX(ctx, msg, tss.TssPubkey)
		require.ErrorContains(t, err, "invalid address")
	})
}

func TestCrossChainTx_Validate(t *testing.T) {
	cctx := sample.CrossChainTx(t, "foo")
	cctx.InboundTxParams = nil
	require.ErrorContains(t, cctx.Validate(), "inbound tx params cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundTxParams = nil
	require.ErrorContains(t, cctx.Validate(), "outbound tx params cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.CctxStatus = nil
	require.ErrorContains(t, cctx.Validate(), "cctx status cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundTxParams = make([]*types.OutboundTxParams, 3)
	require.ErrorContains(t, cctx.Validate(), "outbound tx params cannot be more than 2")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.Index = "0"
	require.ErrorContains(t, cctx.Validate(), "invalid index length 1")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.InboundTxParams = sample.InboundTxParamsValidChainID(rand.New(rand.NewSource(42)))
	cctx.InboundTxParams.SenderChainId = 1000
	require.ErrorContains(t, cctx.Validate(), "invalid sender chain id 1000")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParamsValidChainID(rand.New(rand.NewSource(42)))}
	cctx.InboundTxParams = sample.InboundTxParamsValidChainID(rand.New(rand.NewSource(42)))
	cctx.InboundTxParams.InboundTxObservedHash = sample.Hash().String()
	cctx.InboundTxParams.InboundTxBallotIndex = sample.ZetaIndex(t)
	cctx.OutboundTxParams[0].ReceiverChainId = 1000
	require.ErrorContains(t, cctx.Validate(), "invalid receiver chain id 1000")
}

func TestCrossChainTx_GetCurrentOutTxParam(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundTxParams = []*types.OutboundTxParams{}
	require.Equal(t, &types.OutboundTxParams{}, cctx.GetCurrentOutTxParam())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[0], cctx.GetCurrentOutTxParam())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r), sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[1], cctx.GetCurrentOutTxParam())
}

func TestCrossChainTx_IsCurrentOutTxRevert(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundTxParams = []*types.OutboundTxParams{}
	require.False(t, cctx.IsCurrentOutTxRevert())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r)}
	require.False(t, cctx.IsCurrentOutTxRevert())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r), sample.OutboundTxParams(r)}
	require.True(t, cctx.IsCurrentOutTxRevert())
}

func TestCrossChainTx_OriginalDestinationChainID(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundTxParams = []*types.OutboundTxParams{}
	require.Equal(t, int64(-1), cctx.OriginalDestinationChainID())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[0].ReceiverChainId, cctx.OriginalDestinationChainID())

	cctx.OutboundTxParams = []*types.OutboundTxParams{sample.OutboundTxParams(r), sample.OutboundTxParams(r)}
	require.Equal(t, cctx.OutboundTxParams[0].ReceiverChainId, cctx.OriginalDestinationChainID())
}

func TestCrossChainTx_AddOutbound(t *testing.T) {
	t.Run("successfully get outbound tx", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.AddOutbound(ctx, types.MsgVoteOnObservedOutboundTx{
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

		err := cctx.AddOutbound(ctx, types.MsgVoteOnObservedOutboundTx{
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

		err := cctx.AddOutbound(ctx, types.MsgVoteOnObservedOutboundTx{
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

func Test_SetRevertOutboundValues(t *testing.T) {
	t.Run("successfully set revert outbound values", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.OutboundTxParams = cctx.OutboundTxParams[:1]
		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundTxParams, 2)
		require.Equal(t, cctx.GetCurrentOutTxParam().Receiver, cctx.InboundTxParams.Sender)
		require.Equal(t, cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.InboundTxParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount, cctx.OutboundTxParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxGasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutTxParam().TssPubkey, cctx.OutboundTxParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundTxParams[0].TxFinalizationStatus)
	})

	t.Run("failed to set revert outbound values if revert outbound already exists", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		err := cctx.AddRevertOutbound(100)
		require.ErrorContains(t, err, "cannot revert a revert tx")
	})

	t.Run("failed to set revert outbound values if revert outbound already exists", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.OutboundTxParams = make([]*types.OutboundTxParams, 0)
		err := cctx.AddRevertOutbound(100)
		require.ErrorContains(t, err, "cannot revert before trying to process an outbound tx")
	})
}

func TestCrossChainTx_SetAbort(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.SetAbort("test")
	require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	require.Equal(t, "test", "test")
}

func TestCrossChainTx_SetPendingRevert(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetPendingRevert("test")
	require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
}

func TestCrossChainTx_SetPendingOutbound(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
	cctx.SetPendingOutbound("test")
	require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
}

func TestCrossChainTx_SetOutBoundMined(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetOutBoundMined("test")
	require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
}

func TestCrossChainTx_SetReverted(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
	cctx.SetReverted("test")
	require.Equal(t, types.CctxStatus_Reverted, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
}
