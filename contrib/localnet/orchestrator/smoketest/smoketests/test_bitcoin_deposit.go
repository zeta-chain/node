package smoketests

import (
	"fmt"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestBitcoinDeposit(sm *runner.SmokeTestRunner) {

	sm.SetBtcAddress(sm.Name, false)

	txHash := sm.DepositBTCWithAmount(0.001)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, txHash.String(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	sm.Logger.CCTX(*cctx, "deposit")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected mined status; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}
}
