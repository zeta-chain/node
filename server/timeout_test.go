package server

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_OverWriteConfig(t *testing.T) {
	tests := []struct {
		name          string
		setupTimeouts func(*server.Context)
		expectUpdate  bool
		expectValues  func(t *testing.T, ctx *server.Context)
	}{
		{
			name: "values different from defaults get updated",
			setupTimeouts: func(serverCtx *server.Context) {
				serverCtx.Config.Consensus.TimeoutPropose = 3 * time.Second
				serverCtx.Config.Consensus.TimeoutProposeDelta = 500 * time.Millisecond
				serverCtx.Config.Consensus.TimeoutPrevote = 1 * time.Second
				serverCtx.Config.Consensus.TimeoutPrevoteDelta = 500 * time.Millisecond
				serverCtx.Config.Consensus.TimeoutPrecommit = 1 * time.Second
				serverCtx.Config.Consensus.TimeoutPrecommitDelta = 500 * time.Millisecond
				serverCtx.Config.Consensus.TimeoutCommit = 3 * time.Second
			},
			expectUpdate: true,
			expectValues: func(t *testing.T, ctx *server.Context) {
				require.Equal(t, DefaultTimeoutPropose, ctx.Config.Consensus.TimeoutPropose)
				require.Equal(t, DefaultTimeoutProposeDelta, ctx.Config.Consensus.TimeoutProposeDelta)
				require.Equal(t, DefaultTimeoutPrevote, ctx.Config.Consensus.TimeoutPrevote)
				require.Equal(t, DefaultTimeoutPrevoteDelta, ctx.Config.Consensus.TimeoutPrevoteDelta)
				require.Equal(t, DefaultTimeoutPrecommit, ctx.Config.Consensus.TimeoutPrecommit)
				require.Equal(t, DefaultTimeoutPrecommitDelta, ctx.Config.Consensus.TimeoutPrecommitDelta)
				require.Equal(t, DefaultTimeoutCommit, ctx.Config.Consensus.TimeoutCommit)
			},
		},
		{
			name: "no update needed - values already match defaults",
			setupTimeouts: func(serverCtx *server.Context) {
				serverCtx.Config.Consensus.TimeoutPropose = DefaultTimeoutPropose
				serverCtx.Config.Consensus.TimeoutProposeDelta = DefaultTimeoutProposeDelta
				serverCtx.Config.Consensus.TimeoutPrevote = DefaultTimeoutPrevote
				serverCtx.Config.Consensus.TimeoutPrevoteDelta = DefaultTimeoutPrevoteDelta
				serverCtx.Config.Consensus.TimeoutPrecommit = DefaultTimeoutPrecommit
				serverCtx.Config.Consensus.TimeoutPrecommitDelta = DefaultTimeoutPrecommitDelta
				serverCtx.Config.Consensus.TimeoutCommit = DefaultTimeoutCommit
			},
			expectUpdate: false,
			expectValues: func(t *testing.T, ctx *server.Context) {
				require.Equal(t, DefaultTimeoutPropose, ctx.Config.Consensus.TimeoutPropose)
				require.Equal(t, DefaultTimeoutProposeDelta, ctx.Config.Consensus.TimeoutProposeDelta)
				require.Equal(t, DefaultTimeoutPrevote, ctx.Config.Consensus.TimeoutPrevote)
				require.Equal(t, DefaultTimeoutPrevoteDelta, ctx.Config.Consensus.TimeoutPrevoteDelta)
				require.Equal(t, DefaultTimeoutPrecommit, ctx.Config.Consensus.TimeoutPrecommit)
				require.Equal(t, DefaultTimeoutPrecommitDelta, ctx.Config.Consensus.TimeoutPrecommitDelta)
				require.Equal(t, DefaultTimeoutCommit, ctx.Config.Consensus.TimeoutCommit)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cmd := StartCmd(StartOptions{})
			serverCtx := server.NewDefaultContext()
			ctx := context.WithValue(context.Background(), server.ServerContextKey, serverCtx)
			cmd.SetContext(ctx)

			genesis := sample.AppGenesis(t)
			genesisFilePath := filepath.Join(serverCtx.Config.RootDir, "config", "genesis.json")
			err := genesis.SaveAs(genesisFilePath)
			require.NoError(t, err)

			tt.setupTimeouts(serverCtx)
			err = server.SetCmdServerContext(cmd, serverCtx)
			require.NoError(t, err)

			// Act
			err = overWriteConfig(cmd)

			// Assert
			require.NoError(t, err)
			updatedCtx := server.GetServerContextFromCmd(cmd)
			tt.expectValues(t, updatedCtx)
		})
	}
}

func Test_updateConfigFile(t *testing.T) {

	t.Run("should create config.toml in specified directory", func(t *testing.T) {
		// Arrange
		cmd := StartCmd(StartOptions{})
		serverCtx := server.NewDefaultContext()

		// Create a temporary directory for testing
		tempDir := t.TempDir()
		configDir := filepath.Join(tempDir, "config")
		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)
		cmd.Flags().Set(flags.FlagHome, tempDir)

		// Act
		err = updateConfigFile(cmd, serverCtx.Config)

		// Assert
		require.NoError(t, err)
		configPath := filepath.Join(tempDir, "config", "config.toml")
		_, err = os.Stat(configPath)
		require.NoError(t, err, "config.toml file should exist")
	})

	t.Run("fail if home directory not set", func(t *testing.T) {
		// Arrange
		cmd := StartCmd(StartOptions{})
		serverCtx := server.NewDefaultContext()

		// Act
		err := updateConfigFile(cmd, serverCtx.Config)

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get home directory")
	})
}
