package e2etests

import (
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

const payloadMessageWithdrawERC20 = "this is a test ERC20 withdraw and call payload"

func TestV2ERC20WithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ERC20WithdrawAndCall")

	r.AssertTestDAppEVMValues(false, payloadMessageWithdrawETH, amount)

	r.ApproveERC20ZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.V2ERC20WithdrawAndCall(r.EVMAddress(), amount, r.EncodeERC20Call(r.ERC20ZRC20Addr, amount, payloadMessageWithdrawERC20))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	r.AssertTestDAppEVMValues(true, payloadMessageWithdrawERC20, amount)
}
