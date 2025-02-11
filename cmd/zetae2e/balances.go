package main

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

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
	cmd.Flags().String(flagERC20Network, "", "network from /zeta-chain/observer/supportedChains")
	cmd.Flags().String(flagERC20Symbol, "", "symbol from /zeta-chain/fungible/foreign_coins")
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

	// initialize context
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	// update config with dynamic ERC20
	erc20ChainName, err := cmd.Flags().GetString(flagERC20Network)
	if err != nil {
		return err
	}
	erc20Symbol, err := cmd.Flags().GetString(flagERC20Symbol)
	if err != nil {
		return err
	}
	if erc20ChainName != "" && erc20Symbol != "" {
		erc20Asset, zrc20ContractAddress, err := findERC20(
			cmd.Context(),
			conf,
			erc20ChainName,
			erc20Symbol,
		)
		if err != nil {
			return err
		}
		conf.Contracts.EVM.ERC20 = config.DoubleQuotedString(erc20Asset)
		conf.Contracts.ZEVM.ERC20ZRC20Addr = config.DoubleQuotedString(zrc20ContractAddress)
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

	balances, err := r.GetAccountBalances(skipBTC)
	if err != nil {
		cancel(err)
		return err
	}
	r.PrintAccountBalances(balances)

	return nil
}
