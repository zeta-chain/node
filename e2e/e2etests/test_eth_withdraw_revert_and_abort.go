package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHWithdrawRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// deploy testabort contract
	testAbortAddr, txDeploy, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	// perform the withdraw
	tx := r.ETHWithdrawAndCall(
		sample.EthAddress(), // non-existing address
		amount,
		[]byte("revert"),
		gatewayzevm.RevertOptions{
			RevertAddress: sample.EthAddress(), // non-existing address
			CallOnRevert:  true,
			RevertMessage: []byte(
				"withdraw",
			), // withdraw is passed as message to create a withdraw in onAbort and test cctx can be created
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		},
		gasLimit,
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "withdraw")
	require.NoError(r, err)
	require.EqualValues(r, r.ETHZRC20Addr.Hex(), abortContext.Asset.Hex())

	// check the create withdraw get mined
	cctxWithdrawFromAbort := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		cctx.Index,
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)

	// check the cctx status
	utils.RequireCCTXStatus(r, cctxWithdrawFromAbort, crosschaintypes.CctxStatus_OutboundMined)
}
