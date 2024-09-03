package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestCrossChainTx_GetEVMRevertAddress(t *testing.T) {
	t.Run("use revert address if revert options", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		addr := sample.EthAddress()
		cctx.RevertOptions.RevertAddress = addr.Hex()
		require.EqualValues(t, addr, cctx.GetEVMRevertAddress())
	})

	t.Run("use sender address if no revert options", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		addr := sample.EthAddress()
		cctx.InboundParams.Sender = addr.Hex()
		require.EqualValues(t, addr, cctx.GetEVMRevertAddress())
	})

}

func TestCrossChainTx_GetEVMAbortAddress(t *testing.T) {
	t.Run("use revert address if abort options", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		addr := sample.EthAddress()
		cctx.RevertOptions.AbortAddress = addr.Hex()
		require.EqualValues(t, addr, cctx.GetEVMAbortAddress())
	})

	t.Run("use sender address if no abort options", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		addr := sample.EthAddress()
		cctx.InboundParams.Sender = addr.Hex()
		require.EqualValues(t, addr, cctx.GetEVMAbortAddress())
	})
}

func TestCrossChainTx_SetOutboundBallot(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	ballotIndex := sample.ZetaIndex(t)
	cctx.SetOutboundBallotIndex(ballotIndex)
	require.Equal(t, ballotIndex, cctx.GetCurrentOutboundParam().BallotIndex)
}

func TestCrossChainTx_GetCCTXIndexBytes(t *testing.T) {
	cctx := sample.CrossChainTx(t, "sample")
	indexBytes, err := cctx.GetCCTXIndexBytes()
	require.NoError(t, err)
	require.Equal(t, cctx.Index, types.GetCctxIndexFromBytes(indexBytes))
}

func Test_NewCCTX(t *testing.T) {
	t.Run("should return a cctx with correct values", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		receiverChain := chains.Goerli
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteInbound{
			Creator:                 creator,
			Sender:                  sender.String(),
			SenderChainId:           senderChain.ChainId,
			Receiver:                receiver.String(),
			ReceiverChain:           receiverChain.ChainId,
			Amount:                  amount,
			Message:                 message,
			InboundHash:             inboundHash.String(),
			InboundBlockHeight:      inboundBlockHeight,
			GasLimit:                gasLimit,
			CoinType:                cointType,
			TxOrigin:                sender.String(),
			Asset:                   asset,
			EventIndex:              eventIndex,
			ProtocolContractVersion: types.ProtocolContractVersion_V2,
		}
		cctx, err := types.NewCCTX(ctx, msg, tss.TssPubkey)
		require.NoError(t, err)
		require.Equal(t, receiver.String(), cctx.GetCurrentOutboundParam().Receiver)
		require.Equal(t, receiverChain.ChainId, cctx.GetCurrentOutboundParam().ReceiverChainId)
		require.Equal(t, sender.String(), cctx.GetInboundParams().Sender)
		require.Equal(t, senderChain.ChainId, cctx.GetInboundParams().SenderChainId)
		require.Equal(t, amount, cctx.GetInboundParams().Amount)
		require.Equal(t, message, cctx.RelayedMessage)
		require.Equal(t, inboundHash.String(), cctx.GetInboundParams().ObservedHash)
		require.Equal(t, inboundBlockHeight, cctx.GetInboundParams().ObservedExternalHeight)
		require.Equal(t, gasLimit, cctx.GetCurrentOutboundParam().GasLimit)
		require.Equal(t, asset, cctx.GetInboundParams().Asset)
		require.Equal(t, cointType, cctx.InboundParams.CoinType)
		require.Equal(t, uint64(0), cctx.GetCurrentOutboundParam().TssNonce)
		require.Equal(t, sdkmath.ZeroUint(), cctx.GetCurrentOutboundParam().Amount)
		require.Equal(t, types.CctxStatus_PendingInbound, cctx.CctxStatus.Status)
		require.Equal(t, false, cctx.CctxStatus.IsAbortRefunded)
		require.Equal(t, types.ProtocolContractVersion_V2, cctx.ProtocolContractVersion)
	})

	t.Run("should return an error if the cctx is invalid", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		receiverChain := chains.Goerli
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             "",
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			GasLimit:           gasLimit,
			CoinType:           cointType,
			TxOrigin:           sender.String(),
			Asset:              asset,
			EventIndex:         eventIndex,
		}
		_, err := types.NewCCTX(ctx, msg, tss.TssPubkey)
		require.ErrorContains(t, err, "sender cannot be empty")
	})

	t.Run("zero value for protocol contract version gives V1", func(t *testing.T) {
		cctx := types.CrossChainTx{}
		require.Equal(t, types.ProtocolContractVersion_V1, cctx.ProtocolContractVersion)
	})
}

