package ton

import (
	"cosmossdk.io/math"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/utils"
)

// Coins takes amount in TON and returns it in nano tons.
// Example Coins(5) return  math.Uint(5 * 10^9) nano tons.
//
//nolint:revive // in this context TON means 10^9 nano tons.
//goland:noinspection GoNameStartsWithPackageName
func Coins(amount uint64) math.Uint {
	// 1 ton = 10^9 nano tons
	const mul = 1_000_000_000

	return math.NewUint(amount).MulUint64(mul)
}

func UintToCoins(v math.Uint) tlb.Coins {
	return tlb.Coins(v.Uint64())
}

func FormatCoins(v math.Uint) string {
	// #nosec G115 always in range
	return utils.HumanFriendlyCoinsRepr(int64(v.Uint64()))
}
