package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

const (
	flagNetwork = "network"
	flagSkipBTC = "skip-btc"
)

// NewBalancesCmd returns the balances command
// which shows from the key and rpc, the balance of the account on different network
func NewBalancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances [config-file]",
		Short: "Show account balances on networks for E2E tests",
		RunE:  runBalances,
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().String(
		flagNetwork,
		"",
		"external chain network to query native balances for (e.g. polygon, bsc, eth, base, arbitrum, avalanche, btc, solana, sui, ton)",
	)
	// --skip-btc is kept for backward compat but intentionally not read.
	// The zt/e2e workflow is updated in lockstep to use --network instead.
	cmd.Flags().Bool(
		flagSkipBTC,
		false,
		"[deprecated] use --network instead",
	)
	_ = cmd.Flags().MarkDeprecated(flagSkipBTC, "use --network instead")
	_ = cmd.Flags().MarkHidden(flagSkipBTC)

	registerERC20Flags(cmd)
	return cmd
}

func runBalances(cmd *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0], true)
	if err != nil {
		return err
	}

	network, err := cmd.Flags().GetString(flagNetwork)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(false, color.FgHiCyan, "")

	// initialize context
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	err = processZRC20Flags(cmd, &conf)
	if err != nil {
		return fmt.Errorf("process ZRC20 flags: %w", err)
	}

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

	balances, err := r.GetAccountBalances(network)
	if err != nil {
		cancel(err)
		return err
	}
	r.PrintAccountBalances(balances, network)

	return nil
}
