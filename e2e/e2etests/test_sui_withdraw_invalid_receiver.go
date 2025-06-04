package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawInvalidReceiver tests that a withdrawal to a invalid receiver address that reverts
func TestSuiWithdrawInvalidReceiver(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given amount, receiver, revert address
	receiver := args[0]
	amount := utils.ParseBigInt(r, args[1])
	revertAddress := r.EVMAddress()

	// approve the ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// ACT
	// perform the withdraw to invalid receiver
	tx := r.SuiWithdrawSUI(
		receiver,
		amount,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(*tx, "withdraw to invalid sui address")

	// wait for the withdraw tx to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// revert address balance before
	revertBalanceBefore, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// revert address should receive the amount
	revertBalanceAfter, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.EqualValues(r, new(big.Int).Add(revertBalanceBefore, amount), revertBalanceAfter)
}
