package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestOutboundParams_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	outTxParams := sample.OutboundParamsValidChainID(r)
	outTxParams.Receiver = ""
	require.ErrorContains(t, outTxParams.Validate(), "receiver cannot be empty")
	outTxParams = sample.OutboundParamsValidChainID(r)
	outTxParams.ReceiverChainId = 1000
	require.ErrorContains(t, outTxParams.Validate(), "invalid receiver chain id 1000")
	outTxParams = sample.OutboundParamsValidChainID(r)
	outTxParams.Receiver = "0x123"
	require.ErrorContains(t, outTxParams.Validate(), "invalid address 0x123")
	outTxParams = sample.OutboundParamsValidChainID(r)
	outTxParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, outTxParams.Validate(), "amount cannot be nil")
	outTxParams = sample.OutboundParamsValidChainID(r)
	outTxParams.BallotIndex = "12"
	require.ErrorContains(t, outTxParams.Validate(), "invalid index length 2")
	outTxParams = sample.OutboundParamsValidChainID(r)
	outTxParams.BallotIndex = sample.ZetaIndex(t)
	outTxParams.Hash = sample.Hash().String()
	require.NoError(t, outTxParams.Validate())
}

func TestOutboundTxParams_GetGasPrice(t *testing.T) {
	// #nosec G404 - random seed is not used for security purposes
	r := rand.New(rand.NewSource(42))
	outTxParams := sample.OutboundParams(r)

	outTxParams.GasPrice = "42"
	gasPrice, err := outTxParams.GetGasPriceUInt64()
	require.NoError(t, err)
	require.EqualValues(t, uint64(42), gasPrice)

	outTxParams.GasPrice = "invalid"
	_, err = outTxParams.GetGasPriceUInt64()
	require.Error(t, err)
}
