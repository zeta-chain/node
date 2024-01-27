package smoketests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestWithdrawERC20(sm *runner.SmokeTestRunner) {
	// approve
	tx, err := sm.ETHZRC20.Approve(sm.ZevmAuth, sm.USDTZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// withdraw
	tx, err = sm.USDTZRC20.Withdraw(sm.ZevmAuth, sm.DeployerAddress.Bytes(), big.NewInt(1000))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	sm.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.USDTZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info(
			"  logs: from %s, to %x, value %d, gasfee %d",
			event.From.Hex(),
			event.To,
			event.Value,
			event.Gasfee,
		)
	}

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, receipt.TxHash.Hex(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	verifyTransferAmountFromCCTX(sm, cctx, 1000)
}

func TestMultipleWithdraws(sm *runner.SmokeTestRunner) {
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
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT ZRC20 approve receipt: status %d", receipt.Status)

	// approve gas token
	tx, err = sm.ETHZRC20.Approve(sm.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
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
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}

	cctxs := utils.WaitCctxsMinedByInTxHash(sm.Ctx, tx.Hash().Hex(), sm.CctxClient, 3, sm.Logger, sm.CctxTimeout)
	if len(cctxs) != 3 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// verify the withdraw value
	for _, cctx := range cctxs {
		verifyTransferAmountFromCCTX(sm, cctx, 100)
	}
}

// verifyTransferAmountFromCCTX verifies the transfer amount from the CCTX on Goerli
func verifyTransferAmountFromCCTX(sm *runner.SmokeTestRunner, cctx *crosschaintypes.CrossChainTx, amount int64) {
	sm.Logger.Info("outTx hash %s", cctx.GetCurrentOutTxParam().OutboundTxHash)

	receipt, err := sm.GoerliClient.TransactionReceipt(
		sm.Ctx,
		ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash),
	)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.USDTERC20.ParseTransfer(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info("  logs: from %s, to %s, value %d", event.From.Hex(), event.To.Hex(), event.Value)
		if event.Value.Int64() != amount {
			panic("value is not correct")
		}
	}
}
