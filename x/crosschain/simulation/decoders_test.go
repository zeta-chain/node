package simulation_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/simulation"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestDecodeStore(t *testing.T) {
	k, _, _, _ := keepertest.CrosschainKeeper(t)
	cdc := k.GetCodec()
	dec := simulation.NewDecodeStore(cdc)
	cctx := sample.CrossChainTx(t, "sample")
	lastBlockHeight := sample.LastBlockHeight(t, "sample")
	gasPrice := sample.GasPrice(t, "sample")
	outboundTracker := sample.OutboundTracker(t, "sample")
	inboundTracker := sample.InboundTracker(t, "sample")
	zetaAccounting := sample.ZetaAccounting(t, "sample")
	rateLimiterFlags := sample.RateLimiterFlags()

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.KeyPrefix(types.CCTXKey), Value: cdc.MustMarshal(cctx)},
			{Key: types.KeyPrefix(types.LastBlockHeightKey), Value: cdc.MustMarshal(lastBlockHeight)},
			{Key: types.KeyPrefix(types.GasPriceKey), Value: cdc.MustMarshal(gasPrice)},
			{Key: types.KeyPrefix(types.OutboundTrackerKeyPrefix), Value: cdc.MustMarshal(&outboundTracker)},
			{Key: types.KeyPrefix(types.InboundTrackerKeyPrefix), Value: cdc.MustMarshal(&inboundTracker)},
			{Key: types.KeyPrefix(types.ZetaAccountingKey), Value: cdc.MustMarshal(&zetaAccounting)},
			{Key: types.KeyPrefix(types.RateLimiterFlagsKey), Value: cdc.MustMarshal(&rateLimiterFlags)},
			{Key: types.KeyPrefix(types.FinalizedInboundsKey), Value: []byte{1}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"CrossChainTx", fmt.Sprintf("key %s value A %v value B %v", types.CCTXKey, *cctx, *cctx)},
		{
			"LastBlockHeight",
			fmt.Sprintf("key %s value A %v value B %v", types.LastBlockHeightKey, *lastBlockHeight, *lastBlockHeight),
		},
		{"GasPrice", fmt.Sprintf("key %s value A %v value B %v", types.GasPriceKey, *gasPrice, *gasPrice)},
		{
			"OutboundTracker",
			fmt.Sprintf(
				"key %s value A %v value B %v",
				types.OutboundTrackerKeyPrefix,
				outboundTracker,
				outboundTracker,
			),
		},
		{
			"InboundTracker",
			fmt.Sprintf("key %s value A %v value B %v", types.InboundTrackerKeyPrefix, inboundTracker, inboundTracker),
		},
		{
			"ZetaAccounting",
			fmt.Sprintf("key %s value A %v value B %v", types.ZetaAccountingKey, zetaAccounting, zetaAccounting),
		},
		{
			"RateLimiterFlags",
			fmt.Sprintf("key %s value A %v value B %v", types.RateLimiterFlagsKey, rateLimiterFlags, rateLimiterFlags),
		},
		{
			"FinalizedInbounds",
			fmt.Sprintf("key %s value A %v value B %v", types.FinalizedInboundsKey, []byte{1}, []byte{1}),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]))
		})
	}
}
