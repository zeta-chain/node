package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/sdkconfig"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/mode"
)

func TestValidate(t *testing.T) {
	var sampleTestConfig = config.Config{
		KeyringBackend:     "test",
		ChainID:            "athens_7001-1",
		ZetaCoreURL:        "127.0.0.1",
		AuthzGranter:       "zeta1dkzcws63tttgd0alp6cesk2hlqagukauypc3qs",
		AuthzHotkey:        "hotkey",
		ConfigUpdateTicker: 6,
	}

	// set SDK config to use "zeta" address prefix
	sdkconfig.SetDefault(false)

	tests := []struct {
		name        string
		config      config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name:   "valid config",
			config: sampleTestConfig,
		},
		{
			name: "invalid public IP address",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.PublicIP = "192.168.1"
				return cfg
			}(),
			errorMsg: "reason: invalid public IP, got: 192.168.1",
		},
		{
			name: "public DNS is not supported",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.PublicDNS = "my.zetaclient.com"
				return cfg
			}(),
			errorMsg: "reason: public DNS is not supported, got: my.zetaclient.com",
		},
		{
			name: "invalid chain ID",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ChainID = "zeta1nvalid"
				return cfg
			}(),
			errorMsg: "reason: invalid chain id, got: zeta1nvalid",
		},
		{
			name: "invalid zetacore URL",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ZetaCoreURL = "     "
				return cfg
			}(),
			errorMsg: "reason: invalid zetacore URL, got:     ",
		},
		{
			name: "invalid granter address",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.AuthzGranter = "cosmos1dkzcws63tttgd0alp6cesk2hlqagukauypc3qs" // not ZetaChain address
				return cfg
			}(),
			errorMsg: "reason: invalid bech32 granter address, got: cosmos1dkzcws63tttgd0alp6cesk2hlqagukauypc3qs",
		},
		{
			name: "empty AuthzHotkey (grantee) name",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.AuthzHotkey = ""
				return cfg
			}(),
			errorMsg: "reason: grantee name is empty",
		},
		{
			name: "invalid log level",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.LogLevel = 6
				return cfg
			}(),
			errorMsg: "reason: log level must be between 0 and 5, got: 6",
		},
		{
			name: "invalid config update ticker",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ConfigUpdateTicker = 0
				return cfg
			}(),
			errorMsg: "reason: config update ticker is 0",
		},
		{
			name: "invalid keyring backend",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.KeyringBackend = "invalid"
				return cfg
			}(),
			errorMsg: "reason: invalid keyring backend, got: invalid",
		},
		{
			name: "invalid max base fee",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.MaxBaseFee = -1
				return cfg
			}(),
			errorMsg: "reason: max base fee cannot be negative, got: -1",
		},
		{
			name: "invalid mempool congestion threshold",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.MempoolCongestionThreshold = -1
				return cfg
			}(),
			errorMsg: "reason: mempool congestion threshold cannot be negative, got: -1",
		},
		{
			name: "empty ChaosProfilePath",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ClientMode = mode.ChaosMode
				return cfg
			}(),
			errorMsg: "ChaosProfilePath is a required field",
		},
		{
			name: "invalid ChaosProfilePath",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ClientMode = mode.ChaosMode
				cfg.ChaosProfilePath = "invalid/path"
				return cfg
			}(),
			errorMsg: `invalid ChaosProfilePath "invalid/path"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.errorMsg != "" {
				require.ErrorContains(t, err, tt.errorMsg)
				return
			}
			require.NoError(t, err, "expected no error, got %v", err)
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
