package main

import (
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

var RootCmd = &cobra.Command{
	Use:   "zetaclientd",
	Short: "ZetaClient CLI",
}
var rootArgs = rootArguments{}

type rootArguments struct {
	zetaCoreHome string
}

func setHomeDir() {
	rootArgs.zetaCoreHome, _ = RootCmd.Flags().GetString(tmcli.HomeFlag)
}
