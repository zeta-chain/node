package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

// TestEtherDeposit tests deposit of ethers
func TestEtherDeposit(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestEtherDeposit requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestEtherDeposit.")
	}

	hash := r.DepositEtherWithAmount(false, amount) // in wei
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}
