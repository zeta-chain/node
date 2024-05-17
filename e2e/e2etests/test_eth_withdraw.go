package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestEtherWithdraw tests the withdraw of ether
func TestEtherWithdraw(r *runner.E2ERunner, args []string) {
	r.Logger.Info("TestEtherWithdraw")

	approvedAmount := big.NewInt(1e18)
	if len(args) != 1 {
		panic("TestEtherWithdraw requires exactly one argument for the withdrawal amount.")
	}

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		panic("Invalid withdrawal amount specified for TestEtherWithdraw.")
	}

	if withdrawalAmount.Cmp(approvedAmount) >= 0 {
		panic("Withdrawal amount must be less than the approved amount (1e18).")
	}

	// approve
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, approvedAmount)
	if err != nil {
		panic(err)
	}
	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	tx = r.WithdrawEther(withdrawalAmount)

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic("cctx status is not outbound mined")
	}

	r.Logger.Info("TestEtherWithdraw completed")
}
