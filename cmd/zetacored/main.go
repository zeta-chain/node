package main

import (
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/cmd/zetacored/cmd"
)

func main() {

	cmdcfg.RegisterDenoms()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
