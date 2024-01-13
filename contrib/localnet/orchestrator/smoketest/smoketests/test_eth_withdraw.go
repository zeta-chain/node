package smoketests

import (
	"math/big"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestEtherWithdraw tests the withdraw of ether
func TestEtherWithdraw(sm *runner.SmokeTestRunner) {
	// approve
	tx, err := sm.ETHZRC20.Approve(sm.ZevmAuth, sm.USDTZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	sm.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	tx, err = sm.ETHZRC20.Withdraw(sm.ZevmAuth, sm.DeployerAddress.Bytes(), big.NewInt(100000))
	if err != nil {
		panic(err)
	}
	sm.Logger.EVMTransaction(*tx, "withdraw")

	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}
	sm.Logger.EVMReceipt(*receipt, "withdraw")

	for _, log := range receipt.Logs {
		event, err := sm.ETHZRC20.ParseWithdrawal(*log)
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
	sm.Logger.CCTX(*cctx, "withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic("cctx status is not outbound mined")
	}
}
