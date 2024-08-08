package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const payloadMessageERC20 = "this is a test ERC20 deposit and call payload"

func TestV2ERC20DepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ERC20DepositAndCall")

	r.ApproveERC20(r.GatewayEVMAddr)

	r.AssertTestDAppValues(false, payloadMessageERC20, amount)

	// perform the deposit
	tx := r.V2ERC20DepositAndCall(r.TestDAppV2ZEVMAddr, amount, []byte(payloadMessageERC20))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppValues(true, payloadMessageERC20, amount)
}
