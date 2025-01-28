package signer

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/config"
)

func Test_NewOutboundData(t *testing.T) {
	// sample address
	chain := chains.BitcoinMainnet
	receiver, err := chains.DecodeBtcAddress("bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y", chain.ChainId)
	require.NoError(t, err)

	// setup compliance config
	cfg := config.Config{
		ComplianceConfig: sample.ComplianceConfig(),
	}
	config.LoadComplianceConfig(cfg)

	// test cases
	tests := []struct {
		name         string
		cctx         *crosschaintypes.CrossChainTx
		cctxModifier func(cctx *crosschaintypes.CrossChainTx)
		chainID      int64
		height       uint64
		minRelayFee  float64
		expected     *OutboundData
		errMsg       string
	}{
		{
			name: "create new outbound data successfully",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().Receiver = receiver.String()
				cctx.GetCurrentOutboundParam().ReceiverChainId = chain.ChainId
				cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(1e7) // 0.1 BTC
				cctx.GetCurrentOutboundParam().CallOptions.GasLimit = 254    // 254 bytes
				cctx.GetCurrentOutboundParam().GasPrice = "10"               // 10 sats/vByte
				cctx.GetCurrentOutboundParam().TssNonce = 1
			},
			chainID:     chain.ChainId,
			height:      101,
			minRelayFee: 0.00001, // 1000 sat/KB
			expected: &OutboundData{
				chainID:    chain.ChainId,
				to:         receiver,
				amount:     0.1,
				amountSats: 10000000,
				feeRate:    8, // Round(7.5)
				txSize:     254,
				nonce:      1,
				height:     101,
				cancelTx:   false,
			},
			errMsg: "",
		},
		{
			name: "create new outbound data using current gas rate instead of old rate",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().Receiver = receiver.String()
				cctx.GetCurrentOutboundParam().ReceiverChainId = chain.ChainId
				cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(1e7) // 0.1 BTC
				cctx.GetCurrentOutboundParam().CallOptions.GasLimit = 254    // 254 bytes
				cctx.GetCurrentOutboundParam().GasPrice = "10"               // 10 sats/vByte
				cctx.GetCurrentOutboundParam().GasPriorityFee = "15"         // 15 sats/vByte
				cctx.GetCurrentOutboundParam().TssNonce = 1
			},
			chainID:     chain.ChainId,
			height:      101,
			minRelayFee: 0.00001, // 1000 sat/KB
			expected: &OutboundData{
				chainID:    chain.ChainId,
				to:         receiver,
				amount:     0.1,
				amountSats: 10000000,
				feeRate:    11, // Round(11.25)
				txSize:     254,
				nonce:      1,
				height:     101,
				cancelTx:   false,
			},
			errMsg: "",
		},
		{
			name:         "cctx is nil",
			cctx:         nil,
			cctxModifier: nil,
			expected:     nil,
			errMsg:       "cctx is nil",
		},
		{
			name: "coin type is not gas",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_ERC20
			},
			expected: nil,
			errMsg:   "can only send gas token to a Bitcoin network",
		},
		{
			name: "invalid gas price",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().GasPrice = "invalid"
			},
			expected: nil,
			errMsg:   "invalid fee rate",
		},
		{
			name: "zero fee rate",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().GasPrice = "0"
			},
			expected: nil,
			errMsg:   "invalid fee rate",
		},
		{
			name: "invalid receiver address",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().Receiver = "invalid"
			},
			expected: nil,
			errMsg:   "cannot decode receiver address",
		},
		{
			name: "unsupported receiver address",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().Receiver = "035e4ae279bd416b5da724972c9061ec6298dac020d1e3ca3f06eae715135cdbec"
				cctx.GetCurrentOutboundParam().ReceiverChainId = chain.ChainId
			},
			expected: nil,
			errMsg:   "unsupported receiver address",
		},
		{
			name: "invalid gas limit",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().Receiver = receiver.String()
				cctx.GetCurrentOutboundParam().ReceiverChainId = chain.ChainId
				cctx.GetCurrentOutboundParam().CallOptions.GasLimit = math.MaxInt64 + 1
			},
			expected: nil,
			errMsg:   "invalid gas limit",
		},
		{
			name: "should cancel restricted CCTX",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.InboundParams.Sender = sample.RestrictedEVMAddressTest
				cctx.GetCurrentOutboundParam().Receiver = receiver.String()
				cctx.GetCurrentOutboundParam().ReceiverChainId = chain.ChainId
				cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(1e7) // 0.1 BTC
				cctx.GetCurrentOutboundParam().CallOptions.GasLimit = 254    // 254 bytes
				cctx.GetCurrentOutboundParam().GasPrice = "10"               // 10 sats/vByte
				cctx.GetCurrentOutboundParam().TssNonce = 1
			},
			chainID:     chain.ChainId,
			height:      101,
			minRelayFee: 0.00001, // 1000 sat/KB
			expected: &OutboundData{
				chainID:    chain.ChainId,
				to:         receiver,
				amount:     0, // should cancel the tx
				amountSats: 0,
				feeRate:    8, // Round(7.5)
				txSize:     254,
				nonce:      1,
				height:     101,
				cancelTx:   true,
			},
		},
		{
			name: "should cancel dust amount CCTX",
			cctx: sample.CrossChainTx(t, "0x123"),
			cctxModifier: func(cctx *crosschaintypes.CrossChainTx) {
				cctx.InboundParams.CoinType = coin.CoinType_Gas
				cctx.GetCurrentOutboundParam().Receiver = receiver.String()
				cctx.GetCurrentOutboundParam().ReceiverChainId = chain.ChainId
				cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(constant.BTCWithdrawalDustAmount - 1)
				cctx.GetCurrentOutboundParam().CallOptions.GasLimit = 254 // 254 bytes
				cctx.GetCurrentOutboundParam().GasPrice = "10"            // 10 sats/vByte
				cctx.GetCurrentOutboundParam().TssNonce = 1
			},
			chainID:     chain.ChainId,
			height:      101,
			minRelayFee: 0.00001, // 1000 sat/KB
			expected: &OutboundData{
				chainID:    chain.ChainId,
				to:         receiver,
				amount:     0, // should cancel the tx
				amountSats: 0,
				feeRate:    8, // Round(7.5)
				txSize:     254,
				nonce:      1,
				height:     101,
				cancelTx:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// modify cctx if needed
			if tt.cctxModifier != nil {
				tt.cctxModifier(tt.cctx)
			}

			outboundData, err := NewOutboundData(tt.cctx, tt.chainID, tt.height, tt.minRelayFee, log.Logger, log.Logger)
			if tt.errMsg != "" {
				require.Nil(t, outboundData)
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, outboundData)
			}
		})
	}
}
