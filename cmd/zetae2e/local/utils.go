package local

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// getConfig returns config from file from the command line flag
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

// setCosmosConfig set account prefix to zeta
func setCosmosConfig() {
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()
}

// initTestRunner initializes a runner form tests
// it creates a runner with an account and copy contracts from deployer runner
func initTestRunner(
	name string,
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	userAddress ethcommon.Address,
	userPrivKey string,
	logger *runner.Logger,
) (*runner.SmokeTestRunner, error) {
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

// waitKeygenHeight waits for keygen height
func waitKeygenHeight(
	ctx context.Context,
	cctxClient crosschaintypes.QueryClient,
	logger *runner.Logger,
) {
	// wait for keygen to be completed. ~ height 30
	keygenHeight := int64(60)
	logger.Print("⏳ wait height %v for keygen to be completed", keygenHeight)
	for {
		time.Sleep(2 * time.Second)
		response, err := cctxClient.LastZetaHeight(ctx, &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			logger.Error("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= keygenHeight {
			break
		}
		logger.Info("Last ZetaHeight: %d", response.Height)
	}
}
