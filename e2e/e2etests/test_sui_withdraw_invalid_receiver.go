package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestSuiWithdrawInvalidReceiver tests that a withdrawal to a invalid receiver address that fails in the ZEVM
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
	tx := r.SuiWithdraw(
		receiver,
		amount,
		r.SUIZRC20Addr,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(tx, "withdraw to invalid sui address")

	// ASSERT
	// withdraw tx should fail in ZEVM
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)
}
