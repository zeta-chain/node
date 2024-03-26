package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestOutboundTxParams_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	outTxParams := sample.OutboundTxParamsValidChainID(r)
	outTxParams.Receiver = ""
	require.ErrorContains(t, outTxParams.Validate(), "receiver cannot be empty")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.ReceiverChainId = 1000
	require.ErrorContains(t, outTxParams.Validate(), "invalid receiver chain id 1000")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.Receiver = "0x123"
	require.ErrorContains(t, outTxParams.Validate(), "invalid address 0x123")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, outTxParams.Validate(), "amount cannot be nil")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.OutboundTxBallotIndex = "12"
	require.ErrorContains(t, outTxParams.Validate(), "invalid index length 2")
	outTxParams = sample.OutboundTxParamsValidChainID(r)
	outTxParams.OutboundTxBallotIndex = sample.ZetaIndex(t)
	outTxParams.OutboundTxHash = sample.Hash().String()
	require.NoError(t, outTxParams.Validate())
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
