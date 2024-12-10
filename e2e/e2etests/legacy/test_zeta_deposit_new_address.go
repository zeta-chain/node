package legacy

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestZetaDepositNewAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount
	amount := utils.ParseBigInt(r, args[0])

	newAddress := sample.EthAddress()
	hash := r.LegacyDepositZetaWithAmount(newAddress, amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}
