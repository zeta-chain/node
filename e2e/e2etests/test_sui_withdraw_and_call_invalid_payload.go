package e2etests

import (
	"encoding/hex"
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawAndCallInvalidPayload executes withdrawAndCall on zevm gateway with invalid payload.
// The outbound authenticated call will be cancelled by the zetaclient.
func TestSuiWithdrawAndCallInvalidPayload(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package), SUI amount and gas limit
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])
	gasLimit := big.NewInt(100000)

	// create an invalid 'on_call' payload that cannot be unpacked by zetaclient
	message, err := hex.DecodeString("deadbeef")
	require.NoError(r, err)

	// given TSS balance in Sui network
	tssBalanceBefore := r.SuiGetSUIBalance(r.SuiTSSAddress)

	// ACT
	// approve SUI ZRC20 token
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and authenticated call with revert options
	tx := r.SuiWithdrawAndCall(
		targetPackageID,
		amount,
		r.SUIZRC20Addr,
		message,
		gasLimit,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	r.Logger.EVMTransaction(tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// the TSS balance in Sui network should be higher or equal to the balance before
	// reason is that the max budget is refunded to the TSS
	tssBalanceAfter := r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore)
}
