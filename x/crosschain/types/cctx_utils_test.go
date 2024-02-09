package types_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

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

func TestOutboundTxParams_GetGasPrice(t *testing.T) {
	// #nosec G404 - random seed is not used for security purposes
	r := rand.New(rand.NewSource(42))
	outTxParams := sample.OutboundTxParams(r)

	outTxParams.OutboundTxGasPrice = "42"
	gasPrice, err := outTxParams.GetGasPrice()
	require.NoError(t, err)
	require.EqualValues(t, uint64(42), gasPrice)

	outTxParams.OutboundTxGasPrice = "invalid"
	_, err = outTxParams.GetGasPrice()
	require.Error(t, err)
}
