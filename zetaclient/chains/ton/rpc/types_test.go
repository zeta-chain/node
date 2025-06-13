package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/tlb"
)

func TestAccountParsing(t *testing.T) {
	for _, tt := range []struct {
		name        string
		json        string
		errContains string
		assert      func(t *testing.T, acc Account)
	}{
		{
			name: "e2e-wallet-v3",
			json: `{
				"@type": "raw.fullAccountState",
				"balance": "1000000995239999",
				"extra_currencies": [],
				"code": "te6cckEBAQEAcQAA3v8AIN0gggFMl7ohggEznLqxn3Gw7UTQ0x/THzHXC//jBOCk8mCDCNcYINMf0x/TH/gjE7vyY+1E0NMf0x/T/9FRMrryoVFEuvKiBPkBVBBV+RDyo/gAkyDXSpbTB9QC+wDo0QGkyMsfyx/L/8ntVBC9ba0=",
				"data": "te6cckEBAQEAKgAAUAAAAAEAAAAqO9UYDUzSWygmFCZvXVUzhZdShgpP96QjaKV3FYP6dbx78kax",
				"last_transaction_id": {
					"@type": "internal.transactionId",
					"lt": "57000001",
					"hash": "cXi5gD4Z5fQdjqNskBXzfFDBLdtyJttcssdfhQqZq9c="
				},
				"block_id": {
					"@type": "ton.blockIdExt",
					"workchain": -1,
					"shard": "-9223372036854775808",
					"seqno": 539,
					"root_hash": "m2Jmjd7wYPOnJCqcrIi4jBNPfqjjiC80zWg5xyZKeTc=",
					"file_hash": "bAB0CxZ8q8yDdggk3ZknA2J5Yu6f4CVvJ8u3nOXhBZM="
				},
				"frozen_hash": "",
				"sync_utime": 1749116411,
				"@extra": "1749116849.21231:0:0.47571094006764125",
				"state": "active"
			}`,
			assert: func(t *testing.T, acc Account) {
				assert.Equal(t, tlb.AccountActive, acc.Status)
				assert.Equal(t, tlb.AccountActive, acc.ToShardAccount().Account.Status())
				assert.Equal(t, uint64(57000001), acc.LastTxLT)
				assert.Equal(t, uint64(1000000995239999), acc.Balance)
				assert.NotNil(t, acc.Code)
				assert.NotNil(t, acc.Data)
			},
		},
		{
			name: "unknown",
			json: `{
				"@type": "raw.fullAccountState",
				"balance": 0,
				"extra_currencies": [],
				"code": "",
				"data": "",
				"last_transaction_id": {
					"@type": "internal.transactionId",
					"lt": "0",
					"hash": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
				},
				"block_id": {
					"@type": "ton.blockIdExt",
					"workchain": -1,
					"shard": "-9223372036854775808",
					"seqno": 539,
					"root_hash": "m2Jmjd7wYPOnJCqcrIi4jBNPfqjjiC80zWg5xyZKeTc=",
					"file_hash": "bAB0CxZ8q8yDdggk3ZknA2J5Yu6f4CVvJ8u3nOXhBZM="
				},
				"frozen_hash": "",
				"sync_utime": 1749116411,
				"@extra": "1749117069.922797:0:0.6812065918617207",
				"state": "uninitialized"
			}`,
			assert: func(t *testing.T, acc Account) {
				assert.Equal(t, tlb.AccountUninit, acc.Status)
				assert.Equal(t, tlb.AccountUninit, acc.ToShardAccount().Account.Status())
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var acc Account

			err := json.Unmarshal([]byte(tt.json), &acc)
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
				return
			}

			require.NoError(t, err)
			tt.assert(t, acc)
		})
	}
}
