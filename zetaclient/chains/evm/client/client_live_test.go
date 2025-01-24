package client

import (
	"context"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/common"
)

const (
	URLEthMainnet     = "https://rpc.ankr.com/eth"
	URLEthSepolia     = "https://rpc.ankr.com/eth_sepolia"
	URLBscMainnet     = "https://rpc.ankr.com/bsc"
	URLPolygonMainnet = "https://rpc.ankr.com/polygon"
)

func TestLiveClient(t *testing.T) {
	if !common.LiveTestEnabled() {
		return
	}

	t.Run("IsTxConfirmed", func(t *testing.T) {
		ts := newTestSuite(t, URLEthMainnet)

		// check if the transaction is confirmed
		txHash := "0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146"

		t.Run("should confirm tx", func(t *testing.T) {
			confirmed, err := ts.IsTxConfirmed(ts.ctx, txHash, 12)
			require.NoError(t, err)
			require.True(t, confirmed)
		})

		t.Run("should not confirm tx if confirmations is not enough", func(t *testing.T) {
			confirmed, err := ts.IsTxConfirmed(ts.ctx, txHash, math.MaxUint64)
			require.NoError(t, err)
			require.False(t, confirmed)
		})
	})

	t.Run("HealthCheck", func(t *testing.T) {
		ts := newTestSuite(t, URLEthMainnet)

		_, err := ts.HealthCheck(ts.ctx)
		require.NoError(t, err)
	})
}

type testSuite struct {
	t *testing.T

	ctx       context.Context
	ethClient *ethclient.Client

	*Client
}

func newTestSuite(t *testing.T, endpoint string) *testSuite {
	ctx := context.Background()

	client, err := NewFromEndpoint(ctx, endpoint)
	require.NoError(t, err)

	return &testSuite{t: t, ctx: ctx, Client: client}
}
