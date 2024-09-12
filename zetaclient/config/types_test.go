package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/config"
)

func Test_GetRelayerKeyPath(t *testing.T) {
	// create config
	cfg := config.New(false)

	// should return default relayer key path
	require.Equal(t, config.DefaultRelayerKeyPath, cfg.GetRelayerKeyPath())
}

func Test_GetEVMConfig(t *testing.T) {
	chain := chains.Sepolia
	chainID := chains.Sepolia.ChainId

	t.Run("should find non-empty evm config", func(t *testing.T) {
		// create config with defaults
		cfg := config.New(true)

		// set valid evm endpoint
		cfg.EVMChainConfigs[chainID] = config.EVMConfig{
			Chain:    chain,
			Endpoint: "localhost",
		}

		// should return non-empty evm config
		evmCfg, found := cfg.GetEVMConfig(chainID)
		require.True(t, found)
		require.False(t, evmCfg.Empty())
	})

	t.Run("should not find evm config if endpoint is empty", func(t *testing.T) {
		// create config with defaults
		cfg := config.New(true)

		// should not find evm config because endpoint is empty
		_, found := cfg.GetEVMConfig(chainID)
		require.False(t, found)
	})

	t.Run("should not find evm config if chain is empty", func(t *testing.T) {
		// create config with defaults
		cfg := config.New(true)

		// set empty chain
		cfg.EVMChainConfigs[chainID] = config.EVMConfig{
			Chain:    chains.Chain{},
			Endpoint: "localhost",
		}

		// should not find evm config because chain is empty
		_, found := cfg.GetEVMConfig(chainID)
		require.False(t, found)
	})
}

func Test_GetBTCConfig(t *testing.T) {
	tests := []struct {
		name     string
		chainID  int64
		chainCfg config.BTCConfig
		want     bool
	}{
		{
			name:    "should find non-empty btc config",
			chainID: chains.BitcoinRegtest.ChainId,
			chainCfg: config.BTCConfig{
				RPCUsername: "",
				RPCPassword: "",
				RPCHost:     "localhost",
				RPCParams:   "",
			},
			want: true,
		},
		{
			name:    "should not find btc config if rpc host is empty",
			chainID: chains.BitcoinRegtest.ChainId,
			chainCfg: config.BTCConfig{
				RPCUsername: "user",
				RPCPassword: "pass",
				RPCHost:     "",
				RPCParams:   "regtest",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create config with defaults
			cfg := config.New(true)

			// set chain config
			cfg.BTCChainConfigs[tt.chainID] = tt.chainCfg

			// should return btc config
			btcCfg, found := cfg.GetBTCConfig(tt.chainID)
			require.Equal(t, tt.want, found)
			require.Equal(t, tt.want, !btcCfg.Empty())
		})
	}
}
