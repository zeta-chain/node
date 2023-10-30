package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func docsCmd(cmd *cobra.Command, args []string) error {
	var path string

	// If path is provided as an argument, use it. Else, get from flag.
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0750)
		if err != nil {
			return err
		}
	}

	err := doc.GenMarkdownTree(cmd.Root(), path)
	if err != nil {
		return err
	}
	return nil
}

func docsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs [path]",
		Short: "Generate markdown documentation for zetacored",
		RunE:  docsCmd,
		Args:  cobra.MaximumNArgs(1),
	}

	cmd.Flags().String("path", "docs/cli/zetacored", "Path where the docs will be generated")

	return cmd
}
