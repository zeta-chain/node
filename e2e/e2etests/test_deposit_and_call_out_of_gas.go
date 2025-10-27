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

	// Deploy the GasConsumer contract
	// gas limit is currently 4M
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
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)
}
