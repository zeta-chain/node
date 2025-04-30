package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
)

const flagVerbose = "verbose"
const flagConfig = "config"
const flagFailFast = "fail-fast"
const flagOpTimeout = "op-timeout"

// NewRunCmd returns the run command
// which runs the E2E from a config file describing the tests, networks, and accounts
func NewRunCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "run [testname1]:[arg1],[arg2] [testname2]:[arg1],[arg2]...",
		Short: "Run one or more E2E tests with optional arguments",
		Long: `Run one or more E2E tests specified by their names and optional arguments.
For example: zetae2e run deposit:1000 withdraw: --config config.yml`,
		RunE:         runE2ETest,
		Args:         cobra.MinimumNArgs(1), // Ensures at least one test is provided
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&configFile, flagConfig, "c", "", "path to the configuration file")
	if err := cmd.MarkFlagRequired(flagConfig); err != nil {
		fmt.Println("Error marking flag as required")
		os.Exit(1)
	}

	registerERC20Flags(cmd)

	// Retain the verbose flag
	cmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")

	cmd.Flags().Bool(flagFailFast, false, "should a failure in one test cause an immediate halt")
	cmd.Flags().Duration(flagOpTimeout, time.Minute*60, "timeout for single operation (CCTX or get receipt)")

	return cmd
}

func runE2ETest(cmd *cobra.Command, args []string) error {
	// read the config file
	configPath, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	conf, err := config.ReadConfig(configPath, true)
	if err != nil {
		return err
	}

	// read flag
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		return err
	}

	failFast, err := cmd.Flags().GetBool(flagFailFast)
	if err != nil {
		return err
	}

	timeout, err := cmd.Flags().GetDuration(flagOpTimeout)
	if err != nil {
		return err
	}

	// initialize logger
	logger := runner.NewLogger(verbose, color.FgHiCyan, "e2e")

	err = processZRC20Flags(cmd, &conf)
	if err != nil {
		return fmt.Errorf("process ZRC20 flags: %w", err)
	}

	// initialize context
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	// if failFast option is not specified, overwrite context cancellation function
	// so that it is a no-op
	if !failFast {
		cancel = func(_ error) {}
	}

	var runnerOpts []runner.E2ERunnerOption

	// if keys are defined for all policy accounts, we initialize a ZETA tx server allowing to send admin actions
	emergencyKey := conf.PolicyAccounts.EmergencyPolicyAccount.RawPrivateKey.String()
	operationalKey := conf.PolicyAccounts.OperationalPolicyAccount.RawPrivateKey.String()
	adminKey := conf.PolicyAccounts.AdminPolicyAccount.RawPrivateKey.String()
	if emergencyKey != "" && operationalKey != "" && adminKey != "" {
		zetaTxServer, err := txserver.NewZetaTxServer(
			conf.RPCs.ZetaCoreRPC,
			[]string{utils.EmergencyPolicyName, utils.OperationalPolicyName, utils.AdminPolicyName},
			[]string{
				emergencyKey,
				operationalKey,
				adminKey,
			},
			conf.ZetaChainID,
		)
		if err != nil {
			return err
		}
		runnerOpts = append(runnerOpts, runner.WithZetaTxServer(zetaTxServer))
	}

	// initialize deployer runner with config
	testRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"e2e",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
		runnerOpts...,
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
	testRunner.CctxTimeout = timeout
	testRunner.ReceiptTimeout = timeout

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
	logger.Print("tests finished in %s", time.Since(testStartTime).String())
	testRunner.Logger.SetColor(color.FgHiRed)
	testRunner.Logger.SetColor(color.FgHiGreen)
	testRunner.PrintTestReports(reports)

	anyTestFailed := lo.ContainsBy(reports, func(r runner.TestReport) bool { return !r.Success })
	if anyTestFailed {
		return errors.New("tests failed")
	}

	return nil
}

// parseCmdArgsToE2ETestRunConfig parses command-line arguments into a slice of E2ETestRunConfig structs.
func parseCmdArgsToE2ETestRunConfig(args []string) ([]runner.E2ETestRunConfig, error) {
	tests := make([]runner.E2ETestRunConfig, 0, len(args))

	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		testName := parts[0]
		if testName == "" {
			return nil, errors.New("missing testName")
		}

		var testArgs []string
		if len(parts) > 1 {
			if parts[1] != "" {
				testArgs = strings.Split(parts[1], ",")
			}
		}

		tests = append(tests, runner.E2ETestRunConfig{
			Name: testName,
			Args: testArgs,
		})
	}
	return tests, nil
}
