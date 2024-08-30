package base_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/config"
)

func TestInitLogger(t *testing.T) {
	// test cases
	tests := []struct {
		name string
		t    *testing.T
		cfg  config.Config
		fail bool
	}{
		{
			name: "should be able to initialize json formatted logger",
			cfg: config.Config{
				LogFormat: "json",
				LogLevel:  1, // zerolog.InfoLevel,
				ComplianceConfig: config.ComplianceConfig{
					LogPath: sample.CreateTempDir(t),
				},
			},
			fail: false,
		},
		{
			name: "should be able to initialize plain text logger",
			cfg: config.Config{
				LogFormat: "text",
				LogLevel:  2, // zerolog.WarnLevel,
				ComplianceConfig: config.ComplianceConfig{
					LogPath: sample.CreateTempDir(t),
				},
			},
			fail: false,
		},
		{
			name: "should be able to initialize default formatted logger",
			cfg: config.Config{
				LogFormat: "unknown",
				LogLevel:  3, // zerolog.ErrorLevel,
				ComplianceConfig: config.ComplianceConfig{
					LogPath: sample.CreateTempDir(t),
				},
			},
			fail: false,
		},
		{
			name: "should be able to initialize sampled logger",
			cfg: config.Config{
				LogFormat:  "json",
				LogLevel:   4, // zerolog.DebugLevel,
				LogSampler: true,
				ComplianceConfig: config.ComplianceConfig{
					LogPath: sample.CreateTempDir(t),
				},
			},
		},
		{
			name: "should fail on invalid compliance log path",
			cfg: config.Config{
				LogFormat: "json",
				LogLevel:  1, // zerolog.InfoLevel,
				ComplianceConfig: config.ComplianceConfig{
					LogPath: "/invalid/123path",
				},
			},
			fail: true,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init logger
			logger, err := base.InitLogger(tt.cfg)

			// check if error is expected
			if tt.fail {
				require.Error(t, err)
				return
			}

			// check if logger is initialized
			require.NoError(t, err)

			// should be able to print log
			logger.Std.Info().Msg("print standard log")
			logger.Compliance.Info().Msg("print compliance log")
		})
	}
}
