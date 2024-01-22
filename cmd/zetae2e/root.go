package main

import (
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd/zetae2e/local"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zetae2e",
		Short: "E2E tests CLI",
	}
	cmd.AddCommand(
		NewRunCmd(),
		local.NewLocalCmd(),
		NewStressTestCmd(),
		NewInitCmd(),
	)

	return cmd
}
