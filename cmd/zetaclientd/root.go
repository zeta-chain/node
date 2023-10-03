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

func setHomeDir() error {
	var err error
	rootArgs.zetaCoreHome, err = RootCmd.Flags().GetString(tmcli.HomeFlag)
	return err
}
