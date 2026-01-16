package base

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mathpkg "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_NewTSSKeysignInfo(t *testing.T) {
	tests := []struct {
		name      string
		digest    []byte
		signature [65]byte
	}{
		{
			name:      "new keysign info",
			digest:    []byte{1, 2, 3},
			signature: sample.Signature65B(t),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			info := NewTSSKeysignInfo(tc.digest, tc.signature)
			require.NotNil(t, info)
			require.Equal(t, tc.digest, info.digest)
			require.Equal(t, tc.signature, info.signature)
		})
	}
}

func Test_NewTSSKeysignBatch(t *testing.T) {
	batch := NewTSSKeysignBatch()
	require.NotNil(t, batch)
	require.NotNil(t, batch.Digests())
	require.Empty(t, batch.Digests())
	require.Zero(t, batch.NonceLow())
	require.Zero(t, batch.NonceHigh())
}

func Test_BatchNumber(t *testing.T) {
	tests := []struct {
		name      string
		nonce     uint64
		expected  uint64
		setupFunc func(*TSSKeysignBatch)
	}{
		{
			name:     "batch 0, nonce 9",
			nonce:    9,
			expected: 0,
		},
		{
			name:     "batch 1, nonce 10",
			nonce:    10,
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			info := NewTSSKeysignInfo(sample.Digest32B(t), sample.Signature65B(t))
			batch := NewTSSKeysignBatch()
			batch.AddKeysignInfo(tc.nonce, *info)

			require.Equal(t, tc.expected, batch.BatchNumber())
		})
	}
}

func Test_AddKeysignInfo(t *testing.T) {
	// sample keysign infos
	infoList := createKeysignInfoList(t, 4)

	tests := []struct {
		name          string
		nonces        []uint64
		expectedCount int
		expectedLow   uint64
		expectedHigh  uint64
	}{
		{
			name:          "empty batch",
			nonces:        []uint64{},
			expectedLow:   0,
			expectedHigh:  0,
			expectedCount: 0,
		},
		{
			name:          "single nonce",
			nonces:        []uint64{1},
			expectedLow:   1,
			expectedHigh:  1,
			expectedCount: 1,
		},
		{
			name:          "multiple nonces",
			nonces:        []uint64{0, 1, 2, 3},
			expectedLow:   0,
			expectedHigh:  3,
			expectedCount: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			batch := NewTSSKeysignBatch()
			for i, nonce := range tc.nonces {
				info := infoList[i]
				batch.AddKeysignInfo(nonce, *info)
			}

			require.Len(t, batch.Digests(), tc.expectedCount)
			require.Equal(t, tc.expectedLow, batch.NonceLow())
			require.Equal(t, tc.expectedHigh, batch.NonceHigh())
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	batch := NewTSSKeysignBatch()
	require.True(t, batch.IsEmpty())
}

func Test_IsSequential(t *testing.T) {
	infoList := createKeysignInfoList(t, 10)

	tests := []struct {
		name     string
		nonces   []uint64
		expected bool
	}{
		{
			name:     "empty batch is not sequential",
			nonces:   []uint64{},
			expected: false,
		},
		{
			name:     "single nonce is sequential",
			nonces:   []uint64{1},
			expected: true,
		},
		{
			name:     "sequential nonces [1,2,3,4]",
			nonces:   []uint64{1, 2, 3, 4},
			expected: true,
		},
		{
			name:     "sequential nonces [0,1,3,4]",
			nonces:   []uint64{0, 1, 3, 4},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			batch := NewTSSKeysignBatch()
			for i, nonce := range tc.nonces {
				info := infoList[i]
				batch.AddKeysignInfo(nonce, *info)
			}
			require.Equal(t, tc.expected, batch.IsSequential())
		})
	}
}

func Test_IsEnd(t *testing.T) {
	infoList := createKeysignInfoList(t, 20)

	tests := []struct {
		name     string
		nonces   []uint64
		expected bool
	}{
		{
			name:     "batch 0, nonce 9 (end of batch)",
			nonces:   []uint64{9},
			expected: true,
		},
		{
			name:     "batch 0, nonce [7,8,9] (end of batch)",
			nonces:   []uint64{7, 8, 9},
			expected: true,
		},
		{
			name:     "batch 0, nonce [6,7,8] (not end of batch)",
			nonces:   []uint64{6, 7, 8},
			expected: false,
		},
		{
			name:     "batch 1, nonce [18, 19] (end of batch)",
			nonces:   []uint64{18, 19},
			expected: true,
		},
		{
			name:     "batch 1, nonce [14,15] (not end of batch)",
			nonces:   []uint64{14, 15},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			batch := NewTSSKeysignBatch()
			for i, nonce := range tc.nonces {
				info := infoList[i]
				batch.AddKeysignInfo(nonce, *info)
			}
			require.Equal(t, tc.expected, batch.IsEnd())
		})
	}
}

