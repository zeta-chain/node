package main

import (
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
)

var RootCmd = &cobra.Command{
	Use:   "zetaclientd",
	Short: "ZetaClient CLI",
}

var rootArgs = rootArguments{}

type rootArguments struct {
	zetaCoreHome string
}

func init() {
	rootArgs.zetaCoreHome = app.DefaultNodeHome
}
