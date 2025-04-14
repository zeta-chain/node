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
	tssBalanceBefore := r.SuiGetSUIBalance(r.SuiTSSAddress)

	amount := utils.ParseBigInt(r, args[0])

	r.SuiApproveSUIZRC20(r.GatewayZEVMAddr)

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

	// check the TSS balance after transaction is higher or equal to the balance before
	// reason is that the max budget is refunded to the TSS
	tssBalanceAfter := r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore)

	// PATCH: v29

	tssBalanceBefore = tssBalanceAfter

	// Check that an invalid withdraw doesn't block the outbound
	// Use the same address as in the current incident
	tx = r.SuiWithdrawSUI("0x307832356462313663336361353535663637303263303738363035303331303762623733636365396636633164366466303034363435323964623135643561356162", amount)
	r.Logger.EVMTransaction(*tx, "invalid_withdraw")

	// wait for the cctx to be mined
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the TSS receive the amount
	tssBalanceAfter = r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore+amount.Uint64())
}
