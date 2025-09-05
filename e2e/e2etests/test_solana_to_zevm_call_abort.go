package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSolanaToZEVMCallAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// execute call transaction, receiver is non existing address
	data := []byte("hello")
	sig := r.SOLCall(nil, sample.EthAddress(), data, &solana.RevertOptions{
		OnRevertGasLimit: 0,
		AbortAddress:     testAbortAddr,
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)
}
