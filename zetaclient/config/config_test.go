package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/config"
)

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
