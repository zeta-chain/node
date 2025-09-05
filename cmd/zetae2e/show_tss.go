package main

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

// NewShowTSSCmd returns the show TSS command
// which shows the TSS address in the network
func NewShowTSSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tss [config]",
		Short: "Show address of the TSS",
		RunE:  runShowTSS,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runShowTSS(_ *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0], true)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(true, color.FgHiCyan, "")

	// initialize context
	ctx, cancel := context.WithCancelCause(context.Background())

	// initialize deployer runner with config
	testRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"tss",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
	)
	if err != nil {
		cancel(err)
		return err
	}

	// fetch the TSS address
	if err := testRunner.SetTSSAddresses(); err != nil {
		return err
	}

	// print the TSS address
	logger.Info("TSS EVM address: %s\n", testRunner.TSSAddress.Hex())
	logger.Info("TSS BTC address: %s\n", testRunner.BTCTSSAddress.EncodeAddress())

	return nil
}
