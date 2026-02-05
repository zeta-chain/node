package legacy

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestZetaWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse withdraw amount
	amount := utils.ParseBigInt(r, args[0])

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.LegacyDepositAndApproveWZeta(amount)
	tx := r.LegacyWithdrawZeta(amount, true)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// Get chain params for stability pool percentage
	chainParams, err := r.ObserverClient.GetChainParamsForChain(
		r.Ctx,
		&observertypes.QueryGetChainParamsForChainRequest{ChainId: evmChainID.Int64()},
	)
	require.NoError(r, err)

	// Verify gas accounting and get refund amounts
	refunds := utils.VerifyOutboundGasAccounting(r, cctx, chainParams.ChainParams.StabilityPoolPercentage)
	r.Logger.Info("Gas refund - StabilityPool: %s, UserRefund: %s",
		refunds.StabilityPoolAmount.String(), refunds.UserRefundAmount.String())
}
