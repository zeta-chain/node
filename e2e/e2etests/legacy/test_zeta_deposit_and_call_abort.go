package legacy

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// TestZetaDepositAndCallAbort tests a deposit with a payload that causes the cctx to be aborted.
func TestZetaDepositAndCallAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount
	amount := utils.ParseBigInt(r, args[0])

	hash := r.LegacyDepositZetaWithAmountAndPayload(r.ZevmTestDAppAddr, amount, []byte("test payload"))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")

	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Aborted, "cctx should be aborted")
}
