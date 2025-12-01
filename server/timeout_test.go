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
				serverCtx.Config.Consensus.TimeoutProposeDelta = 500 * time.Millisecond
				serverCtx.Config.Consensus.TimeoutPrevoteDelta = 500 * time.Millisecond
				serverCtx.Config.Consensus.TimeoutPrecommitDelta = 500 * time.Millisecond
			},
			expectUpdate: true,
			expectValues: func(t *testing.T, ctx *server.Context) {
				require.Equal(t, DefaultTimeoutProposeDelta, ctx.Config.Consensus.TimeoutProposeDelta)
				require.Equal(t, DefaultTimeoutPrevoteDelta, ctx.Config.Consensus.TimeoutPrevoteDelta)
				require.Equal(t, DefaultTimeoutPrecommitDelta, ctx.Config.Consensus.TimeoutPrecommitDelta)
			},
		},
		{
			name: "no update needed - values already match defaults",
			setupTimeouts: func(serverCtx *server.Context) {
				serverCtx.Config.Consensus.TimeoutProposeDelta = DefaultTimeoutProposeDelta
				serverCtx.Config.Consensus.TimeoutPrevoteDelta = DefaultTimeoutPrevoteDelta
				serverCtx.Config.Consensus.TimeoutPrecommitDelta = DefaultTimeoutPrecommitDelta
			},
			expectUpdate: false,
			expectValues: func(t *testing.T, ctx *server.Context) {
				require.Equal(t, DefaultTimeoutProposeDelta, ctx.Config.Consensus.TimeoutProposeDelta)
				require.Equal(t, DefaultTimeoutPrevoteDelta, ctx.Config.Consensus.TimeoutPrevoteDelta)
				require.Equal(t, DefaultTimeoutPrecommitDelta, ctx.Config.Consensus.TimeoutPrecommitDelta)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cmd := StartCmd(nil, "")
			serverCtx := server.NewDefaultContext()
			tempDir := t.TempDir()
			serverCtx.Config.RootDir = tempDir
			ctx := context.WithValue(context.Background(), server.ServerContextKey, serverCtx)
			cmd.SetContext(ctx)

			genesis := sample.AppGenesis(t)
			rootConfigDir := filepath.Join(serverCtx.Config.RootDir, "config")
			require.NoError(t, os.MkdirAll(rootConfigDir, 0755))
			err := genesis.SaveAs(filepath.Join(rootConfigDir, "genesis.json"))
			require.NoError(t, err)

			tt.setupTimeouts(serverCtx)
			err = server.SetCmdServerContext(cmd, serverCtx)
			require.NoError(t, err)
			require.NoError(t, cmd.Flags().Set(flags.FlagHome, serverCtx.Config.RootDir))

			// Act
			err = overWriteConfig(cmd)

			// Assert
			require.NoError(t, err)
			updatedCtx := server.GetServerContextFromCmd(cmd)
			tt.expectValues(t, updatedCtx)
		})
	}
}

//
//func Test_UpdateConfigFile(t *testing.T) {
//	t.Run("should create config.toml in specified directory", func(t *testing.T) {
//		// Arrange
//		cmd := StartCmd(StartOptions{})
//		serverCtx := server.NewDefaultContext()
//
//		// Create a temporary directory for testing
//		tempDir := t.TempDir()
//		configDir := filepath.Join(tempDir, "config")
//		err := os.MkdirAll(configDir, 0755)
//		require.NoError(t, err)
//		require.NoError(t, cmd.Flags().Set(flags.FlagHome, tempDir))
//
//		// Act
//		err = updateConfigFile(cmd, serverCtx.Config)
//
//		// Assert
//		require.NoError(t, err)
//		configPath := filepath.Join(tempDir, "config", "config.toml")
//		_, err = os.Stat(configPath)
//		require.NoError(t, err, "config.toml file should exist")
//	})
//
//	t.Run("fail if home directory not set", func(t *testing.T) {
//		// Arrange
//		cmd := StartCmd(StartOptions{})
//		serverCtx := server.NewDefaultContext()
//
//		// Act
//		err := updateConfigFile(cmd, serverCtx.Config)
//
//		// Assert
//		require.Error(t, err)
//		require.Contains(t, err.Error(), "failed to get home directory")
//	})
//}
//
//func Test_GenesisChainID(t *testing.T) {
//	t.Run("get chainID for ZetaChainPrivnet", func(t *testing.T) {
//		genesis := sample.AppGenesis(t)
//		genesis.ChainID = fmt.Sprintf("test_%d-%d", chains.ZetaChainPrivnet.ChainId, 1)
//		tempDir := t.TempDir()
//		rootConfigDir := filepath.Join(tempDir, "config")
//		require.NoError(t, os.MkdirAll(rootConfigDir, 0755))
//		genesisFile := filepath.Join(rootConfigDir, "genesis.json")
//		err := genesis.SaveAs(genesisFile)
//		require.NoError(t, err)
//
//		id, err := genesisChainID(genesisFile)
//
//		require.NoError(t, err)
//		require.Equal(t, id, chains.ZetaChainPrivnet.ChainId)
//	})
//
//	t.Run("fail to get chainID if genesis file does not exist", func(t *testing.T) {
//		tempDir := t.TempDir()
//		genesisFile := filepath.Join(tempDir, "config", "genesis.json")
//
//		// Act
//		_, err := genesisChainID(genesisFile)
//
//		// Assert
//		require.Error(t, err)
//		require.Contains(t, err.Error(), "failed to get genesis state from genesis file")
//	})
//
//	t.Run("fail to get chain Id if genesis file has invalid chainID", func(t *testing.T) {
//		// Arrange
//		tempDir := t.TempDir()
//		rootConfigDir := filepath.Join(tempDir, "config")
//		require.NoError(t, os.MkdirAll(rootConfigDir, 0755))
//		genesisFile := filepath.Join(rootConfigDir, "genesis.json")
//
//		genesis := sample.AppGenesis(t)
//		genesis.ChainID = "invalid_chain_id"
//		err := genesis.SaveAs(genesisFile)
//		require.NoError(t, err)
//
//		// Act
//		_, err = genesisChainID(genesisFile)
//
//		// Assert
//		require.Error(t, err)
//		require.Contains(t, err.Error(), "failed to convert cosmos chain ID to ethereum chain ID")
//	})
//}
