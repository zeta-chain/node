package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
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
	maxDescriptionLength := 0
	for _, test := range tests {
		if len(test.Name) > maxNameLength {
			maxNameLength = len(test.Name)
		}
		if len(test.Description) > maxDescriptionLength {
			maxDescriptionLength = len(test.Description)
		}
	}

	// Formatting and printing the table
	formatString := fmt.Sprintf("%%-%ds | %%-%ds | %%s", maxNameLength, maxDescriptionLength)
	logger.Print(formatString, "Name", "Description", "Arguments (default)")
	for _, test := range tests {
		logger.Print(formatString, test.Name, test.Description, test.ArgsDescription())
	}
}
