package sample

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func OutTxTracker(t *testing.T, index string) types.OutTxTracker {
	r := newRandFromStringSeed(t, index)

	return types.OutTxTracker{
		Index:   index,
		ChainId: r.Int63(),
		Nonce:   r.Uint64(),
	}
}

func Tss() *types.TSS {
	return &types.TSS{
		TssPubkey:           ed25519.GenPrivKey().PubKey().String(),
		FinalizedZetaHeight: 1000,
		KeyGenZetaHeight:    1000,
	}
}

func GasPrice(t *testing.T, index string) *types.GasPrice {
	r := newRandFromStringSeed(t, index)

	return &types.GasPrice{
		Creator:     AccAddress(),
		Index:       index,
		ChainId:     r.Int63(),
		Signers:     []string{AccAddress(), AccAddress()},
		BlockNums:   []uint64{r.Uint64(), r.Uint64()},
		Prices:      []uint64{r.Uint64(), r.Uint64()},
		MedianIndex: 0,
	}
}

func ChainNonces(t *testing.T, index string) *types.ChainNonces {
	r := newRandFromStringSeed(t, index)

	return &types.ChainNonces{
		Creator:         AccAddress(),
		Index:           index,
		ChainId:         r.Int63(),
		Nonce:           r.Uint64(),
		Signers:         []string{AccAddress(), AccAddress()},
		FinalizedHeight: r.Uint64(),
	}
}

func InboundTxParams(r *rand.Rand) *types.InboundTxParams {
	return &types.InboundTxParams{
		Sender:                          EthAddress().String(),
		SenderChainId:                   r.Int63(),
		TxOrigin:                        EthAddress().String(),
		CoinType:                        common.CoinType(r.Intn(100)),
		Asset:                           StringRandom(r, 32),
		Amount:                          math.NewUint(uint64(r.Int63())),
		InboundTxObservedHash:           StringRandom(r, 32),
		InboundTxObservedExternalHeight: r.Uint64(),
		InboundTxBallotIndex:            StringRandom(r, 32),
		InboundTxFinalizedZetaHeight:    r.Uint64(),
	}
}

func OutboundTxParams(r *rand.Rand) *types.OutboundTxParams {
	return &types.OutboundTxParams{
		Receiver:                         EthAddress().String(),
		ReceiverChainId:                  r.Int63(),
		CoinType:                         common.CoinType(r.Intn(100)),
		Amount:                           math.NewUint(uint64(r.Int63())),
		OutboundTxTssNonce:               r.Uint64(),
		OutboundTxGasLimit:               r.Uint64(),
		OutboundTxGasPrice:               String(),
		OutboundTxHash:                   StringRandom(r, 32),
		OutboundTxBallotIndex:            StringRandom(r, 32),
		OutboundTxObservedExternalHeight: r.Uint64(),
	}
}

func Status(t *testing.T, index string) *types.Status {
	r := newRandFromStringSeed(t, index)

	return &types.Status{
		Status:              types.CctxStatus(r.Intn(100)),
		StatusMessage:       String(),
		LastUpdateTimestamp: r.Int63(),
	}
}

func CrossChainTx(t *testing.T, index string) *types.CrossChainTx {
	r := newRandFromStringSeed(t, index)

	return &types.CrossChainTx{
		Creator:          AccAddress(),
		Index:            index,
		ZetaFees:         math.NewUint(uint64(r.Int63())),
		RelayedMessage:   StringRandom(r, 32),
		CctxStatus:       Status(t, index),
		InboundTxParams:  InboundTxParams(r),
		OutboundTxParams: []*types.OutboundTxParams{OutboundTxParams(r), OutboundTxParams(r)},
	}
}

func LastBlockHeight(t *testing.T, index string) *types.LastBlockHeight {
	r := newRandFromStringSeed(t, index)

	return &types.LastBlockHeight{
		Creator:           AccAddress(),
		Index:             index,
		Chain:             StringRandom(r, 32),
		LastSendHeight:    r.Uint64(),
		LastReceiveHeight: r.Uint64(),
	}
}

func InTxHashToCctx(t *testing.T, inTxHash string) types.InTxHashToCctx {
	r := newRandFromStringSeed(t, inTxHash)

	return types.InTxHashToCctx{
		InTxHash:  inTxHash,
		CctxIndex: []string{StringRandom(r, 32), StringRandom(r, 32)},
	}
}
