package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
)

func main() {
	cmdcfg.RegisterDenoms()

	rootCmd, _ := NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {

		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			processError(e)
			os.Exit(1)
		}
	}
}

func processError(err error) {
	// --ledger flag can't be used with Ethereum HD path
	if strings.Contains(err.Error(), "cannot set custom bip32 path with ledger") {
		printNotice([]string{
			"note: --ledger flag can't be used with Ethereum HD path (used by default)",
			"Please set a blank path with --hd-path=\"\" to use Cosmos HD path instead.",
		})
		os.Exit(1)
	}
}

func printNotice(messages []string) {
	if len(messages) == 0 {
		return
	}
	border := strings.Repeat("*", len(messages[0])+4) // 4 to account for padding
	fmt.Println(border)
	for _, message := range messages {
		fmt.Printf("* %s \n", message)
	}
	fmt.Println(border)
}
