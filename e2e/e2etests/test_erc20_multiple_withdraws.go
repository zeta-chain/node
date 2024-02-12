package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func TestMultipleWithdraws(r *runner.E2ERunner) {
	// deploy withdrawer
	withdrawerAddr, _, withdrawer, err := testcontract.DeployWithdrawer(r.ZevmAuth, r.ZevmClient)
	if err != nil {
		panic(err)
	}

	// approve
	tx, err := r.USDTZRC20.Approve(r.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.Info("USDT ZRC20 approve receipt: status %d", receipt.Status)

	// approve gas token
	tx, err = r.ETHZRC20.Approve(r.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve gas token failed")
	}
	r.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// check the balance
	bal, err := r.USDTZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("balance of deployer on USDT ZRC20: %d", bal)

	if bal.Int64() < 1000 {
		panic("not enough USDT ZRC20 balance!")
	}

	// withdraw
	tx, err = withdrawer.RunWithdraws(
		r.ZevmAuth,
		r.DeployerAddress.Bytes(),
		r.USDTZRC20Addr,
		big.NewInt(100),
		big.NewInt(3),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}

	cctxs := utils.WaitCctxsMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, 3, r.Logger, r.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// verify the withdraw value
	for _, cctx := range cctxs {
		verifyTransferAmountFromCCTX(r, cctx, 100)
	}
}
