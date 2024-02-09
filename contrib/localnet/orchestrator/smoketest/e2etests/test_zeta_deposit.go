package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestZetaDeposit(sm *runner.E2ERunner) {
	// Deposit 1 Zeta
	hash := sm.DepositZetaWithAmount(big.NewInt(1e18))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, hash.Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	sm.Logger.CCTX(*cctx, "deposit")
}
