package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestZetaDepositNewAddress(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestZetaDepositNewAddress requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestZetaDepositNewAddress.")
	}

	newAddress := sample.EthAddress()
	hash := r.DepositZetaWithAmount(newAddress, amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}
