package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestEtherWithdraw tests the withdrawal of ether
func TestEtherWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok, "Invalid withdrawal amount specified for TestEtherWithdraw.")

	// add one unit to the withdrawal amount to account for the approval because withdrawal fees are automatically deducted
	approvalAmount := new(big.Int).Add(withdrawalAmount, big.NewInt(1e18))

	// approve
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, approvalAmount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	tx = r.WithdrawEther(withdrawalAmount)

	// verify the withdrawal value
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")

	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	r.Logger.Info("TestEtherWithdraw completed")
}
