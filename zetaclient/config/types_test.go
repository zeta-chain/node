package config_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/config"
)

func Test_ResolvePublicIP(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.Config
		expectIP string
		errorMsg string
	}{
		{
			name: "public IP is set",
			cfg: config.Config{
				PublicIP:  "127.0.0.1",
				PublicDNS: "my.zetaclient.com",
			},
			expectIP: "127.0.0.1",
		},
		{
			name:     "no public IP or DNS is set",
			cfg:      config.Config{},
			errorMsg: "no public IP or DNS is provided",
		},
		{
			name: "only public DNS is set",
			cfg: config.Config{
				// a real world example (Blockdaemon):
				// cjisrdnil456rejvnhd0.bdnodes.net => 202.8.10.137
				PublicDNS: "localhost",
			},
			expectIP: "127.0.0.1",
		},
		{
			name: "unable to resolve public DNS",
			cfg: config.Config{
				PublicDNS: "my.zetaclient.com",
			},
			errorMsg: "unable to resolve IP addresses for public DNS \"my.zetaclient.com\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publicIP, err := tt.cfg.ResolvePublicIP(zerolog.Nop())
			if tt.errorMsg != "" {
				require.Empty(t, publicIP)
				require.ErrorContains(t, err, tt.errorMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectIP, publicIP)
		})
	}
}

func Test_GetRelayerKeyPath(t *testing.T) {
	// create config
	cfg := config.New(false)

	// should return default relayer key path
	require.Equal(t, config.DefaultRelayerKeyPath, cfg.GetRelayerKeyPath())
}

func Test_GetMaxBaseFee(t *testing.T) {
	t.Run("should return zero max base fee", func(t *testing.T) {
		cfg := config.New(false)
		require.Zero(t, cfg.GetMaxBaseFee())
	})

	t.Run("should return configured max base fee", func(t *testing.T) {
		cfg := config.New(false)
		cfg.MaxBaseFee = 1000
		require.EqualValues(t, 1000, cfg.GetMaxBaseFee())
	})
}

func Test_GetMempoolCongestionThreshold(t *testing.T) {
	t.Run("should return zero mempool congestion threshold", func(t *testing.T) {
		cfg := config.New(false)
		require.Zero(t, cfg.GetMempoolCongestionThreshold())
	})

	t.Run("should return configured mempool congestion threshold", func(t *testing.T) {
		cfg := config.New(false)
		cfg.MempoolCongestionThreshold = 5000
		require.EqualValues(t, 5000, cfg.GetMempoolCongestionThreshold())
	})
}

func Test_GetEVMConfig(t *testing.T) {
	chainID := chains.Sepolia.ChainId

	t.Run("should find non-empty evm config", func(t *testing.T) {
		// create config with defaults
		cfg := config.New(true)

		// set valid evm endpoint
		cfg.EVMChainConfigs[chainID] = config.EVMConfig{
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
}

func Test_GetBTCConfig(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		oldCfg  config.BTCConfig
		btcCfg  *config.BTCConfig
		want    bool
	}{
		{
			name:    "should find non-empty btc config",
			chainID: chains.BitcoinRegtest.ChainId,
			btcCfg: &config.BTCConfig{
				RPCHost: "localhost",
			},
			want: true,
		},
		{
			name:    "should fallback to old config but still can't find btc config as it's empty",
			chainID: chains.BitcoinRegtest.ChainId,
			oldCfg: config.BTCConfig{
				RPCUsername: "user",
				RPCPassword: "pass",
				RPCHost:     "", // empty config
				RPCParams:   "regtest",
			},
			btcCfg: &config.BTCConfig{
				RPCUsername: "user",
				RPCPassword: "pass",
				RPCHost:     "", // empty config
				RPCParams:   "regtest",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create config with defaults
			cfg := config.New(true)

			if tt.btcCfg != nil {
				cfg.BTCChainConfigs[tt.chainID] = *tt.btcCfg
			}

			// should return btc config
			btcCfg, found := cfg.GetBTCConfig(tt.chainID)
			require.Equal(t, tt.want, found)
			require.Equal(t, tt.want, !btcCfg.Empty())
		})
	}
}

func Test_StringMasked(t *testing.T) {
	// create config with defaults
	cfg := config.New(true)

	cfg.SolanaConfig.Endpoint += "?api-key=123"

	// mask the config JSON string
	masked := cfg.StringMasked()
	require.NotEmpty(t, masked)

	// should contain necessary fields
	require.Contains(t, masked, "EVMChainConfigs")
	require.Contains(t, masked, "BTCChainConfigs")

	// should not contain endpoint
	require.NotContains(t, masked, "?api-key=123")
}
