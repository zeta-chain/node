package client

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/zeta-chain/node/zetaclient/common"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

const (
	RPCMainnet = "https://sui-mainnet.public.blastapi.io"
	RPCTestnet = "https://sui-testnet.public.blastapi.io"
)

func TestClientLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("skipping live test")
		return
	}

	t.Run("HealthCheck", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCMainnet)

		// ACT
		timestamp, err := ts.HealthCheck(ts.ctx)

		// ASSERT
		require.NoError(t, err)
		require.NotZero(t, timestamp)

		t.Logf("HealthCheck timestamp: %s (%s ago)", timestamp, time.Since(timestamp).String())
	})

	t.Run("QueryEvents", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCMainnet)

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

	t.Run("GetOwnedObjectID", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCMainnet)

		// Given admin wallet us Cetus DEX team
		// (yeah, it took some time to find it)
		const ownerAddress = "0xdbfd0b17fa804c98f51d552b050fb7f850b85db96fa2a0d79e50119525814a47"

		// Given AdminCap struct type of Cetus DEX
		// (they use it for upgrades and stuff)
		const structType = "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb::config::AdminCap"

		// ACT
		// Get owned object id as we would fetch Gateway's WithdrawCap
		// that should belong to TSS
		objectID, err := ts.GetOwnedObjectID(ts.ctx, ownerAddress, structType)

		// ASSERT
		// https://suiscan.xyz/mainnet/object/0x89c1a321291d15ddae5a086c9abc533dff697fde3d89e0ca836c41af73e36a75
		require.NoError(t, err)
		require.Equal(t, "0x89c1a321291d15ddae5a086c9abc533dff697fde3d89e0ca836c41af73e36a75", objectID)
	})

	// examples taken from Cetus docs: https://cetus-1.gitbook.io/cetus-developer-docs/developer/via-contract/getting-started
	t.Run("CheckSharedObjects", func(t *testing.T) {
		ts := newTestSuite(t, RPCMainnet)

		// no object
		// all these objects are shared
		require.NoError(t, ts.CheckObjectIDsShared(ts.ctx, []string{}))

		// all these objects are shared
		objectIds := []string{
			"0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f", // Cetus global config
			"0x0000000000000000000000000000000000000000000000000000000000000006", // Sui universal clock object
			"0xf699e7f2276f5c9a75944b37a0c5b5d9ddfd2471bf6242483b03ab2887d198d0", // Cetus pool factory
		}
		require.NoError(t, ts.CheckObjectIDsShared(ts.ctx, objectIds))

		// contains a owned object
		objectIds = []string{
			"0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f",
			"0x6c31859275c1962b3e32bef11d9d60e7082eee86afe517e994685c62bc968082", // An owned NFT
			"0x0000000000000000000000000000000000000000000000000000000000000006",
			"0xf699e7f2276f5c9a75944b37a0c5b5d9ddfd2471bf6242483b03ab2887d198d0",
		}
		require.Error(t, ts.CheckObjectIDsShared(ts.ctx, objectIds))

		// contains a non existing object
		objectIds = []string{
			"0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f",
			"0x000000000000000000000000000000000000000000000000000000000000aaaa", // doesn't exist
			"0x0000000000000000000000000000000000000000000000000000000000000006",
			"0xf699e7f2276f5c9a75944b37a0c5b5d9ddfd2471bf6242483b03ab2887d198d0",
		}
		require.Error(t, ts.CheckObjectIDsShared(ts.ctx, objectIds))
	})

	t.Run("GetTransactionBlock successful tx", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCMainnet)

		// ACT
		res, err := ts.SuiGetTransactionBlock(ts.ctx, models.SuiGetTransactionBlockRequest{
			Digest:  "4PDngZHNfN79AvgB2VxNZcchDVdhNxNampVjFTFQUmzq",
			Options: models.SuiTransactionBlockOptions{ShowEffects: true},
		})

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, TxStatusSuccess, res.Effects.Status.Status)
	})

	t.Run("GetTransactionBlock failed tx", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCMainnet)

		// ACT
		res, err := ts.SuiGetTransactionBlock(ts.ctx, models.SuiGetTransactionBlockRequest{
			Digest:  "DUtYBP2UX4tFkXH1p4TWCCW2wkAbR6qPfvWsK55v5puq",
			Options: models.SuiTransactionBlockOptions{ShowEffects: true},
		})

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, TxStatusFailure, res.Effects.Status.Status)
	})

	t.Run("GetObjectParsedData", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCTestnet)

		// ACT
		// use testnet gateway object for testing
		gatewayID := "0x6fc08f682551e52c2cc34362a20f744ba6a3d8d17f6583fa2f774887c4079700"
		data, err := ts.GetObjectParsedData(ts.ctx, gatewayID)

		// ASSERT
		require.NoError(t, err)
		require.NotEmpty(t, data.Fields)
	})

	t.Run("GetObjectParsedData failed", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCTestnet)

		// ACT
		nonExistentID := "0x674d2b7396f2484dda53249ab5e4d4dee304e93a0037fd5d5d86aabd029fae98"
		data, err := ts.GetObjectParsedData(ts.ctx, nonExistentID)

		// ASSERT
		require.Error(t, err)
		require.Empty(t, data)
	})

	t.Run("GetSuiCoinObjectRefs", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, RPCTestnet)

		// Given TSS balance
		resp, err := ts.SuiXGetBalance(ts.ctx, models.SuiXGetBalanceRequest{
			Owner:    testutils.TSSAddressSuiTestnet,
			CoinType: string(sui.SUI),
		})
		require.NoError(t, err)

		tssBalance, err := strconv.ParseUint(resp.TotalBalance, 10, 64)
		require.NoError(t, err)
		require.Positive(t, tssBalance)

		// ACT-1
		// should be able to use all owned SUI coin objects
		coinRefs, err := ts.GetSuiCoinObjectRefs(ts.ctx, testutils.TSSAddressSuiTestnet, tssBalance)

		// ASSERT
		require.NoError(t, err)
		require.NotEmpty(t, coinRefs)

		// ACT-2
		// should NOT be able to cover the big amount (balance + 1)
		coinRefs, err = ts.GetSuiCoinObjectRefs(ts.ctx, testutils.TSSAddressSuiTestnet, tssBalance+1)

		// ASSERT
		require.ErrorContains(t, err, "SUI balance is too low")
		require.Empty(t, coinRefs)
	})

	t.Run("GetTransactionBlock successful tx on testnet with a deposit event", func(t *testing.T) {
		ts := newTestSuite(t, RPCTestnet)

		res, err := ts.SuiGetTransactionBlock(ts.ctx, models.SuiGetTransactionBlockRequest{
			Digest:  "BtVGRved1cvW3PHHeeMqeU96cwFxim5W6pNuHZpEuUQF",
			Options: models.SuiTransactionBlockOptions{ShowEvents: true, ShowEffects: true},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, TxStatusSuccess, res.Effects.Status.Status)

		gw, err := sui.NewGatewayFromAddress(
			"0x6b2fe12c605d64e14ca69f9aba51550593ba92ff43376d0a6cc26a5ca226f9bd,0x6fc08f682551e52c2cc34362a20f744ba6a3d8d17f6583fa2f774887c4079700",
		)
		require.NoError(t, err)

		require.Len(t, res.Events, 1)

		_, err = gw.ParseEvent(res.Events[0])
		require.NoError(t, err)
	})
}

func TestParseRPCResponse(t *testing.T) {
	// ARRANGE
	const raw = `{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
			"digest": "8995Wsnjv3udPYGgkfWhfNsu62W2UcT7Zd2tY83MDAyG",
			"confirmedLocalExecution": false
		}
	}`

	// ACT
	out, err := parseRPCResponse[models.SuiTransactionBlockResponse]([]byte(raw))

	// ASSERT
	require.NoError(t, err)
	require.Equal(t,
		models.SuiTransactionBlockResponse{
			Digest:                  "8995Wsnjv3udPYGgkfWhfNsu62W2UcT7Zd2tY83MDAyG",
			ConfirmedLocalExecution: false,
		},
		out,
	)

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
