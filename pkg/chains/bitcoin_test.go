package chains

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func TestBitcoinNetParamsFromChainID(t *testing.T) {
	tests := []struct {
		name     string
		chainID  int64
		expected *chaincfg.Params
		wantErr  bool
	}{
		{"Regnet", BitcoinRegtest.ChainId, BitcoinRegnetParams, false},
		{"Mainnet", BitcoinMainnet.ChainId, BitcoinMainnetParams, false},
		{"Testnet", BitcoinTestnet.ChainId, BitcoinTestnetParams, false},
		{"Unknown", -1, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := BitcoinNetParamsFromChainID(tt.chainID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, params)
			}
		})
	}
}

func TestBitcoinChainIDFromNetParams(t *testing.T) {
	tests := []struct {
		name            string
		networkName     string
		expectedChainID int64
		wantErr         bool
	}{
		{"Regnet", BitcoinRegnetParams.Name, BitcoinRegtest.ChainId, false},
		{"Mainnet", BitcoinMainnetParams.Name, BitcoinMainnet.ChainId, false},
		{"Testnet", BitcoinTestnetParams.Name, BitcoinTestnet.ChainId, false},
		{"Unknown", "Unknown", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chainID, err := BitcoinChainIDFromNetworkName(tt.networkName)
			if tt.wantErr {
				require.Error(t, err)
				require.Zero(t, chainID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedChainID, chainID)
			}
		})
	}
}

func TestIsBitcoinRegnet(t *testing.T) {
	require.True(t, IsBitcoinRegnet(BitcoinRegtest.ChainId))
	require.False(t, IsBitcoinRegnet(BitcoinMainnet.ChainId))
	require.False(t, IsBitcoinRegnet(BitcoinTestnet.ChainId))
}

func TestIsBitcoinMainnet(t *testing.T) {
	require.True(t, IsBitcoinMainnet(BitcoinMainnet.ChainId))
	require.False(t, IsBitcoinMainnet(BitcoinRegtest.ChainId))
	require.False(t, IsBitcoinMainnet(BitcoinTestnet.ChainId))
}
