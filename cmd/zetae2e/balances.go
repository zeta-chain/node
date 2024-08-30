package main

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

const flagSkipBTC = "skip-btc"

// NewBalancesCmd returns the balances command
// which shows from the key and rpc, the balance of the account on different network
func NewBalancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances [config-file]",
		Short: "Show account balances on networks for E2E tests",
		RunE:  runBalances,
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().Bool(
		flagSkipBTC,
		false,
		"skip the BTC network",
	)
	return cmd
}

func runBalances(cmd *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0])
	if err != nil {
		return err
	}

	skipBTC, err := cmd.Flags().GetBool(flagSkipBTC)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(false, color.FgHiCyan, "")

	// set config
	app.SetConfig()

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

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
		cancel()
		return err
	}

	balances, err := r.GetAccountBalances(skipBTC)
	if err != nil {
		cancel()
		return err
	}
	r.PrintAccountBalances(balances)

	return nil
}
