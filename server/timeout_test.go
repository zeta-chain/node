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
