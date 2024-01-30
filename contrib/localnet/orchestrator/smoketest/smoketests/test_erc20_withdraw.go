package smoketests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestERC20Withdraw(sm *runner.SmokeTestRunner) {
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
