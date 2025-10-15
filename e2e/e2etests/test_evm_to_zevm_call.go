package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestEVMToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	payload := randomPayload(r)
	sender := r.EVMAddress().Bytes()

	r.AssertTestDAppZEVMCalled(false, payload, sender, big.NewInt(0))

	// perform the withdraw
	tx := r.EVMToZEMVCall(
		r.TestDAppV2ZEVMAddr,
		[]byte(payload),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, sender, big.NewInt(0))
}
