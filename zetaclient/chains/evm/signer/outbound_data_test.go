package signer

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestNewOutboundData(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))

	ctx := makeCtx(t)

	newOutbound := func(cctx *types.CrossChainTx) (*OutboundData, bool, error) {
		return NewOutboundData(ctx, cctx, 123, logger)
	}

	t.Run("success", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)

		// ACT
		out, skip, err := newOutbound(cctx)

		// ASSERT
		require.NoError(t, err)
		assert.False(t, skip)

		assert.NotEmpty(t, out)

		assert.NotEmpty(t, out.srcChainID)
		assert.NotEmpty(t, out.sender)

		assert.NotEmpty(t, out.toChainID)
		assert.NotEmpty(t, out.to)

		assert.Equal(t, ethcommon.HexToAddress(cctx.InboundParams.Asset), out.asset)
		assert.NotEmpty(t, out.amount)

		assert.NotEmpty(t, out.nonce)
		assert.NotEmpty(t, out.zetaHeight)
		assert.NotEmpty(t, out.gas)
		assert.True(t, out.gas.isLegacy())
		assert.Equal(t, cctx.GetCurrentOutboundParam().CallOptions.GasLimit, out.gas.Limit)

		assert.Empty(t, out.message)
		assert.NotEmpty(t, out.cctxIndex)
		assert.Equal(t, cctx.OutboundParams[0], out.outboundParams)
	})

	t.Run("pending revert", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert

		// ACT
		out, skip, err := newOutbound(cctx)

		// ASSERT
		require.NoError(t, err)
		assert.False(t, skip)
		assert.Equal(t, ethcommon.HexToAddress(cctx.InboundParams.Sender), out.to)
		assert.Equal(t, big.NewInt(cctx.InboundParams.SenderChainId), out.toChainID)
	})

	t.Run("pending outbound", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound

		// ACT
		out, skip, err := newOutbound(cctx)

		// ASSERT
		assert.NoError(t, err)
		assert.False(t, skip)
		assert.Equal(t, ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver), out.to)
		assert.Equal(t, big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId), out.toChainID)
	})

	t.Run("skip inbound", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)
		cctx.CctxStatus.Status = types.CctxStatus_PendingInbound

		// ACT
		_, skip, err := newOutbound(cctx)

		// ASSERT
		require.NoError(t, err)
		assert.True(t, skip)
	})

	t.Run("skip aborted", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)
		cctx.CctxStatus.Status = types.CctxStatus_Aborted

		// ACT
		_, skip, err := newOutbound(cctx)

		// ASSERT
		require.NoError(t, err)
		assert.True(t, skip)
	})

	t.Run("invalid gas price", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)
		cctx.GetCurrentOutboundParam().GasPrice = "invalidGasPrice"

		// ACT
		_, _, err := newOutbound(cctx)

		// ASSERT
		assert.ErrorContains(t, err, "unable to parse gasPrice")
	})

	t.Run("unknown chain", func(t *testing.T) {
		// ARRANGE
		cctx := getInvalidCCTX(t)

		// ACT
		_, _, err := newOutbound(cctx)

		// ASSERT
		assert.ErrorContains(t, err, "chain not found")
	})

	t.Run("no outbound params", func(t *testing.T) {
		// ARRANGE
		cctx := getCCTX(t)
		cctx.OutboundParams = nil

		// ACT
		_, _, err := newOutbound(cctx)

		// ASSERT
		assert.ErrorContains(t, err, "outboundParams is empty")
	})
}
