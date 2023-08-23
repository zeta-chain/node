package sample

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/rand"
	"testing"
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
		Asset:                           String(),
		Amount:                          math.NewUint(uint64(r.Int63())),
		InboundTxObservedHash:           String(),
		InboundTxObservedExternalHeight: r.Uint64(),
		InboundTxBallotIndex:            String(),
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
		OutboundTxHash:                   String(),
		OutboundTxBallotIndex:            String(),
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
		RelayedMessage:   String(),
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
		Chain:             String(),
		LastSendHeight:    r.Uint64(),
		LastReceiveHeight: r.Uint64(),
	}
}

func InTxHashToCctx(inTxHash string) types.InTxHashToCctx {
	return types.InTxHashToCctx{
		InTxHash:  inTxHash,
		CctxIndex: []string{String(), String()},
	}
}
