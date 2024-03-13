package main

import (
	"context"
	"errors"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
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
	conf, err := config.ReadConfig(args[0])
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(false, color.FgHiYellow, "")

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

	if err := r.SetTSSAddresses(); err != nil {
		return err
	}

	r.SetupBitcoinAccount(true)

	logger.Print("* BTC setup done")

	return nil
}
