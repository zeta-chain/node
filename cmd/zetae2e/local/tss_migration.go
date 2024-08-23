package local

import (
	"os"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TSSMigration(deployerRunner *runner.E2ERunner, logger *runner.Logger, verbose bool, conf config.Config) {
	migrationStartTime := time.Now()
	logger.Print("🏁 starting tss migration")

	response, err := deployerRunner.CctxClient.LastZetaHeight(
		deployerRunner.Ctx,
		&crosschaintypes.QueryLastZetaHeightRequest{},
	)
	require.NoError(deployerRunner, err)
	err = deployerRunner.ZetaTxServer.UpdateKeygen(response.Height)
	require.NoError(deployerRunner, err)

	// Generate new TSS
	waitKeygenHeight(deployerRunner.Ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 0)

	// Run migration
	// migrationRoutine runs migration e2e test , which migrates funds from the older TSS to the new one
	// The zetaclient restarts required for this process are managed by the background workers in zetaclient (TSSListener)
	fn := migrationRoutine(conf, deployerRunner, verbose, e2etests.TestMigrateTSSName)

	if err := fn(); err != nil {
		logger.Print("❌ %v", err)
		logger.Print("❌ tss migration failed")
		os.Exit(1)
	}
	deployerRunner.UpdateTssAddressForConnector()
	deployerRunner.UpdateTssAddressForErc20custody()
	logger.Print("✅ migration completed in %s ", time.Since(migrationStartTime).String())
}
