package e2etests

import (
	"fmt"
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestZetaWithdraw(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestZetaWithdraw requires exactly one argument for the withdrawal.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("invalid amount specified")
	}

	r.DepositAndApproveWZeta(amount)
	tx := r.WithdrawZeta(amount, true)

	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}
}
