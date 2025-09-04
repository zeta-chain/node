package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaDepositAndCall tests deposit of lamports calling a example contract
func TestSolanaDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// parse deposit amount (in lamports)
	amount := utils.ParseBigInt(r, args[0])

	// Given payload and ZEVM contract address
	payload := randomPayload(r)
	contractAddr := r.TestDAppV2ZEVMAddr
	r.AssertTestDAppZEVMCalled(false, payload, amount)

	// ACT
	// execute the deposit transaction
	sig := r.SOLDepositAndCall(nil, contractAddr, amount, []byte(payload), nil)

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, contractAddr.Hex())

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, amount)
}
