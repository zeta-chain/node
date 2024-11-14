package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/app"
)

func UpgradeHandlerVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade-handler-version",
		Short: "Print the default upgrade handler version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(app.GetDefaultUpgradeHandlerVersion())
		},
	}
}
