package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSUICoinType(t *testing.T) {
	tests := []struct {
		name     string
		coinType CoinType
		want     bool
	}{
		{
			name:     "SUI coin type",
			coinType: SUI,
			want:     true,
		},
		{
			name:     "SUI short coin type",
			coinType: SUIShort,
			want:     true,
		},
		{
			name:     "not SUI coin type",
			coinType: "0xae330284eefb31f37777c78007fc3bc0e88ca69b1be267e1b803d37c9ea52dc6",
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsSUICoinType(test.coinType)
			require.Equal(t, test.want, got)
		})
	}
}
