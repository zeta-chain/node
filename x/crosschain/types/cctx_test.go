package types_test

import (
	"math/rand"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
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
	require.Equal(t, &types.OutboundParams{CallOptions: &types.CallOptions{}}, cctx.GetCurrentOutboundParam())

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
		require.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().TssPubkey, cctx.OutboundParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
	})

	t.Run("successfully set BTC revert address V1", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.SenderChainId = chains.BitcoinTestnet.ChainId
		cctx.OutboundParams = cctx.OutboundParams[:1]
		cctx.RevertOptions.RevertAddress = sample.BtcAddressP2WPKH(t, &chaincfg.TestNet3Params).String()

		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.Equal(t, cctx.GetCurrentOutboundParam().Receiver, cctx.RevertOptions.RevertAddress)
		require.Equal(t, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.InboundParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount, cctx.OutboundParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().TssPubkey, cctx.OutboundParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
	})

	t.Run("successfully set EVM revert address V2", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "test")
		cctx.OutboundParams = cctx.OutboundParams[:1]
		cctx.RevertOptions.RevertAddress = sample.EthAddress().Hex()

		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.Equal(t, cctx.GetCurrentOutboundParam().Receiver, cctx.RevertOptions.RevertAddress)
		require.Equal(t, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.InboundParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount, cctx.OutboundParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, uint64(100))
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
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetAbort("test", "test")
	require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
	require.Contains(t, cctx.CctxStatus.ErrorMessage, "test")
}

func TestCrossChainTx_SetPendingRevert(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetPendingRevert("test", "test")
	require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
	require.Contains(t, cctx.CctxStatus.ErrorMessage, "test")
}

func TestCrossChainTx_SetPendingOutbound(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
	cctx.SetPendingOutbound("test")
	require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
	require.NotContains(t, cctx.CctxStatus.ErrorMessage, "test")
}

func TestCrossChainTx_SetOutboundMined(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetOutboundMined("test")
	require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
	require.NotContains(t, cctx.CctxStatus.ErrorMessage, "test")
}

func TestCrossChainTx_SetReverted(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
	cctx.SetReverted("test", "test")
	require.Equal(t, types.CctxStatus_Reverted, cctx.CctxStatus.Status)
	require.Contains(t, cctx.CctxStatus.StatusMessage, "test")
	require.Contains(t, cctx.CctxStatus.ErrorMessage, "test")
}
