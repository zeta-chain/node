package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
)

const (
	InitCmdId = "init"
)

var initConf = config.Config{}
var configFile = ""

func NewInitCmd() *cobra.Command {
	var InitCmd = &cobra.Command{
		Use:   InitCmdId,
		Short: "Run Local Stress Test",
		Run:   initConfig,
	}

	InitCmd.Flags().StringVar(&initConf.RPCs.EVM, "ethURL", "http://eth:8545", "--ethURL http://eth:8545")
	InitCmd.Flags().StringVar(&initConf.RPCs.ZetaCoreGRPC, "grpcURL", "zetacore0:9090", "--grpcURL zetacore0:9090")
	InitCmd.Flags().StringVar(&initConf.RPCs.ZetaCoreRPC, "rpcURL", "http://zetacore0:26657", "--rpcURL http://zetacore0:26657")
	InitCmd.Flags().StringVar(&initConf.RPCs.Zevm, "zevmURL", "http://zetacore0:8545", "--zevmURL http://zetacore0:8545")
	InitCmd.Flags().StringVar(&initConf.RPCs.Bitcoin, "btcURL", "bitcoin:18443", "--grpcURL bitcoin:18443")

	InitCmd.Flags().StringVar(&initConf.ZetaChainID, "chainID", "athens_101-1", "--chainID athens_101-1")
	InitCmd.Flags().StringVar(&configFile, "cfg", "smoketest.config", "--cfg ./smoketest.config")

	return InitCmd
}

func initConfig(_ *cobra.Command, _ []string) {
	err := config.WriteConfig(configFile, initConf)
	if err != nil {
		fmt.Printf("error writing config file: %s", err.Error())
	}
}
