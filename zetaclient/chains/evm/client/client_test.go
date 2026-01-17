package client

import (
	"context"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/common"
)

// Note that you need a RPC with API key as const RPC rate limits are low
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

	t.Run("BlockByNumberCustom", func(t *testing.T) {
		for _, tt := range []struct {
			name        string
			endpoint    string
			block       uint64
			checkHeader bool
			assert      func(t *testing.T, v1 *geth.Block, errV1 error, v2 *Block, errV2 error)
		}{
			{
				name:     "both work for ETH",
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
				name:        "only v2 works for BASE",
				endpoint:    URLBaseMainnet,
				block:       25609323,
				checkHeader: true,
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

				v2, errV2 := ts.BlockByNumberCustom(ts.ctx, bn)

				// ASSERT
				tt.assert(t, v1, errV1, v2, errV2)

				if !tt.checkHeader {
					return
				}

				t.Run("HeaderByNumber works regardless of L2", func(t *testing.T) {
					// ARRANGE
					// RPC rate limit
					time.Sleep(time.Second)

					// ACT
					header, err := ts.HeaderByNumber(ts.ctx, bn)

					// ASSERT
					require.NoError(t, err)
					require.NotNil(t, header)
					require.Equal(t,
						strings.ToLower(v2.Hash),
						strings.ToLower(header.Hash().String()),
					)
				})
			})
		}
	})

	t.Run("TransactionByHashCustom", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			endpoint string
			txHash   string
			assert   func(t *testing.T, v1 *geth.Transaction, errV1 error, v2 *Transaction, errV2 error)
		}{
			{
				name:     "both work for BASE",
				endpoint: URLBaseMainnet,
				txHash:   "0xc2df77353c26eb282eb988a00f643e96d819fac20bb8bc1cfa3f4d6928be9fca",
				assert: func(t *testing.T, v1 *geth.Transaction, errV1 error, v2 *Transaction, errV2 error) {
					require.NoError(t, errV1)
					require.NoError(t, errV2)

					require.Equal(t, v1.Hash().String(), v2.Hash)

					v1To, v2To := strings.ToLower(v1.To().String()), strings.ToLower(v2.To)

					require.Equal(t, v1To, v2To)
					require.Equal(t, "0x17517d3645c28b4b6da6bc8cee382b16491605b2", v2To)
				},
			},
			{
				name:     "L1 deposit works only for BASE",
				endpoint: URLBaseMainnet,
				txHash:   "0xc373497bbbd9efdaa5c4b328950d148dfc71914a6180e4f3a9c70041334aa5f3",
				assert: func(t *testing.T, v1 *geth.Transaction, errV1 error, v2 *Transaction, errV2 error) {
					require.ErrorContains(t, errV1, "transaction type not supported")
					require.NoError(t, errV2)

					require.Nil(t, v1)
					require.Equal(t, "0x4200000000000000000000000000000000000007", v2.To)
				},
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// ASSERT
				ts := newTestSuite(t, tt.endpoint)

				// ACT
				v1, _, errV1 := ts.TransactionByHash(ts.ctx, gethcommon.HexToHash(tt.txHash))

				time.Sleep(1 * time.Second)

				v2, errV2 := ts.TransactionByHashCustom(ts.ctx, tt.txHash)

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

	ctx context.Context

	*Client
}

func newTestSuite(t *testing.T, endpoint string) *testSuite {
	ctx := context.Background()

	client, err := NewFromEndpoint(ctx, endpoint)
	require.NoError(t, err)

	return &testSuite{t: t, ctx: ctx, Client: client}
}
