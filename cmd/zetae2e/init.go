package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd/zetae2e/local"
	"github.com/zeta-chain/zetacore/e2e/config"
)

var initConf = config.Config{}
var configFile = ""

func NewInitCmd() *cobra.Command {
	var InitCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize config file for e2e tests",
		Run:   initConfig,
	}

	InitCmd.Flags().StringVar(&initConf.RPCs.EVM, "ethURL", "http://eth:8545", "--ethURL http://eth:8545")
	InitCmd.Flags().StringVar(&initConf.RPCs.ZetaCoreGRPC, "grpcURL", "zetacore0:9090", "--grpcURL zetacore0:9090")
	InitCmd.Flags().StringVar(&initConf.RPCs.ZetaCoreRPC, "rpcURL", "http://zetacore0:26657", "--rpcURL http://zetacore0:26657")
	InitCmd.Flags().StringVar(&initConf.RPCs.Zevm, "zevmURL", "http://zetacore0:8545", "--zevmURL http://zetacore0:8545")
	InitCmd.Flags().StringVar(&initConf.RPCs.Bitcoin.Host, "btcURL", "bitcoin:18443", "--grpcURL bitcoin:18443")

	InitCmd.Flags().StringVar(&initConf.ZetaChainID, "chainID", "athens_101-1", "--chainID athens_101-1")
	InitCmd.Flags().StringVar(&configFile, local.FlagConfigFile, "e2e.config", "--cfg ./e2e.config")

	return InitCmd
}

func initConfig(_ *cobra.Command, _ []string) {
	err := config.WriteConfig(configFile, initConf)
	if err != nil {
		fmt.Printf("error writing config file: %s", err.Error())
	}
}
