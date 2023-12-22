package main

import (
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

const (
	flagContractsDeployed = "deployed"
	flagWaitForHeight     = "wait-for"
	flagConfigFile        = "config"
	flagVerbose           = "verbose"
)

var (
	// TODO: make these variables configurable
	// https://github.com/zeta-chain/node-private/issues/41

	SmokeTestTimeout = 30 * time.Minute

	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263" // #nosec G101 - used for testing

	UserERC20Address    = ethcommon.HexToAddress("0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6")
	UserERC20PrivateKey = "fda3be1b1517bdf48615bdadacc1e6463d2865868dc8077d2cdcfa4709a16894" // #nosec G101 - used for testing

	FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass" // #nosec G101 - used for testing
)

// 0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6	fda3be1b1517bdf48615bdadacc1e6463d2865868dc8077d2cdcfa4709a16894
// 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC   d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263

// 0x5cC2fBb200A929B372e3016F1925DcF988E081fd   729a6cdc5c925242e7df92fdeeb94dadbf2d0b9950d4db8f034ab27a3b114ba7

func NewLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run Local Smoketest",
		Run:   localSmokeTest,
	}
	cmd.Flags().Bool(
		flagContractsDeployed,
		false,
		"set to to true if running smoketest again with existing state",
	)
	cmd.Flags().Int64(
		flagWaitForHeight,
		0,
		"block height for smoketest to begin, ex. --wait-for 100",
	)
	cmd.Flags().String(
		flagConfigFile,
		"",
		"config file to use for the smoketest",
	)
	cmd.Flags().Bool(
		flagVerbose,
		false,
		"set to true to enable verbose logging",
	)
	return cmd
}

func localSmokeTest(cmd *cobra.Command, _ []string) {
	// fetch flags
	waitForHeight, err := cmd.Flags().GetInt64(flagWaitForHeight)
	if err != nil {
		panic(err)
	}
	contractsDeployed, err := cmd.Flags().GetBool(flagContractsDeployed)
	if err != nil {
		panic(err)
	}
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		panic(err)
	}
	logger := runner.NewLogger(verbose)

	testStartTime := time.Now()
	logger.Print("starting smoke tests")

	// start timer
	go func() {
		time.Sleep(SmokeTestTimeout)
		logger.Error("Smoke test timed out after %s", SmokeTestTimeout.String())
		os.Exit(1)
	}()

	// initialize smoke tests config
	conf, err := getConfig(cmd)
	if err != nil {
		panic(err)
	}

	// wait for a specific height on ZetaChain
	if waitForHeight != 0 {
		utils.WaitForBlockHeight(waitForHeight, conf.RPCs.ZetaCoreRPC, logger)
	}

	// set account prefix to zeta
	setCosmosConfig()

	// wait for Genesis
	logger.Print("⏳ wait 40s for genesis")
	time.Sleep(40 * time.Second)

	// initialize deployer runner with config
	deployerRunner, err := runnerFromConfig(conf, DeployerAddress, DeployerPrivateKey, logger)
	if err != nil {
		panic(err)
	}

	// wait for keygen to be completed
	waitKeygenHeight(deployerRunner.CctxClient, logger)

	// setting up the networks
	logger.Print("⚙️ setting up networks")
	startTime := time.Now()
	deployerRunner.SetTSSAddresses()
	deployerRunner.SetupEVM(contractsDeployed)
	deployerRunner.SetZEVMContracts()
	logger.Print("✅ setup completed in %s", time.Since(startTime))

	// fund accounts
	deployerRunner.SendZetaOnEvm(UserERC20Address, 1000)
	deployerRunner.SendUSDTOnEvm(UserERC20Address, 10)

	// initialize runner for erc20 test
	erc20Runner, err := runnerFromConfig(conf, UserERC20Address, UserERC20PrivateKey, logger)
	if err != nil {
		panic(err)
	}
	if err := erc20Runner.CopyAddressesFrom(deployerRunner); err != nil {
		panic(err)
	}

	// run erc20 test
	erc20Runner.WG.Add(1)
	go func() {
		erc20Runner.DepositZeta()
		erc20Runner.DepositEther()
		erc20Runner.SetupBitcoinAccount()
		//erc20Runner.DepositBTC()
		erc20Runner.DepositERC20()
		erc20Runner.CheckZRC20ReserveAndSupply()

		// run erc20 test
		if err := erc20Runner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestDepositAndCallRefundName,
			smoketests.TestMultipleERC20DepositName,
			smoketests.TestWithdrawERC20Name,
			smoketests.TestMultipleWithdrawsName,
			smoketests.TestPauseZRC20Name,
			smoketests.TestERC20DepositAndCallRefundName,
			smoketests.TestWhitelistERC20Name,
		); err != nil {
			panic(err)
		}
		erc20Runner.WG.Done()
	}()

	// deploy zevm swap and context apps
	//logger.Print("⚙️ setting up ZEVM swap and context apps")
	//sm.SetupZEVMSwapApp()
	//sm.SetupContextApp()

	deployerRunner.WG.Wait()
	erc20Runner.WG.Wait()
	logger.Print("✅ smoke tests completed in %s", time.Since(testStartTime).String())
}

func getConfig(cmd *cobra.Command) (config.Config, error) {
	configFile, err := cmd.Flags().GetString(flagConfigFile)
	if err != nil {
		return config.Config{}, err
	}

	// use default config if no config file is specified
	if configFile == "" {
		return config.DefaultConfig(), nil
	}

	return config.ReadConfig(configFile)
}
