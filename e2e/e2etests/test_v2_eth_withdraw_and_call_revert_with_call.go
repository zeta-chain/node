package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TODO
func TestV2ETHWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ETHWithdrawAndCall")

	r.AssertTestDAppEVMValues(false, payloadMessageWithdrawETH, amount)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.V2ETHWithdrawAndCall(r.TestDAppV2EVMAddr, amount, r.EncodeGasCall(payloadMessageWithdrawETH))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	r.AssertTestDAppEVMValues(true, payloadMessageWithdrawETH, amount)
}
