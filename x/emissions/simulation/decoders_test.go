package simulation_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/emissions/simulation"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestDecodeStore(t *testing.T) {
	k, _, _, _ := keepertest.EmissionsKeeper(t)
	cdc := k.GetCodec()
	dec := simulation.NewDecodeStore(cdc)
	withdrawableEmissions := sample.WithdrawableEmissions(t)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.KeyPrefix(types.WithdrawableEmissionsKey), Value: cdc.MustMarshal(&withdrawableEmissions)},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{
			"withdrawable emissions",
			fmt.Sprintf(
				"key %s value A %v value B %v",
				types.WithdrawableEmissionsKey,
				withdrawableEmissions,
				withdrawableEmissions,
			),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]))
		})
	}
}
