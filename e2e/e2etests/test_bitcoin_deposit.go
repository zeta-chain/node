package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	depositAmount := utils.ParseFloat(r, args[0])

	txHash := r.DepositBTCWithExactAmount(depositAmount, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
