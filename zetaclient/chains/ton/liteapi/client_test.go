package liteapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashes(t *testing.T) {
	const sample = `48644940000001:e02b8c7cec103e08175ade8106619a8908707623c31451df2a68497c7d23d15a`

	lt, hash, err := TransactionHashFromString(sample)
	require.NoError(t, err)

	require.Equal(t, uint64(48644940000001), lt)
	require.Equal(t, "e02b8c7cec103e08175ade8106619a8908707623c31451df2a68497c7d23d15a", hash.Hex())
	require.Equal(t, sample, TransactionHashToString(lt, hash))
}

func TestTransactionHashFromString(t *testing.T) {
	for _, tt := range []struct {
		name  string
		raw   string
		error bool
		lt    uint64
		hash  string
	}{
		{
			name: "real example",
			raw:  "163000003:d0415f655644db6ee1260b1fa48e9f478e938823e8b293054fbae1f3511b77c5",
			lt:   163000003,
			hash: "d0415f655644db6ee1260b1fa48e9f478e938823e8b293054fbae1f3511b77c5",
		},
		{
			name: "zero lt",
			raw:  "0:0000000000000000000000000000000000000000000000000000000000000000",
			lt:   0,
			hash: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name: "big lt",
			raw:  "999999999999:fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0",
			lt:   999_999_999_999,
			hash: "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0",
		},
		{
			name:  "missing colon",
			raw:   "123456abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
			error: true,
		},
		{
			name:  "missing logical time",
			raw:   ":abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
			error: true,
		},
		{
			name:  "hash length",
			raw:   "123456:abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcde",
			error: true,
		},
		{
			name:  "non-numeric logical time",
			raw:   "notanumber:abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
			error: true,
		},
		{
			name:  "non-hex hash",
			raw:   "123456:xyz123xyz123xyz123xyz123xyz123xyz123xyz123xyz123xyz123xyz123xyz123",
			error: true,
		},
		{
			name:  "empty string",
			raw:   "",
			error: true,
		},
		{
			name:  "Invalid - only logical time, no hash",
			raw:   "123456:",
			error: true,
		},
		{
			name:  "Invalid - too many parts (extra colon)",
			raw:   "123456:abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef:extra",
			error: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			lt, hash, err := TransactionHashFromString(tt.raw)

			if tt.error {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.lt, lt)
			require.Equal(t, hash.Hex(), tt.hash)
		})
	}
}
