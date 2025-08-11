package server

import (
	"fmt"
	"path/filepath"
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

// REF: https://github.com/zeta-chain/node/issues/4032
const (
	// DefaultTimeoutPropose How long we wait for a proposal block before prevoting nil
	DefaultTimeoutPropose = 3 * time.Second

	// DefaultTimeoutProposeDelta How much timeout_propose increases with each round
	DefaultTimeoutProposeDelta = 500 * time.Millisecond

	// DefaultTimeoutPrevote How long we wait after receiving +2/3 prevotes for "anything" (ie. not a single block or nil)
	DefaultTimeoutPrevote = 1 * time.Second

	// DefaultTimeoutPrevoteDelta How much the timeout_prevote increases with each round
	DefaultTimeoutPrevoteDelta = 500 * time.Millisecond

	// DefaultTimeoutPrecommit How long we wait after receiving +2/3 precommits for "anything" (ie. not a single block or nil)
	DefaultTimeoutPrecommit = 1 * time.Second

	// DefaultTimeoutPrecommitDelta How much the timeout_precommit increases with each round
	DefaultTimeoutPrecommitDelta = 500 * time.Millisecond

	// DefaultTimeoutCommit How long we wait after committing a block, before starting on the new
	// height (this gives us a chance to receive some more precommits, even
	// though we already have +2/3).
	DefaultTimeoutCommit = 3 * time.Second

	FlagSkipConfigOverride = "skip-config-override"
)

// timeoutConfig represents a consensus timeout configuration pair
type timeoutConfig struct {
	currentValue *time.Duration
	defaultValue time.Duration
}

func overRideConfigPreRunHandler(cmd *cobra.Command) error {
	serverCtx := server.GetServerContextFromCmd(cmd)

	timeoutConfigs := []timeoutConfig{
		{&serverCtx.Config.Consensus.TimeoutPropose, DefaultTimeoutPropose},
		{&serverCtx.Config.Consensus.TimeoutProposeDelta, DefaultTimeoutProposeDelta},
		{&serverCtx.Config.Consensus.TimeoutPrevote, DefaultTimeoutPrevote},
		{&serverCtx.Config.Consensus.TimeoutPrevoteDelta, DefaultTimeoutPrevoteDelta},
		{&serverCtx.Config.Consensus.TimeoutPrecommit, DefaultTimeoutPrecommit},
		{&serverCtx.Config.Consensus.TimeoutPrecommitDelta, DefaultTimeoutPrecommitDelta},
		{&serverCtx.Config.Consensus.TimeoutCommit, DefaultTimeoutCommit},
	}

	needsUpdate := false
	for _, config := range timeoutConfigs {
		if *config.currentValue != config.defaultValue {
			*config.currentValue = config.defaultValue
			needsUpdate = true
		}
	}

	if needsUpdate {

		err := server.SetCmdServerContext(cmd, serverCtx)
		if err != nil {
			return fmt.Errorf("failed to set server context: %w", err)
		}

		err = updateConfigFile(cmd, serverCtx.Config)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateConfigFile(cmd *cobra.Command, conf *cmtcfg.Config) error {
	rootDir, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configPath := filepath.Join(rootDir, "config")
	cmtCfgFile := filepath.Join(configPath, "config.toml")

	cmtcfg.WriteConfigFile(cmtCfgFile, conf)
	fmt.Printf("Consensus timeouts updated in %s\n", cmtCfgFile)
	return nil
}
