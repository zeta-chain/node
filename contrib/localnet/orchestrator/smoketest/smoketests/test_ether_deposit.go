package smoketests

import "github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"

func TestEtherDeposit(sm *runner.SmokeTestRunner) {
	sm.DepositEther(false)
}
