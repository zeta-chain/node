package sample

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func RateLimiterFlags() types.RateLimiterFlags {
	r := Rand()

	return types.RateLimiterFlags{
		Enabled: true,
		Window:  r.Int63(),
		Rate:    sdk.NewUint(r.Uint64()),
		Conversions: []types.Conversion{
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdk.NewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdk.NewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdk.NewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdk.NewDec(r.Int63()),
			},
			{
				Zrc20: EthAddress().Hex(),
				Rate:  sdk.NewDec(r.Int63()),
			},
		},
	}
}

// CustomRateLimiterFlags creates a custom rate limiter flags with the given parameters
func CustomRateLimiterFlags(enabled bool, window int64, rate math.Uint, conversions []types.Conversion) types.RateLimiterFlags {
	return types.RateLimiterFlags{
		Enabled:     enabled,
		Window:      window,
		Rate:        rate,
		Conversions: conversions,
	}
}

func AssetRate() types.AssetRate {
	r := Rand()

	return types.AssetRate{
		ChainId:  r.Int63(),
		Asset:    EthAddress().Hex(),
		Decimals: uint32(r.Uint64()),
		CoinType: coin.CoinType_ERC20,
		Rate:     sdk.NewDec(r.Int63()),
	}
}

// CustomAssetRate creates a custom asset rate with the given parameters
func CustomAssetRate(chainID int64, asset string, decimals uint32, coinType coin.CoinType, rate sdk.Dec) types.AssetRate {
	return types.AssetRate{
		ChainId:  chainID,
		Asset:    strings.ToLower(asset),
		Decimals: decimals,
		CoinType: coinType,
		Rate:     rate,
	}
}

func OutboundTracker(t *testing.T, index string) types.OutboundTracker {
	r := newRandFromStringSeed(t, index)

	return types.OutboundTracker{
		Index:   index,
		ChainId: r.Int63(),
		Nonce:   r.Uint64(),
	}
}

