package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testdappempty"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHDepositAndCallBigPayload(r *runner.E2ERunner, _ []string) {
	// deploy the TestDAppEmpty contract on the ZetaChain
	testDAppAddr, _, _, err := testdappempty.DeployTestDAppEmpty(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	previousGasLimit := r.EVMAuth.GasLimit
	r.EVMAuth.GasLimit = 10000000
	defer func() {
		r.EVMAuth.GasLimit = previousGasLimit
	}()

	// create a random payload with 2880 bytes which is current max in the gateway
	payload := randomPayloadWithSize(r, 2880)

	// perform the withdraw
	tx := r.ETHDepositAndCall(
		testDAppAddr,
		big.NewInt(1),
		[]byte(payload),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
