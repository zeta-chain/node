package main

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

const flagPrivKey = "privkey"

// NewBitcoinAddressCmd returns the bitcoin address command
// which shows from the used config file, the bitcoin address that can be used to receive funds for the E2E tests
func NewBitcoinAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bitcoin-address [config-file] ",
		Short: "Show Bitcoin address to receive funds for E2E tests",
		RunE:  runBitcoinAddress,
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().Bool(
		flagPrivKey,
		false,
		"show the priv key in WIF format",
	)
	return cmd
}

func runBitcoinAddress(cmd *cobra.Command, args []string) error {
	showPrivKey, err := cmd.Flags().GetBool(flagPrivKey)
	if err != nil {
		return err
	}

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

	addr, privKey := r.GetBtcKeypair()

	logger.Print("* BTC address: %s", addr.EncodeAddress())
	if showPrivKey {
		logger.Print("* BTC privkey: %s", privKey.String())
	}

	return nil
}