func TestCrossChainTx_Validate(t *testing.T) {
	cctx := sample.CrossChainTx(t, "foo")
	cctx.InboundParams = nil
	require.ErrorContains(t, cctx.Validate(), "inbound tx params cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundParams = nil
	require.ErrorContains(t, cctx.Validate(), "outbound tx params cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.CctxStatus = nil
	require.ErrorContains(t, cctx.Validate(), "cctx status cannot be nil")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.OutboundParams = make([]*types.OutboundParams, 3)
	require.ErrorContains(t, cctx.Validate(), "outbound tx params cannot be more than 2")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.Index = "0"
	require.ErrorContains(t, cctx.Validate(), "invalid index length 1")
	cctx = sample.CrossChainTx(t, "foo")
	cctx.InboundParams = sample.InboundParamsValidChainID(rand.New(rand.NewSource(42)))
}

func TestCrossChainTx_GetCurrentOutboundParam(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundParams = []*types.OutboundParams{}
	require.Equal(t, &types.OutboundParams{}, cctx.GetCurrentOutboundParam())

	cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(r)}
	require.Equal(t, cctx.OutboundParams[0], cctx.GetCurrentOutboundParam())

	cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(r), sample.OutboundParams(r)}
	require.Equal(t, cctx.OutboundParams[1], cctx.GetCurrentOutboundParam())
}

func TestCrossChainTx_IsCurrentOutboundRevert(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundParams = []*types.OutboundParams{}
	require.False(t, cctx.IsCurrentOutboundRevert())

	cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(r)}
	require.False(t, cctx.IsCurrentOutboundRevert())

	cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(r), sample.OutboundParams(r)}
	require.True(t, cctx.IsCurrentOutboundRevert())
}

func TestCrossChainTx_OriginalDestinationChainID(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	cctx := sample.CrossChainTx(t, "foo")

	cctx.OutboundParams = []*types.OutboundParams{}
	require.Equal(t, int64(-1), cctx.OriginalDestinationChainID())

	cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(r)}
	require.Equal(t, cctx.OutboundParams[0].ReceiverChainId, cctx.OriginalDestinationChainID())

	cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(r), sample.OutboundParams(r)}
	require.Equal(t, cctx.OutboundParams[0].ReceiverChainId, cctx.OriginalDestinationChainID())
}

func TestCrossChainTx_AddOutbound(t *testing.T) {
	t.Run("successfully get outbound tx", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.AddOutbound(ctx, types.MsgVoteOutbound{
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundHash:              hash,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundGasUsed:           100,
			ObservedOutboundEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutboundEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Hash, hash)
		require.Equal(t, cctx.GetCurrentOutboundParam().GasUsed, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasPrice, sdkmath.NewInt(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasLimit, uint64(20))
		require.Equal(t, cctx.GetCurrentOutboundParam().ObservedExternalHeight, uint64(10))
	})

	t.Run("successfully get outbound tx for failed ballot without amount check", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.AddOutbound(ctx, types.MsgVoteOutbound{
			ObservedOutboundHash:              hash,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundGasUsed:           100,
			ObservedOutboundEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutboundEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_FailureObservation)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Hash, hash)
		require.Equal(t, cctx.GetCurrentOutboundParam().GasUsed, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasPrice, sdkmath.NewInt(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasLimit, uint64(20))
		require.Equal(t, cctx.GetCurrentOutboundParam().ObservedExternalHeight, uint64(10))
	})

	t.Run("failed to get outbound tx if amount does not match value received", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)

		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.AddOutbound(ctx, types.MsgVoteOutbound{
			ValueReceived:                     sdkmath.NewUint(100),
			ObservedOutboundHash:              hash,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundGasUsed:           100,
			ObservedOutboundEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutboundEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})
}

func Test_SetRevertOutboundValues(t *testing.T) {
	t.Run("successfully set revert outbound values", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.OutboundParams = cctx.OutboundParams[:1]
		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.Equal(t, cctx.GetCurrentOutboundParam().Receiver, cctx.InboundParams.Sender)
		require.Equal(t, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.InboundParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount, cctx.OutboundParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutboundParam().GasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().TssPubkey, cctx.OutboundParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
	})

	t.Run("failed to set revert outbound values if revert outbound already exists", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		err := cctx.AddRevertOutbound(100)
		require.ErrorContains(t, err, "cannot revert a revert tx")
	})

	t.Run("failed to set revert outbound values if revert outbound already exists", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.OutboundParams = make([]*types.OutboundParams, 0)
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

func TestCrossChainTx_SetOutboundMined(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetOutboundMined("test")
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
