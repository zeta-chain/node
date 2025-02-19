package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestEVMToZEVMCallAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// perform the withdraw
	tx := r.EVMToZEMVCall(
		sample.EthAddress(), // non-existing address
		[]byte("revert"),
		gatewayevm.RevertOptions{
			CallOnRevert:  false, // revert not supported
			RevertMessage: []byte("revert"),
			AbortAddress:  testAbortAddr,
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_Aborted, cctx.CctxStatus.Status)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
	require.NoError(r, err)
	require.EqualValues(r, r.ERC20ZRC20Addr.Hex(), abortContext.Asset.Hex())
	require.EqualValues(r, int64(0), abortContext.Amount.Int64())
}
