package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
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

	bitcoinBalance, err := getBitcoinBalance(r)
	if err != nil {
		cancel()
		return err
	}
	logger.Print("* BTC balance: %s", bitcoinBalance)

	return nil
}

func getBitcoinBalance(r *runner.SmokeTestRunner) (string, error) {
	addr, err := r.GetBtcAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get BTC address: %w", err)
	}

	address, err := btcutil.DecodeAddress(addr, r.BitcoinParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode BTC address: %w", err)
	}

	//_, err = r.BtcRPCClient.CreateWallet("user")
	//if err != nil {
	//	if !strings.Contains(err.Error(), "Database already exists") {
	//		return "", fmt.Errorf("failed to create wallet: %w", err)
	//	}
	//}
	//
	//skBytes, err := hex.DecodeString(r.DeployerPrivateKey)
	//if err != nil {
	//	return "", fmt.Errorf("failed to decode private key: %w", err)
	//}
	//
	//sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	//privkeyWIF, err := btcutil.NewWIF(sk, r.BitcoinParams, true)
	//if err != nil {
	//	return "", fmt.Errorf("failed to create WIF: %w", err)
	//}
	//
	//// Import the key to the wallet
	//err = r.BtcRPCClient.ImportPrivKeyLabel(privkeyWIF, "user")
	//if err != nil {
	//	return "", fmt.Errorf("failed to import private key: %w", err)
	//}

	//address, err := btcutil.NewAddressWitnessPubKeyHash(
	//	btcutil.Hash160(privkeyWIF.SerializePubKey()),
	//	r.BitcoinParams,
	//)
	//if err != nil {
	//	return "", err
	//}

	unspentList, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
	if err != nil {
		return "", fmt.Errorf("failed to list unspent: %w", err)
	}

	// calculate total amount
	var totalAmount btcutil.Amount
	for _, unspent := range unspentList {
		totalAmount += btcutil.Amount(unspent.Amount * 1e8)
	}

	return totalAmount.String(), nil
}
