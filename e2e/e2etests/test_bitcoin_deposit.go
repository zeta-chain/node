package e2etests

import (
	"fmt"
	"strconv"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestBitcoinDeposit(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestBitcoinDeposit requires exactly one argument for the amount.")
	}

	depositAmount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic("Invalid deposit amount specified for TestBitcoinDeposit.")
	}

	r.SetBtcAddress(r.Name, false)

	txHash := r.DepositBTCWithAmount(depositAmount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected mined status; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}
}
