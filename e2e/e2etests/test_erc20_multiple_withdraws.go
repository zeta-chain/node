package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func TestMultipleERC20Withdraws(r *runner.E2ERunner, args []string) {
	approvedAmount := big.NewInt(1e18)
	if len(args) != 2 {
		panic("TestMultipleWithdraws requires exactly two arguments: the withdrawal amount and the number of withdrawals.")
	}

	withdrawalAmount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok || withdrawalAmount.Cmp(approvedAmount) >= 0 {
		panic("Invalid withdrawal amount specified for TestMultipleWithdraws.")
	}

	numberOfWithdrawals, ok := big.NewInt(0).SetString(args[1], 10)
	if !ok || numberOfWithdrawals.Int64() < 1 {
		panic("Invalid number of withdrawals specified for TestMultipleWithdraws.")
	}

	// calculate total withdrawal to ensure it doesn't exceed approved amount.
	totalWithdrawal := big.NewInt(0).Mul(withdrawalAmount, numberOfWithdrawals)
	if totalWithdrawal.Cmp(approvedAmount) >= 0 {
		panic("Total withdrawal amount exceeds approved limit.")
	}

	// deploy withdrawer
	withdrawerAddr, _, withdrawer, err := testcontract.DeployWithdrawer(r.ZEVMAuth, r.ZEVMClient)
	if err != nil {
		panic(err)
	}

	// approve
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, withdrawerAddr, approvedAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.Info("ERC20 ZRC20 approve receipt: status %d", receipt.Status)

	// approve gas token
	tx, err = r.ETHZRC20.Approve(r.ZEVMAuth, withdrawerAddr, approvedAmount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve gas token failed")
	}
	r.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// check the balance
	bal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("balance of deployer on ERC20 ZRC20: %d", bal)

	if bal.Int64() < totalWithdrawal.Int64() {
		panic("not enough ERC20 ZRC20 balance!")
	}

	// withdraw
	tx, err = withdrawer.RunWithdraws(
		r.ZEVMAuth,
		r.DeployerAddress.Bytes(),
		r.ERC20ZRC20Addr,
		withdrawalAmount,
		numberOfWithdrawals,
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}

	cctxs := utils.WaitCctxsMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, int(numberOfWithdrawals.Int64()), r.Logger, r.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// verify the withdraw value
	for _, cctx := range cctxs {
		verifyTransferAmountFromCCTX(r, cctx, withdrawalAmount.Int64())
	}
}
