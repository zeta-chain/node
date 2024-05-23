package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestInboundTxParams_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	inTxParams := sample.InboundTxParamsValidChainID(r)
	inTxParams.Sender = ""
	require.ErrorContains(t, inTxParams.Validate(), "sender cannot be empty")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.SenderChainId = 1000
	require.ErrorContains(t, inTxParams.Validate(), "invalid sender chain id 1000")

	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, inTxParams.Validate(), "amount cannot be nil")

	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = sample.Hash().String()
	inTxParams.InboundTxBallotIndex = sample.ZetaIndex(t)
	require.NoError(t, inTxParams.Validate())
}
