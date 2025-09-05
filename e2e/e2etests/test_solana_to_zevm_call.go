package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaToZEVMCall tests calling an example contract
func TestSolanaToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// ARRANGE
	// Given payload and ZEVM contract address
	contractAddr := r.TestDAppV2ZEVMAddr
	payload := randomPayload(r)
	r.AssertTestDAppZEVMCalled(false, payload, nil)

	// ACT
	// execute call transaction
	sig := r.SOLCall(nil, contractAddr, []byte(payload), nil)

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, contractAddr.Hex())

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, big.NewInt(0))
}
