package e2etests

import (
	"fmt"
	"github.com/zeta-chain/zetacore/e2e/runner"
	utils2 "github.com/zeta-chain/zetacore/e2e/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func TestMultipleWithdraws(sm *runner.E2ERunner) {
	// deploy withdrawer
	withdrawerAddr, _, withdrawer, err := testcontract.DeployWithdrawer(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// approve
	tx, err := sm.USDTZRC20.Approve(sm.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils2.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT ZRC20 approve receipt: status %d", receipt.Status)

	// approve gas token
	tx, err = sm.ETHZRC20.Approve(sm.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils2.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve gas token failed")
	}
	sm.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// check the balance
	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("balance of deployer on USDT ZRC20: %d", bal)

	if bal.Int64() < 1000 {
		panic("not enough USDT ZRC20 balance!")
	}

	// withdraw
	tx, err = withdrawer.RunWithdraws(
		sm.ZevmAuth,
		sm.DeployerAddress.Bytes(),
		sm.USDTZRC20Addr,
		big.NewInt(100),
		big.NewInt(3),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils2.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}

	cctxs := utils2.WaitCctxsMinedByInTxHash(sm.Ctx, tx.Hash().Hex(), sm.CctxClient, 3, sm.Logger, sm.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// verify the withdraw value
	for _, cctx := range cctxs {
		verifyTransferAmountFromCCTX(sm, cctx, 100)
	}
}
