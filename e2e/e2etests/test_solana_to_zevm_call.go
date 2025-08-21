package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaToZEVMCall tests calling an example contract
func TestSolanaToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Example contract deployed at: %s", contractAddr.String())

	// execute call transaction
	data := []byte("hello")
	sig := r.SOLCall(nil, contractAddr, data, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, contractAddr.Hex())

	// check if example contract has been called, bar value should be set to amount
	utils.WaitAndVerifyExampleContractCallWithMsg(
		r,
		contract,
		big.NewInt(0),
		data,
		[]byte(r.GetSolanaPrivKey().PublicKey().String()),
	)
}
