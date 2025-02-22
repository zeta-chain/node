package constant

import (
	"math/big"

	sdkmath "cosmossdk.io/math"

	"github.com/zeta-chain/node/pkg/chains"
)

var (
	// DefaultInboundFastConfirmationLiquidityMultiplier is the default ZRC20 liquidity cap multiplier for inbound fast confirmation
	// The reason of using decimal type is for easy integration into ChainParams, like field `ballot_threshold`
	DefaultInboundFastConfirmationLiquidityMultiplier = sdkmath.LegacyMustNewDecFromStr("0.00025")

	// InboundFastConfirmationLiquidityMultiplierMap maps chainID to ZRC20 liquidity cap multiplier for inbound fast confirmation.
	// Fast inbound confirmation is enabled only for chains explicitly listed in this map.
	InboundFastConfirmationLiquidityMultiplierMap = map[int64]sdkmath.LegacyDec{
		// Bitcoin
		chains.BitcoinMainnet.ChainId:       DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BitcoinSignetTestnet.ChainId: DefaultInboundFastConfirmationLiquidityMultiplier,
		chains.BitcoinTestnet4.ChainId:      DefaultInboundFastConfirmationLiquidityMultiplier,
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

// GetInboundFastConfirmationLiquidityMultiplier returns the ZRC20 liquidity cap multiplier for inbound fast confirmation.
//   - Fast confirmation applies to chains that use confirmation count (e.g. EVM chains and Bitcoin).
//   - The chains should have their liquidity cap multiplier set explicitly in the map above.
func GetInboundFastConfirmationLiquidityMultiplier(chainID int64) (sdkmath.LegacyDec, bool) {
	multiplier, enabled := InboundFastConfirmationLiquidityMultiplierMap[chainID]
	return multiplier, enabled
}

// CalcInboundFastAmountCap calculates the amount cap for inbound fast confirmation.
func CalcInboundFastAmountCap(liquidityCap sdkmath.Uint, multiplier sdkmath.LegacyDec) *big.Int {
	fastAmountCap := sdkmath.LegacyNewDecFromBigInt(liquidityCap.BigInt()).Mul(multiplier)
	return fastAmountCap.TruncateInt().BigInt()
}
