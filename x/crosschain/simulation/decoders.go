package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/zeta-chain/node/x/crosschain/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding crosschain types.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.CCTXKey)):
			var cctxA, cctxB types.CrossChainTx
			cdc.MustUnmarshal(kvA.Value, &cctxA)
			cdc.MustUnmarshal(kvB.Value, &cctxB)
			return fmt.Sprintf("key %s value A %v value B %v", types.CCTXKey, cctxA, cctxB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.FinalizedInboundsKey)):
			var finalizedInboundsA, finalizedInboundsB []byte
			finalizedInboundsA = kvA.Value
			finalizedInboundsB = kvB.Value
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.FinalizedInboundsKey,
				finalizedInboundsA,
				finalizedInboundsB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.GasPriceKey)):
			var gasPriceA, gasPriceB types.GasPrice
			cdc.MustUnmarshal(kvA.Value, &gasPriceA)
			cdc.MustUnmarshal(kvB.Value, &gasPriceB)
			return fmt.Sprintf("key %s value A %v value B %v", types.GasPriceKey, gasPriceA, gasPriceB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.OutboundTrackerKeyPrefix)):
			var outboundTrackerA, outboundTrackerB types.OutboundTracker
			cdc.MustUnmarshal(kvA.Value, &outboundTrackerA)
			cdc.MustUnmarshal(kvB.Value, &outboundTrackerB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.OutboundTrackerKeyPrefix,
				outboundTrackerA,
				outboundTrackerB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.InboundTrackerKeyPrefix)):
			var inboundTrackerA, inboundTrackerB types.InboundTracker
			cdc.MustUnmarshal(kvA.Value, &inboundTrackerA)
			cdc.MustUnmarshal(kvB.Value, &inboundTrackerB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.InboundTrackerKeyPrefix,
				inboundTrackerA,
				inboundTrackerB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ZetaAccountingKey)):
			var zetaAccountingA, zetaAccountingB types.ZetaAccounting
			cdc.MustUnmarshal(kvA.Value, &zetaAccountingA)
			cdc.MustUnmarshal(kvB.Value, &zetaAccountingB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.ZetaAccountingKey,
				zetaAccountingA,
				zetaAccountingB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.RateLimiterFlagsKey)):
			var rateLimiterFlagsA, rateLimiterFlagsB types.RateLimiterFlags
			cdc.MustUnmarshal(kvA.Value, &rateLimiterFlagsA)
			cdc.MustUnmarshal(kvB.Value, &rateLimiterFlagsB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.RateLimiterFlagsKey,
				rateLimiterFlagsA,
				rateLimiterFlagsB,
			)
		default:
			panic(
				fmt.Sprintf(
					"invalid crosschain key prefix %X (first 8 bytes: %X)",
					kvA.Key[:1],
					kvA.Key[:min(8, len(kvA.Key))],
				),
			)
		}
	}
}
