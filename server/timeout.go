package server

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

// REF: https://github.com/zeta-chain/node/issues/4032
const (

	// DefaultTimeoutProposeDelta How much timeout_propose increases with each round
	DefaultTimeoutProposeDelta = 200 * time.Millisecond

	// DefaultTimeoutPrevoteDelta How much the timeout_prevote increases with each round
	DefaultTimeoutPrevoteDelta = 200 * time.Millisecond

	// DefaultTimeoutPrecommitDelta How much the timeout_precommit increases with each round
	DefaultTimeoutPrecommitDelta = 200 * time.Millisecond

	FlagSkipConfigOverwrite = "skip-config-overwrite"
)

// timeoutConfig represents a consensus timeout configuration pair
type timeoutConfig struct {
	currentValue *time.Duration
	defaultValue time.Duration
}

// overWriteConfigCmd overwrites the consensus timeout configurations to the default values.
func overWriteConfig(cmd *cobra.Command) error {
	serverCtx := server.GetServerContextFromCmd(cmd)

	timeoutConfigs := []timeoutConfig{
		{&serverCtx.Config.Consensus.TimeoutProposeDelta, DefaultTimeoutProposeDelta},
		{&serverCtx.Config.Consensus.TimeoutPrevoteDelta, DefaultTimeoutPrevoteDelta},
		{&serverCtx.Config.Consensus.TimeoutPrecommitDelta, DefaultTimeoutPrecommitDelta},
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
	}
	return nil
}
