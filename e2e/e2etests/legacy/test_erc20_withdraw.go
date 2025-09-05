package legacy

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestERC20Withdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok, "Invalid withdrawal amount specified for TestERC20Withdraw.")

	// approve 1 unit of the gas token to cover the gas fee
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// withdraw
	tx = r.LegacyWithdrawERC20(withdrawalAmount)

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	r.EVMVerifyOutboundTransferAmount(cctx.GetCurrentOutboundParam().Hash, withdrawalAmount.Int64())
}