func Test_ContainsNonce(t *testing.T) {
	infoList := createKeysignInfoList(t, 10)

	// ARRANGE
	batch := NewTSSKeysignBatch()
	for i, info := range infoList {
		// #nosec G115 - always positive
		nonce := uint64(i)
		batch.AddKeysignInfo(nonce, *info)
	}

	// ASSERT
	for i := range infoList {
		// #nosec G115 - always positive
		nonce := uint64(i)
		require.True(t, batch.ContainsNonce(nonce))
	}
}

func Test_KeysignHeight(t *testing.T) {
	tests := []struct {
		name       string
		chainID    int64
		zetaHeight int64
		errorMsg   string
	}{
		{
			name:       "max zeta height bucket and chainID, should work",
			chainID:    mathpkg.MaxPairValue,
			zetaHeight: mathpkg.MaxPairValue*10 - 1,
		},
		{
			name:       "invalid zeta height zero",
			chainID:    1,
			zetaHeight: 0,
			errorMsg:   fmt.Sprintf("invalid zeta height: %d", 0),
		},
		{
			name:       "invalid zeta height too large",
			chainID:    1,
			zetaHeight: mathpkg.MaxPairValue * 10,
			errorMsg:   fmt.Sprintf("invalid zeta height: %d", mathpkg.MaxPairValue*10),
		},
		{
			name:       "invalid chainID zero",
			chainID:    0,
			zetaHeight: 1,
			errorMsg:   "invalid chain ID: 0",
		},
		{
			name:       "invalid chainID too large",
			chainID:    mathpkg.MaxPairValue + 1,
			zetaHeight: 1,
			errorMsg:   fmt.Sprintf("invalid chain ID: %d", mathpkg.MaxPairValue+1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			height, err := KeysignHeight(tc.chainID, tc.zetaHeight)

			// ASSERT
			if tc.errorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				require.Positive(t, height)
			}
		})
	}
}

func Test_NonceToBatchNumber(t *testing.T) {
	tests := []struct {
		name     string
		nonce    uint64
		expected uint64
	}{
		{
			name:     "nonce 0 -> batch 0",
			nonce:    0,
			expected: 0,
		},
		{
			name:     "nonce 9 -> batch 0",
			nonce:    9,
			expected: 0,
		},
		{
			name:     "nonce 10 -> batch 1",
			nonce:    10,
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NonceToBatchNumber(tc.nonce)
			require.Equal(t, tc.expected, result)
		})
	}
}

func Test_BatchNumberToRange(t *testing.T) {
	tests := []struct {
		name         string
		batchNumber  uint64
		expectedLow  uint64
		expectedHigh uint64
	}{
		{
			name:         "batch 0 -> [0, 9]",
			batchNumber:  0,
			expectedLow:  0,
			expectedHigh: 9,
		},
		{
			name:         "batch 1 -> [10, 19]",
			batchNumber:  1,
			expectedLow:  10,
			expectedHigh: 19,
		},
		{
			name:         "batch 9 -> [90, 99]",
			batchNumber:  9,
			expectedLow:  90,
			expectedHigh: 99,
		},
		{
			name:         "batch 10 -> [100, 109]",
			batchNumber:  10,
			expectedLow:  100,
			expectedHigh: 109,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			low, high := BatchNumberToRange(tc.batchNumber)
			require.Equal(t, tc.expectedLow, low)
			require.Equal(t, tc.expectedHigh, high)
		})
	}
}

// createKeysignInfoList creates a list of sample keysign infos.
func createKeysignInfoList(t *testing.T, count int) []*TSSKeysignInfo {
	infos := make([]*TSSKeysignInfo, count)
	for i := range count {
		infos[i] = NewTSSKeysignInfo(sample.Digest32B(t), [65]byte{})
	}
	return infos
}
