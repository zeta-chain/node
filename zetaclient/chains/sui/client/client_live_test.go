package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/common"
)

const (
	RpcMainnet = "https://sui-mainnet.public.blastapi.io"
)

func TestClientLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		// todo
		// t.Skip("skipping live test")
		// return
	}

	t.Run("HealthCheck", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RpcMainnet)

		// ACT
		timestamp, err := ts.HealthCheck(ts.ctx)

		// ASSERT
		require.NoError(t, err)
		require.NotZero(t, timestamp)

		t.Logf("HealthCheck timestamp: %s (%s ago)", timestamp, time.Since(timestamp).String())
	})
}

type testSuite struct {
	t   *testing.T
	ctx context.Context
	*Client
}

func newTestSuite(t *testing.T, endpoint string) *testSuite {
	ctx := context.Background()
	client := NewFromEndpoint(endpoint)

	return &testSuite{t: t, ctx: ctx, Client: client}
}
