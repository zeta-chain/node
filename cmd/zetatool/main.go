package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/cmd/zetatool/filterdeposit"
)

var rootCmd = &cobra.Command{
	Use:   "zetatool",
	Short: "utility tool for zeta-chain",
}

func init() {
	rootCmd.AddCommand(filterdeposit.NewFilterDepositCmd())
	rootCmd.PersistentFlags().String(config.FlagConfig, "", "custom config file: --config filename.json")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
