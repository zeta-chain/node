package smoketests

import (
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

func TestBitcoinDeposit(sm *runner.SmokeTestRunner) {

	sm.SetBtcAddress(sm.Name, false)

	sm.DepositBTCWithAmount(0.001)
}
