package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestInboundTxParams_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	inboundParams := sample.InboundParamsValidChainID(r)
	inboundParams.Sender = ""
	require.ErrorContains(t, inboundParams.Validate(), "sender cannot be empty")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.SenderChainId = 1000
	require.ErrorContains(t, inboundParams.Validate(), "invalid sender chain id 1000")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.SenderChainId = chains.GoerliChain.ChainId
	inboundParams.Sender = "0x123"
	require.ErrorContains(t, inboundParams.Validate(), "invalid address 0x123")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.SenderChainId = chains.GoerliChain.ChainId
	inboundParams.TxOrigin = "0x123"
	require.ErrorContains(t, inboundParams.Validate(), "invalid address 0x123")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, inboundParams.Validate(), "amount cannot be nil")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.ObservedHash = "12"
	require.ErrorContains(t, inboundParams.Validate(), "hash must be a valid ethereum hash 12")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.ObservedHash = sample.Hash().String()
	inboundParams.BallotIndex = "12"
	require.ErrorContains(t, inboundParams.Validate(), "invalid index length 2")
	inboundParams = sample.InboundParamsValidChainID(r)
	inboundParams.ObservedHash = sample.Hash().String()
	inboundParams.BallotIndex = sample.ZetaIndex(t)
	require.NoError(t, inboundParams.Validate())
}
