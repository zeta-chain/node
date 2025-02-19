package constant

import (
	sdkmath "cosmossdk.io/math"

	"github.com/zeta-chain/node/pkg/chains"
)

var (
	// DefaultInboundFastConfirmationLiquidityMultiplier is the default ZRC20 liquidity cap multiplier for inbound fast confirmation
	DefaultInboundFastConfirmationLiquidityMultiplier = sdkmath.LegacyMustNewDecFromStr("0.0001")

	// InboundFastConfirmationLiquidityMultiplierMap maps chainID to ZRC20 liquidity cap multiplier for inbound fast confirmation.
	// Fast inbound confirmation is enabled only for chains explicitly listed in this map.
	InboundFastConfirmationLiquidityMultiplierMap = map[int64]sdkmath.LegacyDec{
		// Bitcoin
		chains.BitcoinMainnet.ChainId:       DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BitcoinSignetTestnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BitcoinTestnet.ChainId:       DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BitcoinRegtest.ChainId:       DefaultInboundFastConfirmationLiquidityMultiplier,

		// Ethereum
		chains.Ethereum.ChainId:       DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.Sepolia.ChainId:        DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.GoerliLocalnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,

		// Binance Smart Chain
		chains.BscMainnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BscTestnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,

		// Polygon
		chains.Polygon.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.Amoy.ChainId:    DefaultInboundFastConfirmationLiquidityMultiplier,

		// Base chain
		chains.BaseMainnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BaseSepolia.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,

		// Optimism
		chains.OptimismMainnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.OptimismSepolia.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,

		// Avalanche
		chains.AvalancheMainnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.AvalancheTestnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,

		// Arbitrum
		chains.ArbitrumMainnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.ArbitrumSepolia.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,

		// World Chain
		chains.WorldMainnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.WorldTestnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
	}
)
