package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/zeta-chain/node/x/fungible/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding fungible types.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		if !bytes.Equal(kvA.Key[:1], kvB.Key[:1]) {
			return fmt.Sprintf("key prefixes do not match. A: %X, B: %X", kvA.Key[:1], kvB.Key[:1])
		}
		switch {
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.SystemContractKey)):
			var systemContractA, systemContractB types.SystemContract
			if err := cdc.Unmarshal(kvA.Value, &systemContractA); err != nil {
				return fmt.Sprintf("failed to unmarshal A: %v", err)
			}
			if err := cdc.Unmarshal(kvB.Value, &systemContractB); err != nil {
				return fmt.Sprintf("failed to unmarshal B: %v", err)
			}
			return fmt.Sprintf("%v\n%v", systemContractA, systemContractB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ForeignCoinsKeyPrefix)):
			var foreignCoinsA, foreignCoinsB types.ForeignCoins
			if err := cdc.Unmarshal(kvA.Value, &foreignCoinsA); err != nil {
				return fmt.Sprintf("failed to unmarshal A: %v", err)
			}
			if err := cdc.Unmarshal(kvB.Value, &foreignCoinsB); err != nil {
				return fmt.Sprintf("failed to unmarshal B: %v", err)
			}
			return fmt.Sprintf("%v\n%v", foreignCoinsA, foreignCoinsB)
		default:
			return fmt.Sprintf("invalid fungible key prefix %X", kvA.Key[:1])
		}
	}
}
