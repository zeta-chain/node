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
		{"Regnet", BitcoinRegtest.ChainId, &chaincfg.RegressionNetParams, false},
		{"Mainnet", BitcoinMainnet.ChainId, &chaincfg.MainNetParams, false},
		{"Testnet", BitcoinTestnet.ChainId, &chaincfg.TestNet3Params, false},
		{"Signet", BitcoinSignetTestnet.ChainId, &chaincfg.SigNetParams, false},
		{"Testnet4", BitcoinTestnet4.ChainId, &TestNet4Params, false},
		{"Unknown", -1, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := BitcoinNetParamsFromChainID(tt.chainID)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, params)
			} else {
				require.NoError(t, err)
				require.EqualValues(t, tt.expected, params)
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
		{"Regnet", chaincfg.RegressionNetParams.Name, BitcoinRegtest.ChainId, false},
		{"Mainnet", chaincfg.MainNetParams.Name, BitcoinMainnet.ChainId, false},
		{"Testnet", chaincfg.TestNet3Params.Name, BitcoinTestnet.ChainId, false},
		{"Signet", chaincfg.SigNetParams.Name, BitcoinSignetTestnet.ChainId, false},
		{"Testnet4", TestNet4Params.Name, BitcoinTestnet4.ChainId, false},
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
	require.False(t, IsBitcoinRegnet(BitcoinSignetTestnet.ChainId))
	require.False(t, IsBitcoinRegnet(BitcoinTestnet4.ChainId))
}

func TestIsBitcoinMainnet(t *testing.T) {
	require.True(t, IsBitcoinMainnet(BitcoinMainnet.ChainId))
	require.False(t, IsBitcoinMainnet(BitcoinRegtest.ChainId))
	require.False(t, IsBitcoinMainnet(BitcoinTestnet.ChainId))
	require.False(t, IsBitcoinMainnet(BitcoinSignetTestnet.ChainId))
	require.False(t, IsBitcoinMainnet(BitcoinTestnet4.ChainId))
}
