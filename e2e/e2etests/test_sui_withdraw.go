package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	balanceBefore := r.SuiGetSUIBalance(signer.Address())

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.SuiWithdrawSUI(signer.Address(), amount)
	r.Logger.EVMTransaction(*tx, "withdraw")

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the balance after the withdraw
	balanceAfter := r.SuiGetSUIBalance(signer.Address())
	require.EqualValues(r, balanceBefore+amount.Uint64(), balanceAfter)
}
