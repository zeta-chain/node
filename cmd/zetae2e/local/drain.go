package local

import (
	"os"
	"time"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// triggerDrain runs the emergency drain e2e test in isolation. The test self-funds the
// TSS, arms the drain payload, and asserts the funds are swept to the safe receivers.
func triggerDrain(deployerRunner *runner.E2ERunner, logger *runner.Logger, _ config.Config) {
	logger.Print("🏁 starting emergency drain test")
	startTime := time.Now()

	testsToRun, err := deployerRunner.GetE2ETestsToRunByName(e2etests.AllE2ETests, e2etests.TestDrainTSSName)
	if err != nil {
		logger.Print("❌ %v", err)
		os.Exit(1)
	}

	if err := deployerRunner.RunE2ETests(testsToRun); err != nil {
		logger.Print("❌ drain test failed: %v", err)
		os.Exit(1)
	}

	logger.Print("✅ drain test completed in %s", time.Since(startTime).String())
}
