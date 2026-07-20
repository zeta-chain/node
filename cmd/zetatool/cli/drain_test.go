package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/drain"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestLatestTSS(t *testing.T) {
	t.Run("picks highest finalized height", func(t *testing.T) {
		list := []observertypes.TSS{
			{TssPubkey: "a", FinalizedZetaHeight: 100},
			{TssPubkey: "c", FinalizedZetaHeight: 300},
			{TssPubkey: "b", FinalizedZetaHeight: 200},
		}
		got, err := latestTSS(list)
		require.NoError(t, err)
		require.Equal(t, "c", got.TssPubkey)
	})

	t.Run("errors on empty", func(t *testing.T) {
		_, err := latestTSS(nil)
		require.Error(t, err)
	})
}

func TestPickMedian(t *testing.T) {
	t.Run("returns median-indexed price", func(t *testing.T) {
		gp := &crosschaintypes.GasPrice{Prices: []uint64{10, 20, 30}, MedianIndex: 1}
		got, err := pickMedian(gp, 1)
		require.NoError(t, err)
		require.Equal(t, "20", got.String())
	})

	t.Run("errors on nil / empty / out-of-range", func(t *testing.T) {
		_, err := pickMedian(nil, 1)
		require.Error(t, err)
		_, err = pickMedian(&crosschaintypes.GasPrice{}, 1)
		require.Error(t, err)
		_, err = pickMedian(&crosschaintypes.GasPrice{Prices: []uint64{1}, MedianIndex: 5}, 1)
		require.Error(t, err)
	})
}

func TestDrainNetwork(t *testing.T) {
	require.Equal(t, drain.NetworkMainnet, drainNetwork(config.NetworkMainnet))
	require.Equal(t, drain.NetworkLocalnet, drainNetwork(config.NetworkLocalnet))
	require.Equal(t, drain.NetworkTestnet, drainNetwork(config.NetworkTestnet))
	require.Equal(t, drain.NetworkTestnet, drainNetwork(config.NetworkSignet))
}
