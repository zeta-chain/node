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

func TestCrossChainTx_GetConnectedChainID(t *testing.T) {
	t.Run("no inbound params", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams = nil
		_, _, err := cctx.GetConnectedChainID()
		require.Error(t, err)
		require.ErrorContains(t, err, "inbound params cannot be nil")
	})

	t.Run("no outbound params", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams = &types.InboundParams{
			SenderChainId: chains.ZetaChainPrivnet.GetChainId(),
		}
		cctx.OutboundParams = []*types.OutboundParams{}
		_, _, err := cctx.GetConnectedChainID()
		require.Error(t, err)
		require.ErrorContains(t, err, "outbound params cannot be nil")
	})

	t.Run("no outbound params with nil", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams = &types.InboundParams{
			SenderChainId: chains.ZetaChainPrivnet.GetChainId(),
		}
		cctx.OutboundParams = []*types.OutboundParams{nil}
		_, _, err := cctx.GetConnectedChainID()
		require.Error(t, err)
		require.ErrorContains(t, err, "outbound params cannot be nil")
	})

	t.Run("outgoing cctx", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams = &types.InboundParams{
			SenderChainId: chains.ZetaChainPrivnet.GetChainId(),
		}
		cctx.OutboundParams = []*types.OutboundParams{
			{
				ReceiverChainId: chains.BitcoinTestnet.GetChainId(),
			},
		}
		chainID, outgoing, err := cctx.GetConnectedChainID()
		require.NoError(t, err)
		require.True(t, outgoing)
		require.EqualValues(t, chains.BitcoinTestnet.GetChainId(), chainID)
	})

	t.Run("incoming cctx", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams = &types.InboundParams{
			SenderChainId: chains.BitcoinTestnet.GetChainId(),
		}
		cctx.OutboundParams = []*types.OutboundParams{
			{
				ReceiverChainId: chains.ZetaChainPrivnet.GetChainId(),
			},
		}
		chainID, outgoing, err := cctx.GetConnectedChainID()
		require.NoError(t, err)
		require.False(t, outgoing)
		require.EqualValues(t, chains.BitcoinTestnet.GetChainId(), chainID)
	})
}

