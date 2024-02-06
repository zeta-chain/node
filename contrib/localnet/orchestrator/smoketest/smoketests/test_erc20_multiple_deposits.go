package smoketests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func TestMultipleERC20Deposit(sm *runner.SmokeTestRunner) {
	initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash := MultipleDeposits(sm, big.NewInt(1e9), big.NewInt(3))
	cctxs := utils.WaitCctxsMinedByInTxHash(sm.Ctx, txhash.Hex(), sm.CctxClient, 3, sm.Logger, sm.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// check new balance is increased by 1e9 * 3
	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	diff := big.NewInt(0).Sub(bal, initialBal)
	if diff.Int64() != 3e9 {
		panic(fmt.Sprintf("balance difference is not correct: %d", diff.Int64()))
	}
}

func MultipleDeposits(sm *runner.SmokeTestRunner, amount, count *big.Int) ethcommon.Hash {
	// deploy depositor
	depositorAddr, _, depositor, err := testcontract.DeployDepositor(sm.GoerliAuth, sm.GoerliClient, sm.ERC20CustodyAddr)
	if err != nil {
		panic(err)
	}

	fullAmount := big.NewInt(0).Mul(amount, count)

	// approve
	tx, err := sm.USDTERC20.Approve(sm.GoerliAuth, depositorAddr, fullAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	// deposit
	tx, err = depositor.RunDeposits(sm.GoerliAuth, sm.DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, []byte{}, count)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposits failed")
	}
	sm.Logger.Info("Deposits receipt tx hash: %s", tx.Hash().Hex())

	for _, log := range receipt.Logs {
		event, err := sm.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info("Multiple deposit event: ")
		sm.Logger.Info("  Amount: %d, ", event.Amount)
	}
	sm.Logger.Info("gas limit %d", sm.ZevmAuth.GasLimit)
	return tx.Hash()
}
