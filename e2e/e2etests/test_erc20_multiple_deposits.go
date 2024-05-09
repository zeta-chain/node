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

func TestMultipleERC20Deposit(r *runner.E2ERunner, args []string) {
	if len(args) != 2 {
		panic("TestMultipleERC20Deposit requires exactly two arguments: the deposit amount and the number of deposits.")
	}

	depositAmount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid deposit amount specified for TestMultipleERC20Deposit.")
	}

	numberOfDeposits, ok := big.NewInt(0).SetString(args[1], 10)
	if !ok || numberOfDeposits.Int64() < 1 {
		panic("Invalid number of deposits specified for TestMultipleERC20Deposit.")
	}

	initialBal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash := MultipleDeposits(r, depositAmount, numberOfDeposits)
	cctxs := utils.WaitCctxsMinedByInboundHash(r.Ctx, txhash.Hex(), r.CctxClient, int(numberOfDeposits.Int64()), r.Logger, r.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// check new balance is increased by amount * count
	bal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	diff := big.NewInt(0).Sub(bal, initialBal)
	total := depositAmount.Mul(depositAmount, numberOfDeposits)
	if diff.Cmp(total) != 0 {
		panic(fmt.Sprintf("balance difference is not correct: %d", diff.Int64()))
	}
}

func MultipleDeposits(r *runner.E2ERunner, amount, count *big.Int) ethcommon.Hash {
	// deploy depositor
	depositorAddr, _, depositor, err := testcontract.DeployDepositor(r.EVMAuth, r.EVMClient, r.ERC20CustodyAddr)
	if err != nil {
		panic(err)
	}

	fullAmount := big.NewInt(0).Mul(amount, count)

	// approve
	tx, err := r.ERC20.Approve(r.EVMAuth, depositorAddr, fullAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	// deposit
	tx, err = depositor.RunDeposits(r.EVMAuth, r.DeployerAddress.Bytes(), r.ERC20Addr, amount, []byte{}, count)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
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
	r.Logger.Info("gas limit %d", r.ZEVMAuth.GasLimit)
	return tx.Hash()
}
