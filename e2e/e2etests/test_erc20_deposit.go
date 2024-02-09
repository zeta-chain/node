package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestERC20Deposit(sm *runner.E2ERunner) {
	hash := sm.DepositERC20WithAmountAndMessage(big.NewInt(100000), []byte{})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, hash.Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	sm.Logger.CCTX(*cctx, "deposit")
}
