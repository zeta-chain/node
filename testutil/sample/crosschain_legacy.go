package sample

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CrossChainTxV14(t *testing.T, index string) *types.CrossChainTxV14 {
	r := newRandFromStringSeed(t, index)
	cointType := common.CoinType(r.Intn(100))
	return &types.CrossChainTxV14{
		Creator:          AccAddress(),
		Index:            GetCctxIndexFromString(index),
		ZetaFees:         math.NewUint(uint64(r.Int63())),
		RelayedMessage:   StringRandom(r, 32),
		CctxStatus:       Status(t, index),
		InboundTxParams:  InboundTxParamsV14(r, cointType),
		OutboundTxParams: []*types.OutboundTxParamsV14{OutboundTxParamsV14(r, cointType), OutboundTxParamsV14(r, cointType)},
	}
}

func InboundTxParamsV14(r *rand.Rand, coinType common.CoinType) *types.InboundTxParamsV14 {
	return &types.InboundTxParamsV14{
		Sender:                          EthAddress().String(),
		SenderChainId:                   r.Int63(),
		TxOrigin:                        EthAddress().String(),
		Asset:                           StringRandom(r, 32),
		Amount:                          math.NewUint(uint64(r.Int63())),
		InboundTxObservedHash:           StringRandom(r, 32),
		InboundTxObservedExternalHeight: r.Uint64(),
		InboundTxBallotIndex:            StringRandom(r, 32),
		InboundTxFinalizedZetaHeight:    r.Uint64(),
		CoinType:                        coinType,
	}
}

func OutboundTxParamsV14(r *rand.Rand, coinType common.CoinType) *types.OutboundTxParamsV14 {
	return &types.OutboundTxParamsV14{
		Receiver:                         EthAddress().String(),
		ReceiverChainId:                  r.Int63(),
		Amount:                           math.NewUint(uint64(r.Int63())),
		OutboundTxTssNonce:               r.Uint64(),
		OutboundTxGasLimit:               r.Uint64(),
		OutboundTxGasPrice:               math.NewUint(uint64(r.Int63())).String(),
		OutboundTxHash:                   StringRandom(r, 32),
		OutboundTxBallotIndex:            StringRandom(r, 32),
		OutboundTxObservedExternalHeight: r.Uint64(),
		OutboundTxGasUsed:                r.Uint64(),
		OutboundTxEffectiveGasPrice:      math.NewInt(r.Int63()),
		CoinType:                         coinType,
	}
}
