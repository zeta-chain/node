package observer

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestDetermineCoinType(t *testing.T) {
	zetaTokenAddress := sample.EthAddress().Hex()
	tests := []struct {
		name     string
		asset    ethcommon.Address
		expected coin.CoinType
	}{
		{
			name:     "empty address returns Gas",
			asset:    ethcommon.Address{},
			expected: coin.CoinType_Gas,
		},
		{
			name:     "zeta token address returns Zeta",
			asset:    ethcommon.HexToAddress(zetaTokenAddress),
			expected: coin.CoinType_Zeta,
		},
		{
			name:     "other address returns ERC20",
			asset:    ethcommon.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
			expected: coin.CoinType_ERC20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, determineCoinType(tt.asset, zetaTokenAddress), "Coin type mismatch for %s", tt.name)
		})
	}
}