func TestCrossChainTx_IsWithdrawTx(t *testing.T) {
	t.Run("withdraw tx", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.GetChainId()
		isZeta, err := cctx.IsWithdrawTx()
		require.NoError(t, err)
		require.True(t, isZeta)
	})

	t.Run("not withdraw tx", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams.SenderChainId = chains.BitcoinMainnet.GetChainId()
		isZeta, err := cctx.IsWithdrawTx()
		require.NoError(t, err)
		require.False(t, isZeta)
	})

	t.Run("no inbound params", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "sample")
		cctx.InboundParams = nil
		_, err := cctx.IsWithdrawTx()
		require.Error(t, err)
		require.ErrorContains(t, err, "inbound params cannot be nil")
	})
}

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
		require.Equal(t, cctx.GetCurrentOutboundParam().CoinType, cctx.InboundParams.CoinType)
		require.Equal(t, cctx.GetCurrentOutboundParam().ConfirmationMode, cctx.InboundParams.ConfirmationMode)
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
		require.Equal(t, cctx.GetCurrentOutboundParam().CoinType, cctx.InboundParams.CoinType)
	})

	t.Run("successfully set BTC revert address V2", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "test")
		cctx.OutboundParams = cctx.OutboundParams[:1]
		r := sample.Rand()
		cctx.InboundParams.SenderChainId = chains.BitcoinMainnet.ChainId
		cctx.RevertOptions.RevertAddress = sample.BTCAddressP2WPKH(t, r, &chaincfg.MainNetParams).String()

		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.Equal(t, cctx.GetCurrentOutboundParam().Receiver, cctx.RevertOptions.RevertAddress)
		require.Equal(t, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.InboundParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount, cctx.OutboundParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().TssPubkey, cctx.OutboundParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
		require.Equal(t, cctx.GetCurrentOutboundParam().CoinType, cctx.InboundParams.CoinType)
	})

	t.Run("successfully set SOL revert address V2", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "test")
		cctx.OutboundParams = cctx.OutboundParams[:1]
		cctx.InboundParams.SenderChainId = chains.SolanaDevnet.ChainId
		cctx.RevertOptions.RevertAddress = sample.SolanaAddress(t)

		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.Equal(t, cctx.GetCurrentOutboundParam().Receiver, cctx.RevertOptions.RevertAddress)
		require.Equal(t, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.InboundParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount, cctx.OutboundParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().TssPubkey, cctx.OutboundParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
		require.Equal(t, cctx.GetCurrentOutboundParam().CoinType, cctx.InboundParams.CoinType)
	})

	t.Run("successfully set SOL revert address V2 to inbound sender if revert address is invalid", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "test")
		cctx.OutboundParams = cctx.OutboundParams[:1]
		cctx.InboundParams.SenderChainId = chains.SolanaDevnet.ChainId
		cctx.RevertOptions.RevertAddress = sample.EthAddress().Hex()

		err := cctx.AddRevertOutbound(100)
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.Equal(t, cctx.GetCurrentOutboundParam().Receiver, cctx.InboundParams.Sender)
		require.Equal(t, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.InboundParams.SenderChainId)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount, cctx.OutboundParams[0].Amount)
		require.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().TssPubkey, cctx.OutboundParams[0].TssPubkey)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
		require.Equal(t, cctx.GetCurrentOutboundParam().CoinType, cctx.InboundParams.CoinType)
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
	t.Run("set abort from pending revert", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		cctx.SetAbort(types.StatusMessages{
			StatusMessage:      "status message",
			ErrorMessageRevert: "error revert",
		})
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, cctx.CctxStatus.StatusMessage, "status message")
		require.Equal(t, cctx.CctxStatus.ErrorMessageRevert, "error revert")
	})

	t.Run("set abort from pending outbound", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.SetAbort(types.StatusMessages{
			StatusMessage:        "status message",
			ErrorMessageOutbound: "error outbound",
		})
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, cctx.CctxStatus.StatusMessage, "status message")
		require.Equal(t, cctx.CctxStatus.ErrorMessage, "error outbound")
	})
}

func TestCrossChainTx_SetPendingRevert(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetPendingRevert(types.StatusMessages{
		StatusMessage:        "status message",
		ErrorMessageOutbound: "error outbound",
	})
	require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
	require.Equal(t, cctx.CctxStatus.StatusMessage, "status message")
	require.Equal(t, cctx.CctxStatus.ErrorMessage, "error outbound")
}

func TestCrossChainTx_SetPendingOutbound(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
	cctx.SetPendingOutbound(types.StatusMessages{
		StatusMessage: "status message",
	})
	require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
	require.Equal(t, cctx.CctxStatus.StatusMessage, "status message")
}

func TestCrossChainTx_SetOutboundMined(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	cctx.SetOutboundMined()
	require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(t, cctx.CctxStatus.StatusMessage, "")
}

func TestCrossChainTx_SetReverted(t *testing.T) {
	cctx := sample.CrossChainTx(t, "test")
	cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
	cctx.SetReverted()
	require.Equal(t, types.CctxStatus_Reverted, cctx.CctxStatus.Status)
	require.Equal(t, cctx.CctxStatus.StatusMessage, "")
	require.Equal(t, cctx.CctxStatus.ErrorMessageRevert, "")
}

func TestCrossChainTx_IsWithdrawAndCall(t *testing.T) {
	t.Run("withdraw and call", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.IsCrossChainCall = true
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		require.True(t, cctx.IsWithdrawAndCall())
	})

	t.Run("not withdraw and call", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.IsCrossChainCall = false
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		require.False(t, cctx.IsWithdrawAndCall())
	})

	t.Run("not pending outbound status", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.IsCrossChainCall = true
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		require.False(t, cctx.IsWithdrawAndCall())
	})

	t.Run("nil inbound", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams = nil
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		require.False(t, cctx.IsWithdrawAndCall())
	})

	t.Run("nil status", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.IsCrossChainCall = true
		cctx.CctxStatus = nil
		require.False(t, cctx.IsWithdrawAndCall())
	})
}
