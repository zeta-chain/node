package main

import (
	"fmt"
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"golang.org/x/sync/errgroup"
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

	// DeployerAddress is the address of the account for deploying networks
	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263" // #nosec G101 - used for testing

	// UserERC20Address is the address of the account for testing ERC20
	UserERC20Address    = ethcommon.HexToAddress("0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6")
	UserERC20PrivateKey = "fda3be1b1517bdf48615bdadacc1e6463d2865868dc8077d2cdcfa4709a16894" // #nosec G101 - used for testing

	// UserZetaTestAddress is the address of the account for testing Zeta
	UserZetaTestAddress    = ethcommon.HexToAddress("0x5cC2fBb200A929B372e3016F1925DcF988E081fd")
	UserZetaTestPrivateKey = "729a6cdc5c925242e7df92fdeeb94dadbf2d0b9950d4db8f034ab27a3b114ba7" // #nosec G101 - used for testing

	// UserBitcoinAddress is the address of the account for testing Bitcoin
	UserBitcoinAddress    = ethcommon.HexToAddress("0x283d810090EdF4043E75247eAeBcE848806237fD")
	UserBitcoinPrivateKey = "7bb523963ee2c78570fb6113d886a4184d42565e8847f1cb639f5f5e2ef5b37a" // #nosec G101 - used for testing

	// UserEtherAddress is the address of the account for testing Ether
	UserEtherAddress    = ethcommon.HexToAddress("0x8D47Db7390AC4D3D449Cc20D799ce4748F97619A")
	UserEtherPrivateKey = "098e74a1c2261fa3c1b8cfca8ef2b4ff96c73ce36710d208d1f6535aef42545d" // #nosec G101 - used for testing

	// UserMiscAddress is the address of the account for miscellaneous tests
	UserMiscAddress    = ethcommon.HexToAddress("0x90126d02E41c9eB2a10cfc43aAb3BD3460523Cdf")
	UserMiscPrivateKey = "853c0945b8035a501b1161df65a17a0a20fc848bda8975a8b4e9222cc6f84cd4" // #nosec G101 - used for testing

	// UserERC20AdvancedAddress is the address of the account for testing ERC20 advanced features
	UserERC20AdvancedAddress    = ethcommon.HexToAddress("0xcC8487562AAc220ea4406196Ee902C7c076966af")
	UserERC20AdvancedPrivateKey = "95409f1f0e974871cc26ba98ffd31f613aa1287d40c0aea6a87475fc3521d083" // #nosec G101 - used for testing

	FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass" // #nosec G101 - used for testing
)

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
	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

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
	deployerRunner.SendZetaOnEvm(UserZetaTestAddress, 1000)
	deployerRunner.SendZetaOnEvm(UserBitcoinAddress, 1000)
	deployerRunner.SendZetaOnEvm(UserEtherAddress, 1000)

	// error group for running multiple smoke tests concurrently
	var eg errgroup.Group

	// initialize runner for erc20 test
	erc20Runner, err := initERC20Runner(conf, deployerRunner, verbose)
	if err != nil {
		panic(err)
	}

	// initialize runner for zeta test
	zetaRunner, err := initZetaRunner(conf, deployerRunner, verbose)
	if err != nil {
		panic(err)
	}

	// initialize runner for bitcoin test
	bitcoinRunner, err := initBitcoinRunner(conf, deployerRunner, verbose)
	if err != nil {
		panic(err)
	}

	// initialize runner for ether test
	etherRunner, err := initEtherRunner(conf, deployerRunner, verbose)
	if err != nil {
		panic(err)
	}

	// run tests
	eg.Go(erc20TestRoutine(erc20Runner))
	eg.Go(zetaTestRoutine(zetaRunner))
	eg.Go(bitcoinTestRoutine(bitcoinRunner))
	eg.Go(ethereumTestRoutine(etherRunner))

	// deploy zevm swap and context apps
	//logger.Print("⚙️ setting up ZEVM swap and context apps")
	//sm.SetupZEVMSwapApp()
	//sm.SetupContextApp()

	if err := eg.Wait(); err != nil {
		logger.Print("❌ %v", err)
		os.Exit(1)
	}

	logger.Print("✅ smoke tests completed in %s", time.Since(testStartTime).String())
}

// erc20TestRoutine runs erc20 related smoke tests
func erc20TestRoutine(erc20Runner *runner.SmokeTestRunner) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the smoke tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("erc20 panic: %v", r)
			}
		}()

		erc20Runner.DepositZeta()
		erc20Runner.DepositEther()
		erc20Runner.DepositERC20()
		//erc20Runner.SetupBitcoinAccount()
		//erc20Runner.CheckZRC20ReserveAndSupply()

		// run erc20 test
		if err := erc20Runner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestDepositAndCallRefundName,
			smoketests.TestMultipleERC20DepositName,
			smoketests.TestWithdrawERC20Name,
			smoketests.TestMultipleWithdrawsName,
			smoketests.TestPauseZRC20Name,
			smoketests.TestERC20DepositAndCallRefundName,
			smoketests.TestUpdateBytecodeName,
			smoketests.TestWhitelistERC20Name,
		); err != nil {
			return err
		}

		return err
	}
}

