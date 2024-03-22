package config

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.Equal(t, cfg.EthRPCURL, EthRPCURL)
	require.Equal(t, cfg.ZetaURL, ZetaURL)
	require.Equal(t, cfg.BtcExplorerURL, BtcExplorerURL)
	require.Equal(t, cfg.ConnectorAddress, ConnectorAddress)
	require.Equal(t, cfg.CustodyAddress, CustodyAddress)
}

func TestGetConfig(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defaultCfg := DefaultConfig()

	t.Run("No config file specified", func(t *testing.T) {
		cfg, err := GetConfig("")
		require.NoError(t, err)
		require.Equal(t, cfg, defaultCfg)

		exists, err := afero.Exists(AppFs, defaultCfgFileName)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("config file specified", func(t *testing.T) {
		cfg, err := GetConfig(defaultCfgFileName)
		require.NoError(t, err)
		require.Equal(t, cfg, defaultCfg)
	})
}

func TestConfig_Read(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	cfg, err := GetConfig("")
	require.NoError(t, err)

	t.Run("read existing file", func(t *testing.T) {
		c := &Config{}
		err := c.Read(defaultCfgFileName)
		require.NoError(t, err)
		require.Equal(t, c, cfg)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		err := AppFs.Remove(defaultCfgFileName)
		require.NoError(t, err)
		c := &Config{}
		err = c.Read(defaultCfgFileName)
		require.ErrorContains(t, err, "file does not exist")
		require.NotEqual(t, c, cfg)
	})
}

func TestConfig_Save(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	cfg := DefaultConfig()
	cfg.EtherscanAPIkey = "DIFFERENTAPIKEY"

	t.Run("save modified cfg", func(t *testing.T) {
		err := cfg.Save()
		require.NoError(t, err)

		newCfg, err := GetConfig(defaultCfgFileName)
		require.NoError(t, err)
		require.Equal(t, cfg, newCfg)
	})

	// Should test invalid json encoding but currently not able to without interface
}
