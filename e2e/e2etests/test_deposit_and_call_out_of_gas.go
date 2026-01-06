package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testgasconsumer"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestDepositAndCallOutOfGas tests that a deposit and call that consumer all gas will revert
func TestDepositAndCallOutOfGas(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	// Update the gateway gas limit to 4M
	r.UpdateGatewayGasLimit(uint64(4_000_000))
	// Deploy the GasConsumer contract with a gas limit of 5M
	gasConsumerAddress, txDeploy, _, err := testgasconsumer.DeployTestGasConsumer(
		r.ZEVMAuth,
		r.ZEVMClient,
		big.NewInt(5000000),
	)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	// perform the deposit and call to the GasConsumer contract
	tx := r.ETHDepositAndCall(
		gasConsumerAddress,
		amount,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
}
