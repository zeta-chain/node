package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testgasconsumer"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestDepositAndCallOutOfGas tests that a deposit and call that consumer all gas will revert
func TestDepositAndCallOutOfGas(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ETHDepositAndCall")

	// Deploy the GasConsumer contract
	gasConsumerAddress, _, _, err := testgasconsumer.DeployTestGasConsumer(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// perform the deposit and call to the GasConsumer contract
	tx := r.V2ETHDepositAndCall(
		gasConsumerAddress,
		amount,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)
}
