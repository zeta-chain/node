package main

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

const flagVerbose = "verbose"

// NewRunCmd returns the run command
// which runs the smoketest from a config file describing the tests, networks, and accounts
func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [config-file]",
		Short: "Run E2E tests from a config file",
		RunE:  runE2ETest,
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().Bool(
		flagVerbose,
		false,
		"set to true to enable verbose logging",
	)

	return cmd
}

func runE2ETest(cmd *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0])
	if err != nil {
		return err
	}

	// read flag
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(verbose, color.FgWhite, "e2e")

	// set config
	app.SetConfig()

	_ = conf

	testStartTime := time.Now()
	logger.Print("starting tests")

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	_ = ctx
	_ = cancel
	_ = testStartTime

	cancel()

	return nil
}
