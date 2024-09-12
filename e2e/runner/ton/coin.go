package ton

import (
	"cosmossdk.io/math"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/utils"
)

// TONCoins takes amount in nano tons and returns it in tons.
//
//nolint:revive // in this context TON means 10^9 nano tons.
//goland:noinspection GoNameStartsWithPackageName
func TONCoins(amount uint64) math.Uint {
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
