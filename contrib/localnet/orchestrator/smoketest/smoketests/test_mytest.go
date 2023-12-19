package smoketests

import (
	"fmt"
	"time"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestMyTest(_ *runner.SmokeTestRunner) {
	utils.LoudPrintf("Custom test\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	// add your test here
}
