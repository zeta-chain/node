package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/e2etests"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

// NewListTestsCmd returns the list test cmd
// which list the available tests
func NewListTestsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tests",
		Short: "List available tests",
		RunE:  runListTests,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func runListTests(_ *cobra.Command, _ []string) error {
	logger := runner.NewLogger(false, color.FgHiGreen, "")

	logger.Print("Available tests:")
	renderTests(logger, e2etests.AllE2ETests)

	return nil
}

func renderTests(logger *runner.Logger, tests []runner.E2ETest) {
	// Find the maximum length of the Name field
	maxNameLength := 0
	for _, test := range tests {
		if len(test.Name) > maxNameLength {
			maxNameLength = len(test.Name)
		}
	}

	// Formatting and printing the table
	formatString := fmt.Sprintf("%%-%ds | %%s", maxNameLength)
	logger.Print(formatString, "Name", "Description")
	for _, test := range tests {
		logger.Print(formatString, test.Name, test.Description)
	}
}