// zetaTestRoutine runs Zeta transfer and message passing related smoke tests
func zetaTestRoutine(zetaRunner *runner.SmokeTestRunner) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the smoke tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("zeta panic: %v", r)
			}
		}()

		zetaRunner.DepositZeta()
		zetaRunner.DepositEther()
		//zetaRunner.SetupBitcoinAccount()
		//zetaRunner.CheckZRC20ReserveAndSupply()

		// run erc20 test
		if err := zetaRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestSendZetaOutName,
			smoketests.TestMessagePassingName,
			smoketests.TestMessagePassingRevertFailName,
			smoketests.TestMessagePassingRevertSuccessName,
		); err != nil {
			return err
		}

		return err
	}
}

// bitcoinTestRoutine runs Bitcoin related smoke tests
func bitcoinTestRoutine(bitcoinRunner *runner.SmokeTestRunner) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the smoke tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("bitcoin panic: %v", r)
			}
		}()

		bitcoinRunner.DepositZeta()
		bitcoinRunner.DepositEther()
		bitcoinRunner.SetupBitcoinAccount()
		bitcoinRunner.DepositBTC()
		//bitcoinRunner.CheckZRC20ReserveAndSupply()

		// run bitcoin test
		if err := bitcoinRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestBitcoinWithdrawName,
			smoketests.TestSendZetaOutBTCRevertName,
			smoketests.TestCrosschainSwapName,
		); err != nil {
			return err
		}

		return err
	}
}

// ethereumTestRoutine runs Ethereum related smoke tests
func ethereumTestRoutine(ethereumRunner *runner.SmokeTestRunner) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the smoke tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("ethereum panic: %v", r)
			}
		}()

		ethereumRunner.DepositZeta()
		ethereumRunner.DepositEther()
		//ethereumRunner.SetupBitcoinAccount()
		ethereumRunner.SetupContextApp()
		//ethereumRunner.CheckZRC20ReserveAndSupply()

		// run ethereum test
		if err := ethereumRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestContextUpgradeName,
			smoketests.TestEtherDepositAndCallName,
			smoketests.TestDepositEtherLiquidityCapName,
		); err != nil {
			return err
		}

		return err
	}
}

// initERC20Runner initializes a runner for erc20 tests
func initERC20Runner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for erc20 test
	erc20Runner, err := runnerFromConfig(
		conf,
		UserERC20Address,
		UserERC20PrivateKey,
		runner.NewLogger(verbose, color.FgGreen, "erc20"),
	)
	if err != nil {
		return nil, err
	}
	if err := erc20Runner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return erc20Runner, nil
}

// initZetaRunner initializes a runner for zeta tests
func initZetaRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for zeta test
	zetaRunner, err := runnerFromConfig(
		conf,
		UserZetaTestAddress,
		UserZetaTestPrivateKey,
		runner.NewLogger(verbose, color.FgBlue, "zeta"),
	)
	if err != nil {
		return nil, err
	}
	if err := zetaRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return zetaRunner, nil
}

// initBitcoinRunner initializes a runner for bitcoin tests
func initBitcoinRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for bitcoin test
	bitcoinRunner, err := runnerFromConfig(
		conf,
		UserBitcoinAddress,
		UserBitcoinPrivateKey,
		runner.NewLogger(verbose, color.FgYellow, "bitcoin"),
	)
	if err != nil {
		return nil, err
	}
	if err := bitcoinRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return bitcoinRunner, nil
}

// initEtherRunner initializes a runner for ether tests
func initEtherRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for ether test
	etherRunner, err := runnerFromConfig(
		conf,
		UserEtherAddress,
		UserEtherPrivateKey,
		runner.NewLogger(verbose, color.FgMagenta, "ether"),
	)
	if err != nil {
		return nil, err
	}
	if err := etherRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return etherRunner, nil
}

// initMiscRunner initializes a runner for miscellaneous tests
func initMiscRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for misc test
	miscRunner, err := runnerFromConfig(
		conf,
		UserMiscAddress,
		UserMiscPrivateKey,
		runner.NewLogger(verbose, color.FgCyan, "misc"),
	)
	if err != nil {
		return nil, err
	}
	if err := miscRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return miscRunner, nil
}

// initERC20AdvancedRunner initializes a runner for erc20 advanced tests
func initERC20AdvancedRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for erc20 advanced test
	erc20AdvancedRunner, err := runnerFromConfig(
		conf,
		UserERC20AdvancedAddress,
		UserERC20AdvancedPrivateKey,
		runner.NewLogger(verbose, color.FgHiGreen, "erc20advanced"),
	)
	if err != nil {
		return nil, err
	}
	if err := erc20AdvancedRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return erc20AdvancedRunner, nil
}
