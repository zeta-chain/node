package bitcoin

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/require"
)

func TestProve(t *testing.T) {
	t.Run("returns true if empty block", func(t *testing.T) {
		result := Prove(chainhash.Hash{}, chainhash.Hash{}, []byte{}, 0)
		require.True(t, result)
	})
}

func TestVerifyHash256Merkle(t *testing.T) {
	tests := []struct {
		name  string
		proof []byte
		index uint
		want  bool
	}{
		{
			name:  "valid length but invalid index and content",
			proof: make([]byte, 32),
			index: 0,
			want:  true,
		},
		{
			name:  "invalid length not multiple of 32",
			proof: make([]byte, 34),
			index: 0,
			want:  false,
		},
		{
			name:  "invalid length equal to 64",
			proof: make([]byte, 64),
			index: 0,
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := VerifyHash256Merkle(tc.proof, tc.index)
			require.Equal(t, tc.want, result)
		})
	}
}
