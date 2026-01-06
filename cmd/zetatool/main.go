package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/cli"
	"github.com/zeta-chain/node/cmd/zetatool/config"
)

var rootCmd = &cobra.Command{
	Use:   "zetatool",
	Short: "utility tool for zeta-chain",
}

func init() {
	rootCmd.AddCommand(cli.NewGetInboundBallotCMD())
	rootCmd.AddCommand(cli.NewTrackCCTXCMD())
	rootCmd.AddCommand(cli.NewApplicationDBStatsCMD())
	rootCmd.AddCommand(cli.NewTSSBalancesCMD())
	rootCmd.PersistentFlags().String(config.FlagConfig, "", "custom config file: --config filename.json")
	rootCmd.PersistentFlags().
		Bool(config.FlagDebug, false, "enable debug mode, to show more details on why the command might be failing")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
