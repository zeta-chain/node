package client

import (
	"context"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/common"
)

const (
	RpcMainnet = "https://sui-mainnet.public.blastapi.io"
)

func TestClientLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("skipping live test")
		return
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

	t.Run("QueryEvents", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RpcMainnet)

		// Sleep for a while to avoid rate limiting
		sleep := func() { time.Sleep(time.Second) }

		// Some aliases
		request := func(q EventQuery) ([]models.SuiEventResponse, string) {
			res, cursor, err := ts.QueryModuleEvents(ts.ctx, q)
			require.NoError(t, err, "Unable to query events. Query: %+v", q)

			return res, cursor
		}

		// Given some event query that searches for validator set events
		validatorQuery := func(cursor string, limit uint64) EventQuery {
			return EventQuery{
				PackageID: "0x3",
				Module:    "validator_set",
				Cursor:    cursor,
				Limit:     limit,
			}
		}

		eventsEqual := func(t *testing.T, a, b models.SuiEventResponse) {
			require.Equal(t, a.Id, b.Id)
			require.Equal(t, a.Bcs, b.Bcs)
			require.Equal(t, a.TimestampMs, b.TimestampMs)
		}

		// ACT
		// Let's query some validator events from RPC twice
		// First time, we'd query first 20 events
		res0, _ := request(validatorQuery("", 20))
		sleep()

		// Then we let's query 5 + 12 + 3 events
		res1, cursor1 := request(validatorQuery("", 5))
		sleep()

		res2, cursor2 := request(validatorQuery(cursor1, 12))
		sleep()

		res3, _ := request(validatorQuery(cursor2, 3))
		sleep()

		// ASSERT
		// We should have similar results combined
		resCombined := append(res1, append(res2, res3...)...)

		require.Equal(t, len(res0), 20)
		require.Equal(t, len(resCombined), 20)

		// Make sure that events are actually equal piece by piece
		for i, a := range res0 {
			eventsEqual(t, a, resCombined[i])
		}
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
