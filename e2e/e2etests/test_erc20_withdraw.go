package e2etests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestERC20Withdraw(r *runner.E2ERunner, args []string) {
	approvedAmount := big.NewInt(1e18)
	if len(args) != 1 {
		panic("TestERC20Withdraw requires exactly one argument for the withdrawal amount.")
	}

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		panic("Invalid withdrawal amount specified for TestERC20Withdraw.")
	}

	if withdrawalAmount.Cmp(approvedAmount) >= 0 {
		panic("Withdrawal amount must be less than the approved amount (1e18).")
	}

	// approve
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ERC20ZRC20Addr, approvedAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// withdraw
	tx = r.WithdrawERC20(withdrawalAmount)

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	verifyTransferAmountFromCCTX(r, cctx, withdrawalAmount.Int64())
}
