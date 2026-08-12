package local

import (
	"os"
	"time"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// triggerKeygenReset runs the keygen-reset test in isolation. The test resets the keygen
// record the way zetacore does on an observer set change, then asserts the network still
// signs an outbound with the TSS key it already holds.
func triggerKeygenReset(deployerRunner *runner.E2ERunner, logger *runner.Logger, _ config.Config) {
	logger.Print("🏁 starting keygen reset test")
	startTime := time.Now()

	testsToRun, err := deployerRunner.GetE2ETestsToRunByName(
		e2etests.AllE2ETests,
		e2etests.TestKeygenResetSigningName,
	)
	if err != nil {
		logger.Print("❌ %v", err)
		os.Exit(1)
	}

	if err := deployerRunner.RunE2ETests(testsToRun); err != nil {
		logger.Print("❌ keygen reset test failed: %v", err)
		os.Exit(1)
	}

	logger.Print("✅ keygen reset test completed in %s", time.Since(startTime).String())
}
