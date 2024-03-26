package types_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
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
	inTxParams.SenderChainId = common.GoerliChain().ChainId
	inTxParams.Sender = "0x123"
	require.ErrorContains(t, inTxParams.Validate(), "invalid address 0x123")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.SenderChainId = common.GoerliChain().ChainId
	inTxParams.TxOrigin = "0x123"
	require.ErrorContains(t, inTxParams.Validate(), "invalid address 0x123")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.Amount = sdkmath.Uint{}
	require.ErrorContains(t, inTxParams.Validate(), "amount cannot be nil")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = "12"
	require.ErrorContains(t, inTxParams.Validate(), "hash must be a valid ethereum hash 12")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = sample.Hash().String()
	inTxParams.InboundTxBallotIndex = "12"
	require.ErrorContains(t, inTxParams.Validate(), "invalid index length 2")
	inTxParams = sample.InboundTxParamsValidChainID(r)
	inTxParams.InboundTxObservedHash = sample.Hash().String()
	inTxParams.InboundTxBallotIndex = sample.ZetaIndex(t)
	require.NoError(t, inTxParams.Validate())
}
