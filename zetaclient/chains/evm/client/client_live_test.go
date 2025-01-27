package client

import (
	"context"
	"math"
	"math/big"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/common"
)

const (
	URLEthMainnet     = "https://rpc.ankr.com/eth"
	URLEthSepolia     = "https://rpc.ankr.com/eth_sepolia"
	URLBscMainnet     = "https://rpc.ankr.com/bsc"
	URLPolygonMainnet = "https://rpc.ankr.com/polygon"
	URLBaseMainnet    = "https://rpc.ankr.com/base"
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

	// Note that you need a RPC with API key as const RPC rate limits are low
	t.Run("BlockByNumber2", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			endpoint string
			block    uint64
			assert   func(t *testing.T, v1 *geth.Block, errV1 error, v2 *Block, errV2 error)
		}{
			{
				name:     "Both work for ETH",
				endpoint: URLEthMainnet,
				block:    21718032,
				assert: func(t *testing.T, v1 *geth.Block, errV1 error, v2 *Block, errV2 error) {
					require.NoError(t, errV1)
					require.NoError(t, errV2)

					require.Equal(t, v1.Number().Int64(), int64(v2.Number))
					require.Equal(t, v1.Hash().String(), v2.Hash)
					require.Equal(t, v1.Transactions().Len(), len(v2.Transactions))
				},
			},
			{
				name:     "Only V2 works for BASE",
				endpoint: URLBaseMainnet,
				block:    25609323,
				assert: func(t *testing.T, v1 *geth.Block, errV1 error, v2 *Block, errV2 error) {
					require.ErrorContains(t, errV1, "transaction type not supported")
					require.NoError(t, errV2)

					require.Nil(t, v1)

					require.Equal(t, "0x93ec5d027f703544aa4b44892eacb2d58eb274cac45ed05694681e8e4d150aa6", v2.Hash)
					require.NotEmpty(t, v2.Transactions)
				},
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// ARRANGE
				ts := newTestSuite(t, tt.endpoint)

				// Given block number
				bn := big.NewInt(0).SetUint64(tt.block)

				// ACT
				v1, errV1 := ts.BlockByNumber(ts.ctx, bn)

				time.Sleep(1 * time.Second)

				v2, errV2 := ts.BlockByNumber2(ts.ctx, bn)

				// ASSERT
				tt.assert(t, v1, errV1, v2, errV2)
			})
		}
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
