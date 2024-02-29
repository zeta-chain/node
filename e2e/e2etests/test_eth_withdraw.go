package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestEtherWithdraw tests the withdraw of ether
func TestEtherWithdraw(r *runner.E2ERunner) {
	// approve
	tx, err := r.ETHZRC20.Approve(r.ZevmAuth, r.ETHZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	tx, err = r.ETHZRC20.Withdraw(r.ZevmAuth, r.DeployerAddress.Bytes(), big.NewInt(100000))
	if err != nil {
		panic(err)
	}
	r.Logger.EVMTransaction(*tx, "withdraw")

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}
	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.ZRC20Withdrawal(r.ETHZRC20, *receipt, "withdraw")

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic("cctx status is not outbound mined")
	}
}
