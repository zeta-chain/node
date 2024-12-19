package local

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
)

// statefulPrecompilesTestRoutine runs steateful precompiles related e2e tests
func statefulPrecompilesTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserPrecompile

		precompileRunner, err := initTestRunner(
			"precompiles",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgRed, "precompiles"),
		)
		if err != nil {
			return err
		}

		// Initialize a ZetaTxServer with the precompile user account.
		// It's needed to send messages on behalf of the precompile user.
		zetaTxServer, err := txserver.NewZetaTxServer(
			conf.RPCs.ZetaCoreRPC,
			[]string{
				sdk.AccAddress(conf.AdditionalAccounts.UserPrecompile.EVMAddress().Bytes()).String(),
			},
			[]string{
				conf.AdditionalAccounts.UserPrecompile.RawPrivateKey.String(),
			},
			conf.ZetaChainID,
		)
		if err != nil {
			return err
		}

		precompileRunner.ZetaTxServer = zetaTxServer

		precompileRunner.Logger.Print("🏃 starting stateful precompiled contracts tests")
		startTime := time.Now()

		// Send ERC20 that will be depositted into ERC20ZRC20 tokens.
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 1e7)
		precompileRunner.WaitForTxReceiptOnEVM(txERC20Send)

		testsToRun, err := precompileRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("precompiled contracts tests failed: %v", err)
		}

		if err := precompileRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("precompiled contracts tests failed: %v", err)
		}

		precompileRunner.Logger.Print("🍾 precompiled contracts tests completed in %s", time.Since(startTime).String())

		return err
	}
}
