package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const flagVerbose = "verbose"
const flagConfig = "config"
const flagERC20ChainName = "erc20-chain-name"
const flagERC20Symbol = "erc20-symbol"

// NewRunCmd returns the run command
// which runs the E2E from a config file describing the tests, networks, and accounts
func NewRunCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "run [testname1]:[arg1],[arg2] [testname2]:[arg1],[arg2]...",
		Short: "Run one or more E2E tests with optional arguments",
		Long: `Run one or more E2E tests specified by their names and optional arguments.
For example: zetae2e run deposit:1000 withdraw: --config config.yml`,
		RunE: runE2ETest,
		Args: cobra.MinimumNArgs(1), // Ensures at least one test is provided
	}

	cmd.Flags().StringVarP(&configFile, flagConfig, "c", "", "path to the configuration file")
	if err := cmd.MarkFlagRequired(flagConfig); err != nil {
		fmt.Println("Error marking flag as required")
		os.Exit(1)
	}

	cmd.Flags().String(flagERC20ChainName, "", "chain_name from /zeta-chain/observer/supportedChains")
	cmd.Flags().String(flagERC20Symbol, "", "symbol from /zeta-chain/fungible/foreign_coins")

	// Retain the verbose flag
	cmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")

	return cmd
}

func runE2ETest(cmd *cobra.Command, args []string) error {
	// read the config file
	configPath, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	conf, err := config.ReadConfig(configPath)
	if err != nil {
		return err
	}

	// read flag
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(verbose, color.FgHiCyan, "e2e")

	// update config with dynamic ERC20
	erc20ChainName, err := cmd.Flags().GetString(flagERC20ChainName)
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

	// set config
	app.SetConfig()

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize deployer runner with config
	testRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"e2e",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
	)
	if err != nil {
		return err
	}

	testStartTime := time.Now()
	logger.Print("starting tests")

	// fetch the TSS address
	if err := testRunner.SetTSSAddresses(); err != nil {
		return err
	}

	// set timeout
	testRunner.CctxTimeout = 60 * time.Minute
	testRunner.ReceiptTimeout = 60 * time.Minute

	// parse test names and arguments from cmd args and run them
	userTestsConfigs, err := parseCmdArgsToE2ETestRunConfig(args)
	if err != nil {
		return err
	}

	testsToRun, err := testRunner.GetE2ETestsToRunByConfig(e2etests.AllE2ETests, userTestsConfigs)
	if err != nil {
		return err
	}
	reports, err := testRunner.RunE2ETestsIntoReport(testsToRun)
	if err != nil {
		return err
	}

	// Print tests completion info
	logger.Print("tests finished successfully in %s", time.Since(testStartTime).String())
	testRunner.Logger.SetColor(color.FgHiRed)
	testRunner.Logger.SetColor(color.FgHiGreen)
	testRunner.PrintTestReports(reports)

	return nil
}

// parseCmdArgsToE2ETestRunConfig parses command-line arguments into a slice of E2ETestRunConfig structs.
func parseCmdArgsToE2ETestRunConfig(args []string) ([]runner.E2ETestRunConfig, error) {
	tests := []runner.E2ETestRunConfig{}
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, errors.New("command arguments should be in format: testName:testArgs")
		}
		if parts[0] == "" {
			return nil, errors.New("missing testName")
		}
		testName := parts[0]
		testArgs := []string{}
		if parts[1] != "" {
			testArgs = strings.Split(parts[1], ",")
		}
		tests = append(tests, runner.E2ETestRunConfig{
			Name: testName,
			Args: testArgs,
		})
	}
	return tests, nil
}

// findERC20 loads ERC20 addresses via gRPC given CLI flags
func findERC20(ctx context.Context, conf config.Config, erc20ChainName, erc20Symbol string) (string, string, error) {
	clients, err := zetae2econfig.GetZetacoreClient(conf)
	if err != nil {
		return "", "", fmt.Errorf("get zeta clients: %w", err)
	}

	supportedChainsRes, err := clients.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return "", "", fmt.Errorf("get chain params: %w", err)
	}

	chainID := int64(0)
	for _, chain := range supportedChainsRes.Chains {
		if chain.Name == erc20ChainName {
			chainID = chain.ChainId
			break
		}
	}
	if chainID == 0 {
		return "", "", fmt.Errorf("chain %s not found", erc20ChainName)
	}

	foreignCoinsRes, err := clients.Fungible.ForeignCoinsAll(ctx, &fungibletypes.QueryAllForeignCoinsRequest{})
	if err != nil {
		return "", "", fmt.Errorf("get foreign coins: %w", err)
	}

	for _, coin := range foreignCoinsRes.ForeignCoins {
		if coin.ForeignChainId != chainID {
			continue
		}
		// sometimes symbol is USDT, sometimes it's like USDT.SEPOLIA
		if strings.HasPrefix(coin.Symbol, erc20Symbol) || strings.HasSuffix(coin.Symbol, erc20Symbol) {
			return coin.Asset, coin.Zrc20ContractAddress, nil
		}
	}
	return "", "", fmt.Errorf("erc20 %s not found on %s", erc20Symbol, erc20ChainName)
}
