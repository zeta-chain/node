package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/sdkconfig"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/config"
)

var sampleTestConfig = config.Config{
	KeyringBackend:     "test",
	ChainID:            "athens_7001-1",
	ZetaCoreURL:        "127.0.0.1",
	AuthzGranter:       "zeta1dkzcws63tttgd0alp6cesk2hlqagukauypc3qs",
	AuthzHotkey:        "hotkey",
	ConfigUpdateTicker: 6,
}

func TestValidate(t *testing.T) {
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
			errorMsg: "invalid public IP 192.168.1",
		},
		{
			name: "invalid public DNS name",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.PublicDNS = "invalid..dns"
				return cfg
			}(),
			errorMsg: "invalid public DNS invalid..dns",
		},
		{
			name: "invalid chain ID",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ChainID = "zeta1nvalid"
				return cfg
			}(),
			errorMsg: "invalid chain id zeta1nvalid",
		},
		{
			name: "invalid zetacore URL",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ZetaCoreURL = "     "
				return cfg
			}(),
			errorMsg: "invalid zetacore URL",
		},
		{
			name: "invalid granter address",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.AuthzGranter = "cosmos1dkzcws63tttgd0alp6cesk2hlqagukauypc3qs" // not ZetaChain address
				return cfg
			}(),
			errorMsg: "invalid bech32 granter address",
		},
		{
			name: "empty AuthzHotkey (grantee) name",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.AuthzHotkey = ""
				return cfg
			}(),
			errorMsg: "grantee name cannot be empty",
		},
		{
			name: "invalid log level",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.LogLevel = 6
				return cfg
			}(),
			errorMsg: "log level must be between 0 and 5",
		},
		{
			name: "invalid config update ticker",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.ConfigUpdateTicker = 0
				return cfg
			}(),
			errorMsg: "config update ticker cannot be 0",
		},
		{
			name: "invalid keyring backend",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.KeyringBackend = "invalid"
				return cfg
			}(),
			errorMsg: "invalid keyring backend invalid",
		},
		{
			name: "invalid max base fee",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.MaxBaseFee = -1
				return cfg
			}(),
			errorMsg: "max base fee cannot be negative",
		},
		{
			name: "invalid mempool congestion threshold",
			config: func() config.Config {
				cfg := sampleTestConfig
				cfg.MempoolCongestionThreshold = -1
				return cfg
			}(),
			errorMsg: "mempool congestion threshold cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.config)

			if tt.errorMsg != "" {
				require.ErrorContains(t, err, tt.errorMsg)
				return
			}
			require.NoError(t, err, "expected no error, got %v", err)
		})
	}
}

func Test_LoadRestrictedAddressesConfig(t *testing.T) {
	// Create test addresses
	testAddresses := []string{
		sample.RestrictedEVMAddressTest,
		sample.RestrictedBtcAddressTest,
		sample.RestrictedSolAddressTest,
		sample.RestrictedSuiAddressTest,
	}

	t.Run("should load restricted addresses from config file", func(t *testing.T) {
		// ARRANGE
		// create temporary directory
		basePath := sample.CreateTempDir(t)
		defer os.RemoveAll(basePath) // Clean up after test

		// create restricted addresses config file
		createRestrictedAddressesConfig(t, basePath, testAddresses)

		// ACT
		err := config.LoadRestrictedAddressesConfig(config.New(false), basePath)
		require.NoError(t, err)

		// ASSERT
		addresses := config.GetRestrictedAddresses()
		require.Equal(t, len(testAddresses), len(addresses))
		for _, addr := range testAddresses {
			require.True(t, slices.Contains(addresses, strings.ToLower(addr)))
		}
	})
}

// createRestrictedAddressesConfig creates a restricted addresses config file
func createRestrictedAddressesConfig(t *testing.T, basePath string, addresses []string) {
	// create config directory
	configDir := filepath.Join(basePath, "config")
	err := os.MkdirAll(configDir, 0o700)
	require.NoError(t, err)

	// marshal addresses and write to json file
	jsonData, err := json.Marshal(addresses)
	require.NoError(t, err)

	configFile := filepath.Join(configDir, "zetaclient_restricted_addresses.json")
	err = os.WriteFile(configFile, jsonData, 0600)
	require.NoError(t, err)

}
