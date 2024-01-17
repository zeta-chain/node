package main

import (
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/cmd/smoketest/local"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smoketest",
		Short: "Smoke Test CLI",
	}
	cmd.AddCommand(local.NewLocalCmd())

	return cmd
}
