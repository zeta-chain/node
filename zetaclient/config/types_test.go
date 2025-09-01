package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/config"
)

func Test_GetZetacoreClientConfig(t *testing.T) {
	signerName := "signerA"

	tests := []struct {
		name    string
		cfg     config.Config
		want    config.ZetacoreClientConfig
		wantErr bool
	}{
		{
			name: "should resolve zetacore client config from IP address",
			cfg: func() config.Config {
				cfg := config.New(false)
				cfg.ZetacoreIP = "127.0.0.1"
				cfg.AuthzHotkey = signerName
				return cfg
			}(),
			want: config.ZetacoreClientConfig{
				GRPCURL:     "127.0.0.1:9090",
				WSRemote:    "http://127.0.0.1:26657",
				SignerName:  signerName,
				GRPCDialOpt: config.CredsInsecureGRPC,
			},
		},
		{
			name: "should resolve zetacore client config from hostname",
			cfg: func() config.Config {
				cfg := config.New(false)
				cfg.ZetacoreURLGRPC = "zetachain.lavenderfive.com:443"
				cfg.ZetacoreURLWSS = "wss://rpc.lavenderfive.com:443/zetachain/websocket"
				cfg.AuthzHotkey = signerName
				return cfg
			}(),
			want: config.ZetacoreClientConfig{
				GRPCURL:     "zetachain.lavenderfive.com:443",
				WSRemote:    "https://rpc.lavenderfive.com:443/zetachain",
				SignerName:  signerName,
				GRPCDialOpt: config.CredsTLSGRPC,
			},
		},
		{
			name: "localnet zetacore container names should work",
			cfg: func() config.Config {
				cfg := config.New(false)
				cfg.ZetacoreIP = "zetacore0"
				cfg.AuthzHotkey = signerName
				return cfg
			}(),
			want: config.ZetacoreClientConfig{
				GRPCURL:     "zetacore0:9090",
				WSRemote:    "http://zetacore0:26657",
				SignerName:  signerName,
				GRPCDialOpt: config.CredsInsecureGRPC,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetZetacoreClientConfig()
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_GetRelayerKeyPath(t *testing.T) {
	// create config
	cfg := config.New(false)

	// should return default relayer key path
	require.Equal(t, config.DefaultRelayerKeyPath, cfg.GetRelayerKeyPath())
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
