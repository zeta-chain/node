package main

import (
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zetaclientd",
	Short: "ZetaClient CLI",
}
var rootArgs = rootArguments{}

type rootArguments struct {
	zetaCoreHome string
}

func setHomeDir() error {
	var err error
	rootArgs.zetaCoreHome, err = RootCmd.Flags().GetString(tmcli.HomeFlag)
	return err
}
