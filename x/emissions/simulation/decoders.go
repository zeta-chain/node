package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/zeta-chain/node/x/emissions/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding crosschain types.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.WithdrawableEmissionsKey)):
			var withdrawableEmissionsA, withdrawableEmissionsB types.WithdrawableEmissions
			cdc.MustUnmarshal(kvA.Value, &withdrawableEmissionsA)
			cdc.MustUnmarshal(kvB.Value, &withdrawableEmissionsB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.WithdrawableEmissionsKey,
				withdrawableEmissionsA,
				withdrawableEmissionsB,
			)
		default:
			panic(
				fmt.Sprintf(
					"invalid emissions key prefix %X (first 8 bytes: %X)",
					kvA.Key[:1],
					kvA.Key[:min(8, len(kvA.Key))],
				),
			)
		}
	}
}
