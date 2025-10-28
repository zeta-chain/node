package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testdappempty"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHWithdrawAndCallBigPayload(r *runner.E2ERunner, _ []string) {
	// deploy the TestDAppEmpty contract on the EVM chain
	testDAppAddr, txDeploy, _, err := testdappempty.DeployTestDAppEmpty(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnEVM(txDeploy)

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	// create a random payload with 2880 bytes which is current max in the gateway
	payload := randomPayloadWithSize(r, 2880)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ETHWithdrawAndCall(
		testDAppAddr,
		big.NewInt(1),
		[]byte(payload),
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		big.NewInt(200000),
	)

	r.Logger.EVMTransaction(tx, "withdraw and call big payload")

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
