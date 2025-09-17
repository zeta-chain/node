package server

import (
	"fmt"
	"path/filepath"
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/chains"
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

// overWriteConfig overwrites the consensus timeout configurations to the default values.
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

		err = updateConfigFile(cmd, serverCtx.Config)
		if err != nil {
			return err
		}
	}
	return nil
}

// updateConfigFile updates the config file with the current server context configuration.
func updateConfigFile(cmd *cobra.Command, conf *cmtcfg.Config) error {
	rootDir, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil || rootDir == "" {
		fmt.Println("root directory :", rootDir)
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configPath := filepath.Join(rootDir, "config")
	cmtCfgFile := filepath.Join(configPath, "config.toml")
	cmtcfg.WriteConfigFile(cmtCfgFile, conf)
	return nil
}

// genesisChainID reads the genesis file at the given path and returns the corresponding chain ID in int64 format(EVM)
func genesisChainID(genesisFilePath string) (int64, error) {
	_, genesis, err := genutiltypes.GenesisStateFromGenFile(genesisFilePath)
	if err != nil {
		return -1, fmt.Errorf("failed to get genesis state from genesis file: %w", err)
	}
	evmChainID, err := chains.CosmosToEthChainID(genesis.ChainID)
	if err != nil {
		return -1, fmt.Errorf("failed to convert cosmos chain ID to ethereum chain ID: %w", err)
	}
	return evmChainID, nil
}