func InboundTracker(t *testing.T, index string) types.InboundTracker {
	r := newRandFromStringSeed(t, index)

	return types.InboundTracker{
		ChainId:  r.Int63(),
		CoinType: coin.CoinType_Zeta,
		TxHash:   Hash().Hex(),
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

func InboundParams(r *rand.Rand) *types.InboundParams {
	return &types.InboundParams{
		Sender:                 EthAddress().String(),
		SenderChainId:          r.Int63(),
		TxOrigin:               EthAddress().String(),
		CoinType:               coin.CoinType(r.Intn(100)),
		Asset:                  StringRandom(r, 32),
		Amount:                 math.NewUint(uint64(r.Int63())),
		ObservedHash:           StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		BallotIndex:            StringRandom(r, 32),
		FinalizedZetaHeight:    r.Uint64(),
	}
}

func InboundParamsValidChainID(r *rand.Rand) *types.InboundParams {
	return &types.InboundParams{
		Sender:                 EthAddress().String(),
		SenderChainId:          chains.GoerliChain.ChainId,
		TxOrigin:               EthAddress().String(),
		Asset:                  StringRandom(r, 32),
		Amount:                 math.NewUint(uint64(r.Int63())),
		ObservedHash:           StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		BallotIndex:            StringRandom(r, 32),
		FinalizedZetaHeight:    r.Uint64(),
	}
}

func OutboundParams(r *rand.Rand) *types.OutboundParams {
	return &types.OutboundParams{
		Receiver:               EthAddress().String(),
		ReceiverChainId:        r.Int63(),
		CoinType:               coin.CoinType(r.Intn(100)),
		Amount:                 math.NewUint(uint64(r.Int63())),
		TssNonce:               r.Uint64(),
		GasLimit:               r.Uint64(),
		GasPrice:               math.NewUint(uint64(r.Int63())).String(),
		Hash:                   StringRandom(r, 32),
		BallotIndex:            StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		GasUsed:                r.Uint64(),
		EffectiveGasPrice:      math.NewInt(r.Int63()),
	}
}

func OutboundParamsValidChainID(r *rand.Rand) *types.OutboundParams {
	return &types.OutboundParams{
		Receiver:               EthAddress().String(),
		ReceiverChainId:        chains.GoerliChain.ChainId,
		Amount:                 math.NewUint(uint64(r.Int63())),
		TssNonce:               r.Uint64(),
		GasLimit:               r.Uint64(),
		GasPrice:               math.NewUint(uint64(r.Int63())).String(),
		Hash:                   StringRandom(r, 32),
		BallotIndex:            StringRandom(r, 32),
		ObservedExternalHeight: r.Uint64(),
		GasUsed:                r.Uint64(),
		EffectiveGasPrice:      math.NewInt(r.Int63()),
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

func GetCctxIndexFromString(index string) string {
	return crypto.Keccak256Hash([]byte(index)).String()
}

func CrossChainTx(t *testing.T, index string) *types.CrossChainTx {
	r := newRandFromStringSeed(t, index)

	return &types.CrossChainTx{
		Creator:        AccAddress(),
		Index:          GetCctxIndexFromString(index),
		ZetaFees:       math.NewUint(uint64(r.Int63())),
		RelayedMessage: StringRandom(r, 32),
		CctxStatus:     Status(t, index),
		InboundParams:  InboundParams(r),
		OutboundParams: []*types.OutboundParams{OutboundParams(r), OutboundParams(r)},
	}
}

// CustomCctxsInBlockRange create 1 cctx per block in block range [lowBlock, highBlock] (inclusive)
func CustomCctxsInBlockRange(
	t *testing.T,
	lowBlock uint64,
	highBlock uint64,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	amount uint64,
	status types.CctxStatus,
) (cctxs []*types.CrossChainTx) {
	// create 1 cctx per block
	for i := lowBlock; i <= highBlock; i++ {
		nonce := i - 1
		cctx := CrossChainTx(t, fmt.Sprintf("%d-%d", chainID, nonce))
		cctx.CctxStatus.Status = status
		cctx.InboundParams.CoinType = coinType
		cctx.InboundParams.Asset = asset
		cctx.InboundParams.ObservedExternalHeight = i
		cctx.GetCurrentOutboundParam().ReceiverChainId = chainID
		cctx.GetCurrentOutboundParam().Amount = sdk.NewUint(amount)
		cctx.GetCurrentOutboundParam().TssNonce = nonce
		cctxs = append(cctxs, cctx)
	}
	return cctxs
}

func LastBlockHeight(t *testing.T, index string) *types.LastBlockHeight {
	r := newRandFromStringSeed(t, index)

	return &types.LastBlockHeight{
		Creator:            AccAddress(),
		Index:              index,
		Chain:              StringRandom(r, 32),
		LastInboundHeight:  r.Uint64(),
		LastOutboundHeight: r.Uint64(),
	}
}

func InboundHashToCctx(t *testing.T, inboundHash string) types.InboundHashToCctx {
	r := newRandFromStringSeed(t, inboundHash)

	return types.InboundHashToCctx{
		InboundHash: inboundHash,
		CctxIndex:   []string{StringRandom(r, 32), StringRandom(r, 32)},
	}
}

func ZetaAccounting(t *testing.T, index string) types.ZetaAccounting {
	r := newRandFromStringSeed(t, index)
	return types.ZetaAccounting{
		AbortedZetaAmount: math.NewUint(uint64(r.Int63())),
	}
}

func InboundVote(coinType coin.CoinType, from, to int64) types.MsgVoteInbound {
	return types.MsgVoteInbound{
		Creator:            "",
		Sender:             EthAddress().String(),
		SenderChainId:      Chain(from).GetChainId(),
		Receiver:           EthAddress().String(),
		ReceiverChain:      Chain(to).GetChainId(),
		Amount:             UintInRange(10000000, 1000000000),
		Message:            base64.StdEncoding.EncodeToString(Bytes()),
		InboundBlockHeight: Uint64InRange(1, 10000),
		GasLimit:           1000000000,
		InboundHash:        Hash().String(),
		CoinType:           coinType,
		TxOrigin:           EthAddress().String(),
		Asset:              "",
		EventIndex:         EventIndex(),
	}
}

// receiver is 1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3
func GetInvalidZRC20WithdrawToExternal(t *testing.T) (receipt ethtypes.Receipt) {
	block := "{\n  \"type\": \"0x2\",\n  \"root\": \"0x\",\n  \"status\": \"0x1\",\n  \"cumulativeGasUsed\": \"0x4e7a38\",\n  \"logsBloom\": \"0x00000000000000000000010000020000000000000000000000000000000000020000000100000000000000000000000080000000000000000000000400200000200000000002000000000008000000000000000000000000000000000000000000000000020000000000000000800800000040000000000000000010000000000000000000000000000000000000000000000000000004000000000000000000020000000000000000000000000000000000000000000000000000000000010000000002000000000000000000000000000000000000000000000000000020000010000000000000000001000000000000000000040200000000000000000000\",\n  \"logs\": [\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"\n      ],\n      \"data\": \"0x0000000000000000000000000000000000000000000000000000000000003790\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x46\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x00000000000000000000000013a0c5930c028511dc02665e7285134b6d11a5f4\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000000000006a1217\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x47\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x0000000000000000000000000000000000000000000000000000000000000000\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000000000006a0c70\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x48\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\"\n      ],\n      \"data\": \"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000006a0c700000000000000000000000000000000000000000000000000000000000003790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000223145595676584c7573437874567545776f59765752794e35455a5458775056766f33000000000000000000000000000000000000000000000000000000000000\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x49\",\n      \"removed\": false\n    }\n  ],\n  \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n  \"contractAddress\": \"0x0000000000000000000000000000000000000000\",\n  \"gasUsed\": \"0x12521\",\n  \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n  \"blockNumber\": \"0x1a2ad3\",\n  \"transactionIndex\": \"0x10\"\n}\n"
	err := json.Unmarshal([]byte(block), &receipt)
	require.NoError(t, err)
	return
}

func GetValidZrc20WithdrawToETH(t *testing.T) (receipt ethtypes.Receipt) {
	block := "{\n  \"type\": \"0x2\",\n  \"root\": \"0x\",\n  \"status\": \"0x1\",\n  \"cumulativeGasUsed\": \"0xdbedca\",\n  \"logsBloom\": \"0x00200000001000000000000088020001000001000000000000000000000000000000020100000000000000000000000080000000000000000000000400640000000000000000000008000008020000200000000000000002000000008000000000000000020000000200000000800801000000080000000000000010000000000000000000000000000000000000001000000001000004080001404000000000028002000000000000000040000000000000000000000000000000000000000000000002000000000000008000000000000000800800001000000002000021000010000100000000000010800400000000020000000100400880000000004000\",\n  \"logs\": [\n    {\n      \"address\": \"0x3f641963f3d9adf82d890fd8142313dcec807ba5\",\n      \"topics\": [\n        \"0x3d0ce9bfc3ed7d6862dbb28b2dea94561fe714a1b4d019aa8af39730d1ad7c3d\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000045400a8fd5330000\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x57\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\n      \"topics\": [\n        \"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000001ac7c4159f72b90000\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x58\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\n      \"topics\": [\n        \"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n        \"0x0000000000000000000000002ca7d64a7efe2d62a725e2b35cf7230d6677ffee\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000001ac7c4159f72b90000\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x59\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n        \"0x00000000000000000000000016ef1b018026e389fda93c1e993e987cf6e852e7\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000001ac7c4159f72b90000\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x5a\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0xd97b1de3619ed2c6beb3860147e30ca8a7dc9891\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x00000000000000000000000016ef1b018026e389fda93c1e993e987cf6e852e7\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000002e640d76638740f\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x5b\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x16ef1b018026e389fda93c1e993e987cf6e852e7\",\n      \"topics\": [\n        \"0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1\"\n      ],\n      \"data\": \"0x000000000000000000000000000000000000000000000b3f1da425061770a11600000000000000000000000000000000000000000000000135be3952e251aa40\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x5c\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x16ef1b018026e389fda93c1e993e987cf6e852e7\",\n      \"topics\": [\n        \"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822\",\n        \"0x0000000000000000000000002ca7d64a7efe2d62a725e2b35cf7230d6677ffee\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000001ac7c4159f72b900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002e640d76638740f\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x5d\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0xd97b1de3619ed2c6beb3860147e30ca8a7dc9891\",\n      \"topics\": [\n        \"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n        \"0x000000000000000000000000d97b1de3619ed2c6beb3860147e30ca8a7dc9891\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000015059f36c8ec0\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x5e\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0xd97b1de3619ed2c6beb3860147e30ca8a7dc9891\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n        \"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000015059f36c8ec0\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x5f\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0xd97b1de3619ed2c6beb3860147e30ca8a7dc9891\",\n      \"topics\": [\n        \"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n        \"0x000000000000000000000000d97b1de3619ed2c6beb3860147e30ca8a7dc9891\"\n      ],\n      \"data\": \"0x0000000000000000000000000000000000000000000000000000000000000000\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x60\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0xd97b1de3619ed2c6beb3860147e30ca8a7dc9891\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n        \"0x0000000000000000000000000000000000000000000000000000000000000000\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000002e4f07d72cbe54f\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x61\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0xd97b1de3619ed2c6beb3860147e30ca8a7dc9891\",\n      \"topics\": [\n        \"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\n        \"0x0000000000000000000000008e0f8e7e9e121403e72151d00f4937eacb2d9ef3\"\n      ],\n      \"data\": \"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000002e4f07d72cbe54f00000000000000000000000000000000000000000000000000015059f36c8ec0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000005dabfdd153aaab4a970fd953dcfeee8bf6bb946e\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x62\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x8e0f8e7e9e121403e72151d00f4937eacb2d9ef3\",\n      \"topics\": [\n        \"0x97eb75cc53ffa3f4560fc62e4912dda10ac56c3d12dbc48dc8c27d5ab225cf66\"\n      ],\n      \"data\": \"0x0000000000000000000000005f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf000000000000000000000000d97b1de3619ed2c6beb3860147e30ca8a7dc989100000000000000000000000000000000000000000000001b0d04202f47ec000000000000000000000000000000000000000000000000001ac7c4159f72b900000000000000000000000000005dabfdd153aaab4a970fd953dcfeee8bf6bb946e00000000000000000000000000000000000000000000000045400a8fd5330000\",\n      \"blockNumber\": \"0x17ef22\",\n      \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n      \"transactionIndex\": \"0x1f\",\n      \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n      \"logIndex\": \"0x63\",\n      \"removed\": false\n    }\n  ],\n  \"transactionHash\": \"0x87229bb05d67f42017a697b34ed52d95afc9f5e3285479e845fe088b4c77d8f0\",\n  \"contractAddress\": \"0x0000000000000000000000000000000000000000\",\n  \"gasUsed\": \"0x41c3c\",\n  \"blockHash\": \"0xf49e7039c7f1a81cd46de150980d92fa869cc0d2e2f1fe46aedc6400396137ff\",\n  \"blockNumber\": \"0x17ef22\",\n  \"transactionIndex\": \"0x1f\"\n}"
	err := json.Unmarshal([]byte(block), &receipt)
	require.NoError(t, err)
	return

}

// receiver is bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu
func GetValidZRC20WithdrawToBTC(t *testing.T) (receipt ethtypes.Receipt) {
	block := "{\"type\":\"0x2\",\"root\":\"0x\",\"status\":\"0x1\",\"cumulativeGasUsed\":\"0x1f25ed\",\"logsBloom\":\"0x00000000000000000000000000020000000000000000000000000000000000020000000100000000000000000040000080000000000000000000000400200000200000000002000000000008000000000000000000000000000000000000000000000000020000000000000000800800000000000000000000000010000000000000000000000000000000000000000000000000000004000000000000000000020000000001000000000000000000000000000000000000000000000000010000000002000000000000000010000000000000000000000000000000000020000010000000000000000000000000000000000000040200000000000000000000\",\"logs\":[{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000003d84\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x28\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x00000000000000000000000013a0c5930c028511dc02665e7285134b6d11a5f4\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000978c98\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x29\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x0000000000000000000000000000000000000000000000000000000000000000\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000003039\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x2a\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\"],\"data\":\"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000030390000000000000000000000000000000000000000000000000000000000003d840000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a626331717973643473703971386d793539756c3977736635727673397033383768663876667761747a7500000000000000000000000000000000000000000000\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x2b\",\"removed\":false}],\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"contractAddress\":\"0x0000000000000000000000000000000000000000\",\"gasUsed\":\"0x12575\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"blockNumber\":\"0x1a00f3\",\"transactionIndex\":\"0x8\"}\n"
	err := json.Unmarshal([]byte(block), &receipt)
	require.NoError(t, err)
	return
}

func GetValidZetaSentDestinationExternal(t *testing.T) (receipt ethtypes.Receipt) {
	block := "{\"root\":\"0x\",\"status\":\"0x1\",\"cumulativeGasUsed\":\"0xd75f4f\",\"logsBloom\":\"0x00000000000000000000000000000000800800000000000000000000100000000000002000000100000000000000000000000000000000000000000000240000000000000000000000000008000000000800000000440000000000008080000000000000000000000000000000000000000000000000040000000010000000000000000000000000000000000000000200000001000000000000000040000000020000000000000000000000008200000000000000000000000000000000000000000002000000000000008000000000000000000000000000080002000041000010000000000000000000000000000000000000000000400000000000000000\",\"logs\":[{\"address\":\"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\"topics\":[\"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c\",\"0x000000000000000000000000f0a3f93ed1b126142e61423f9546bf1323ff82df\"],\"data\":\"0x000000000000000000000000000000000000000000000003cb71f51fc5580000\",\"blockNumber\":\"0x1bedc8\",\"transactionHash\":\"0x19d8a67a05998f1cb19fe731b96d817d5b186b62c9430c51679664959c952ef0\",\"transactionIndex\":\"0x5f\",\"blockHash\":\"0x198fdd1f4bc6b910db978602cb15bdb2bcc6fd960e9324e9b9675dc062133794\",\"logIndex\":\"0x13b\",\"removed\":false},{\"address\":\"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x000000000000000000000000f0a3f93ed1b126142e61423f9546bf1323ff82df\",\"0x000000000000000000000000239e96c8f17c85c30100ac26f635ea15f23e9c67\"],\"data\":\"0x000000000000000000000000000000000000000000000003cb71f51fc5580000\",\"blockNumber\":\"0x1bedc8\",\"transactionHash\":\"0x19d8a67a05998f1cb19fe731b96d817d5b186b62c9430c51679664959c952ef0\",\"transactionIndex\":\"0x5f\",\"blockHash\":\"0x198fdd1f4bc6b910db978602cb15bdb2bcc6fd960e9324e9b9675dc062133794\",\"logIndex\":\"0x13c\",\"removed\":false},{\"address\":\"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x000000000000000000000000f0a3f93ed1b126142e61423f9546bf1323ff82df\",\"0x000000000000000000000000239e96c8f17c85c30100ac26f635ea15f23e9c67\"],\"data\":\"0x000000000000000000000000000000000000000000000003cb71f51fc5580000\",\"blockNumber\":\"0x1bedc8\",\"transactionHash\":\"0x19d8a67a05998f1cb19fe731b96d817d5b186b62c9430c51679664959c952ef0\",\"transactionIndex\":\"0x5f\",\"blockHash\":\"0x198fdd1f4bc6b910db978602cb15bdb2bcc6fd960e9324e9b9675dc062133794\",\"logIndex\":\"0x13d\",\"removed\":false},{\"address\":\"0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf\",\"topics\":[\"0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65\",\"0x000000000000000000000000239e96c8f17c85c30100ac26f635ea15f23e9c67\"],\"data\":\"0x000000000000000000000000000000000000000000000003cb71f51fc5580000\",\"blockNumber\":\"0x1bedc8\",\"transactionHash\":\"0x19d8a67a05998f1cb19fe731b96d817d5b186b62c9430c51679664959c952ef0\",\"transactionIndex\":\"0x5f\",\"blockHash\":\"0x198fdd1f4bc6b910db978602cb15bdb2bcc6fd960e9324e9b9675dc062133794\",\"logIndex\":\"0x13e\",\"removed\":false},{\"address\":\"0x239e96c8f17c85c30100ac26f635ea15f23e9c67\",\"topics\":[\"0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4\",\"0x000000000000000000000000f0a3f93ed1b126142e61423f9546bf1323ff82df\",\"0x0000000000000000000000000000000000000000000000000000000000000001\"],\"data\":\"0x00000000000000000000000060983881bdf302dcfa96603a58274d15d596620900000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000003cb71f51fc558000000000000000000000000000000000000000000000000000000000000000186a000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000001460983881bdf302dcfa96603a58274d15d59662090000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000\",\"blockNumber\":\"0x1bedc8\",\"transactionHash\":\"0x19d8a67a05998f1cb19fe731b96d817d5b186b62c9430c51679664959c952ef0\",\"transactionIndex\":\"0x5f\",\"blockHash\":\"0x198fdd1f4bc6b910db978602cb15bdb2bcc6fd960e9324e9b9675dc062133794\",\"logIndex\":\"0x13f\",\"removed\":false}],\"transactionHash\":\"0x19d8a67a05998f1cb19fe731b96d817d5b186b62c9430c51679664959c952ef0\",\"contractAddress\":\"0x0000000000000000000000000000000000000000\",\"gasUsed\":\"0x2406d\",\"blockHash\":\"0x198fdd1f4bc6b910db978602cb15bdb2bcc6fd960e9324e9b9675dc062133794\",\"blockNumber\":\"0x1bedc8\",\"transactionIndex\":\"0x5f\"}\n"
	err := json.Unmarshal([]byte(block), &receipt)
	require.NoError(t, err)
	return
}
