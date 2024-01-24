package main

import (
	"context"
	"errors"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
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
	return cmd
}

func runBalances(_ *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0])
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(false, color.FgHiCyan, "")

	// set config
	app.SetConfig()

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// get EVM address from config
	evmAddr := conf.Accounts.EVMAddress
	if !ethcommon.IsHexAddress(evmAddr) {
		cancel()
		return errors.New("invalid EVM address")
	}

	// initialize deployer runner with config
	r, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"e2e",
		cancel,
		conf,
		ethcommon.HexToAddress(evmAddr),
		conf.Accounts.EVMPrivKey,
		utils.FungibleAdminName, // placeholder value, not used
		FungibleAdminMnemonic,   // placeholder value, not used
		logger,
	)
	if err != nil {
		cancel()
		return err
	}

	balances, err := r.GetAccountBalances()
	if err != nil {
		cancel()
		return err
	}
	r.PrintAccountBalances(balances)

	return nil
}
