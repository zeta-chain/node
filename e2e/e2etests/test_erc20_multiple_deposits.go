package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func TestMultipleERC20Deposit(r *runner.E2ERunner) {
	initialBal, err := r.USDTZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash := MultipleDeposits(r, big.NewInt(1e9), big.NewInt(3))
	cctxs := utils.WaitCctxsMinedByInTxHash(r.Ctx, txhash.Hex(), r.CctxClient, 3, r.Logger, r.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// check new balance is increased by 1e9 * 3
	bal, err := r.USDTZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	diff := big.NewInt(0).Sub(bal, initialBal)
	if diff.Int64() != 3e9 {
		panic(fmt.Sprintf("balance difference is not correct: %d", diff.Int64()))
	}
}

func MultipleDeposits(r *runner.E2ERunner, amount, count *big.Int) ethcommon.Hash {
	// deploy depositor
	depositorAddr, _, depositor, err := testcontract.DeployDepositor(r.GoerliAuth, r.GoerliClient, r.ERC20CustodyAddr)
	if err != nil {
		panic(err)
	}

	fullAmount := big.NewInt(0).Mul(amount, count)

	// approve
	tx, err := r.USDTERC20.Approve(r.GoerliAuth, depositorAddr, fullAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.GoerliClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	// deposit
	tx, err = depositor.RunDeposits(r.GoerliAuth, r.DeployerAddress.Bytes(), r.USDTERC20Addr, amount, []byte{}, count)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.GoerliClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposits failed")
	}
	r.Logger.Info("Deposits receipt tx hash: %s", tx.Hash().Hex())

	for _, log := range receipt.Logs {
		event, err := r.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("Multiple deposit event: ")
		r.Logger.Info("  Amount: %d, ", event.Amount)
	}
	r.Logger.Info("gas limit %d", r.ZevmAuth.GasLimit)
	return tx.Hash()
}
