package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetae2e/local"
	"github.com/zeta-chain/node/e2e/config"
)

var initConf = config.DefaultConfig()
var configFile = ""

func NewInitCmd() *cobra.Command {
	var InitCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize config file for e2e tests",
		RunE:  initConfig,
	}

	InitCmd.Flags().StringVar(&initConf.RPCs.EVM, "ethURL", initConf.RPCs.EVM, "--ethURL http://eth:8545")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.ZetaCoreGRPC, "grpcURL", initConf.RPCs.ZetaCoreGRPC, "--grpcURL zetacore0:9090")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.ZetaCoreRPC, "rpcURL", initConf.RPCs.ZetaCoreRPC, "--rpcURL http://zetacore0:26657")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.Zevm, "zevmURL", initConf.RPCs.Zevm, "--zevmURL http://zetacore0:8545")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.Bitcoin.Host, "btcURL", initConf.RPCs.Bitcoin.Host, "--btcURL bitcoin:18443")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.Solana, "solanaURL", initConf.RPCs.Solana, "--solanaURL http://solana:8899")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.TON, "tonURL", initConf.RPCs.TON, "--tonURL http://ton:8081")
	InitCmd.Flags().StringVar(&initConf.ZetaChainID, "chainID", initConf.ZetaChainID, "--chainID athens_101-1")
	InitCmd.Flags().StringVar(&configFile, local.FlagConfigFile, "e2e.config", "--cfg ./e2e.config")

	return InitCmd
}

func initConfig(_ *cobra.Command, _ []string) error {
	err := initConf.GenerateKeys()
	if err != nil {
		return fmt.Errorf("generating keys: %w", err)
	}
	err = config.WriteConfig(configFile, initConf)
	if err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
