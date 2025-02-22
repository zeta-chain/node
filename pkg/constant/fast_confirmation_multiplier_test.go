package constant_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
)

func Test_GetInboundFastConfirmationLiquidityMultiplier(t *testing.T) {
	defaultMultiplier := constant.DefaultInboundFastConfirmationLiquidityMultiplier

	tests := []struct {
		name     string
		chainID  int64
		expected sdkmath.LegacyDec
	}{
		// Bitcoin
		{name: "Bitcoin mainnet", chainID: chains.BitcoinMainnet.ChainId, expected: defaultMultiplier},
		{name: "Bitcoin signet testnet", chainID: chains.BitcoinSignetTestnet.ChainId, expected: defaultMultiplier},
		{name: "Bitcoin testnet4", chainID: chains.BitcoinTestnet4.ChainId, expected: defaultMultiplier},
		{name: "Bitcoin regtest", chainID: chains.BitcoinRegtest.ChainId, expected: defaultMultiplier},

		// Ethereum
		{name: "Ethereum", chainID: chains.Ethereum.ChainId, expected: defaultMultiplier},
		{name: "Sepolia", chainID: chains.Sepolia.ChainId, expected: defaultMultiplier},
		{name: "Goerli localnet", chainID: chains.GoerliLocalnet.ChainId, expected: defaultMultiplier},

		// Binance Smart Chain
		{name: "Binance Smart Chain mainnet", chainID: chains.BscMainnet.ChainId, expected: defaultMultiplier},
		{name: "Binance Smart Chain testnet", chainID: chains.BscTestnet.ChainId, expected: defaultMultiplier},

		// Polygon
		{name: "Polygon", chainID: chains.Polygon.ChainId, expected: defaultMultiplier},
		{name: "Amoy", chainID: chains.Amoy.ChainId, expected: defaultMultiplier},

		// Base chain
		{name: "Base mainnet", chainID: chains.BaseMainnet.ChainId, expected: defaultMultiplier},
		{name: "Base Sepolia", chainID: chains.BaseSepolia.ChainId, expected: defaultMultiplier},

		// Optimism
		{name: "Optimism mainnet", chainID: chains.OptimismMainnet.ChainId, expected: defaultMultiplier},
		{name: "Optimism Sepolia", chainID: chains.OptimismSepolia.ChainId, expected: defaultMultiplier},

		// Avalanche
		{name: "Avalanche mainnet", chainID: chains.AvalancheMainnet.ChainId, expected: defaultMultiplier},
		{name: "Avalanche testnet", chainID: chains.AvalancheTestnet.ChainId, expected: defaultMultiplier},

		// Arbitrum
		{name: "Arbitrum mainnet", chainID: chains.ArbitrumMainnet.ChainId, expected: defaultMultiplier},
		{name: "Arbitrum Sepolia", chainID: chains.ArbitrumSepolia.ChainId, expected: defaultMultiplier},

		// World Chain
		{name: "World mainnet", chainID: chains.WorldMainnet.ChainId, expected: defaultMultiplier},
		{name: "World testnet", chainID: chains.WorldTestnet.ChainId, expected: defaultMultiplier},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, enabled := constant.GetInboundFastConfirmationLiquidityMultiplier(tt.chainID)
			if tt.expected.IsZero() {
				require.False(t, enabled)
				return
			}
			require.True(t, enabled)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func Test_CalcInboundFastAmountCap(t *testing.T) {
	tests := []struct {
		name         string
		liquidityCap sdkmath.Uint
		multiplier   sdkmath.LegacyDec
		expected     *big.Int
	}{
		{
			name:         "1% of 10000",
			liquidityCap: sdkmath.NewUintFromString("10000"),
			multiplier:   sdkmath.LegacyMustNewDecFromStr("0.01"),
			expected:     big.NewInt(100),
		},
		{
			name:         "0.15% of 10000",
			liquidityCap: sdkmath.NewUintFromString("10000"),
			multiplier:   sdkmath.LegacyMustNewDecFromStr("0.0015"),
			expected:     big.NewInt(15),
		},
		{
			name:         "0.025% of 10000",
			liquidityCap: sdkmath.NewUintFromString("10000"),
			multiplier:   sdkmath.LegacyMustNewDecFromStr("0.00025"),
			expected:     big.NewInt(2), // truncate 2.5 to 2
		},
		{
			name:         "0.0299% of 10000",
			liquidityCap: sdkmath.NewUintFromString("10000"),
			multiplier:   sdkmath.LegacyMustNewDecFromStr("0.000299"),
			expected:     big.NewInt(2), // truncate 2.99 to 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := constant.CalcInboundFastAmountCap(tt.liquidityCap, tt.multiplier)
			require.Equal(t, tt.expected, actual)
		})
	}
}
