package local

import (
	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

// initTestRunner initializes a runner form tests
// it creates a runner with an account and copy contracts from deployer runner
func initTestRunner(
	name string,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	account config.Account,
	logger *runner.Logger,
	opts ...runner.E2ERunnerOption,
) (*runner.E2ERunner, error) {
	// initialize runner for test
	testRunner, err := zetae2econfig.RunnerFromConfig(
		deployerRunner.Ctx,
		name,
		deployerRunner.CtxCancel,
		conf,
		account,
		logger,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	// copy timeouts from deployer runner
	testRunner.CctxTimeout = deployerRunner.CctxTimeout
	testRunner.ReceiptTimeout = deployerRunner.ReceiptTimeout
	testRunner.TestFilter = deployerRunner.TestFilter

	// copy contracts from deployer runner
	if err := testRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}

	return testRunner, nil
}
