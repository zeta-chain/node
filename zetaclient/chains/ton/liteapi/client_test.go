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
