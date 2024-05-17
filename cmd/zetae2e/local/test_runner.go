package local

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

// initTestRunner initializes a runner form tests
// it creates a runner with an account and copy contracts from deployer runner
func initTestRunner(
	name string,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	userAddress ethcommon.Address,
	userPrivKey string,
	logger *runner.Logger,
) (*runner.E2ERunner, error) {
	// initialize runner for test
	testRunner, err := zetae2econfig.RunnerFromConfig(
		deployerRunner.Ctx,
		name,
		deployerRunner.CtxCancel,
		conf,
		userAddress,
		userPrivKey,
		utils.FungibleAdminName,
		FungibleAdminMnemonic,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// copy contracts from deployer runner
	if err := testRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}

	return testRunner, nil
}
