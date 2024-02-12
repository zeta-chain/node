package e2etests

import (
	"fmt"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestBitcoinDeposit(r *runner.E2ERunner) {

	r.SetBtcAddress(r.Name, false)

	txHash := r.DepositBTCWithAmount(0.001)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected mined status; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}
}
