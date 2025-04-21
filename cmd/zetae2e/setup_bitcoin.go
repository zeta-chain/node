package main

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

// NewSetupBitcoinCmd sets up bitcoin wallet for e2e tests
// should be run in case bitcoin e2e tests return load wallet errors
func NewSetupBitcoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup-bitcoin [config-file] ",
		Short: "Setup Bitcoin wallet for e2e tests",
		RunE:  runSetupBitcoin,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runSetupBitcoin(_ *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0], true)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(false, color.FgHiYellow, "")

	// initialize context
	ctx, cancel := context.WithCancelCause(context.Background())

	// initialize deployer runner with config
	r, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"e2e",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
	)
	if err != nil {
		cancel(err)
		return err
	}

	if err := r.SetTSSAddresses(); err != nil {
		return err
	}

	r.SetupBitcoinAccounts(true)

	logger.Print("* BTC setup done")

	return nil
}
