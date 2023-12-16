package main

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smoketest",
		Short: "Smoke Test CLI",
	}
	cmd.AddCommand(NewLocalCmd())

	return cmd
}
