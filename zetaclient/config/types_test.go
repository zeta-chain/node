package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/config"
)

func Test_GetRelayerKeyPath(t *testing.T) {
	// create config
	cfg := config.New(false)

	// should return default relayer key path
	require.Equal(t, config.DefaultRelayerKeyPath, cfg.GetRelayerKeyPath())
}
