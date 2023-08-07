package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/zeta-chain/zetacore/app"
)

func main() {
	cmdcfg.RegisterDenoms()

	rootCmd, _ := NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
