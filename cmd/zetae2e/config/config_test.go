package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/config"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		expectNoError  bool
		expectNotEmpty bool
	}{
		{
			name:           "ReadLocalnetConfig",
			filePath:       "localnet.yml",
			expectNoError:  true,
			expectNotEmpty: true,
		},
		{
			name:           "ReadLocalConfig",
			filePath:       "local.yml",
			expectNoError:  true,
			expectNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, err := config.ReadConfig(tt.filePath, true)
			require.NoError(t, err)
			require.NotEmpty(t, conf.DefaultAccount.RawEVMAddress)
			require.NotEmpty(t, conf.DefaultAccount.RawPrivateKey)
		})
	}
}
