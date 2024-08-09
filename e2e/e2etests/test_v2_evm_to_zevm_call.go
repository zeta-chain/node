package e2etests

import (
	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestV2EVMToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// TODO: set payload
	payload := []byte("")

	// perform the withdraw
	tx := r.V2EVMToZEMVCall(r.EVMAddress(), payload)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
