package chains

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCosmosToEthChainID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		chainID  string
		expected int64
		isErr    bool
	}{
		{
			name:     "valid chain ID",
			chainID:  "cosmoshub_400-1",
			expected: 400,
		},
		{
			name:     "another valid chain ID",
			chainID:  "athens_701-2",
			expected: 701,
		},
		{
			name:    "no underscore",
			chainID: "athens701-2",
			isErr:   true,
		},
		{
			name:    "no dash",
			chainID: "athens_7012",
			isErr:   true,
		},
		{
			name:    "wrong pattern",
			chainID: "athens-701_2",
			isErr:   true,
		},
		{
			name:    "invalid number format",
			chainID: "cosmoshub_abc-1",
			isErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ethChainID, err := CosmosToEthChainID(tc.chainID)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, ethChainID)
			}
		})
	}
}
