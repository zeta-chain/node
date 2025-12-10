package e2etests

import (
	"math/big"

	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	// crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestZetaWithdrawRevertAndAbort tests that ZETA withdraw revert and abort through gateway
// is not supported in V2 - no CCTX should be created.
func TestZetaWithdrawRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	_ = testAbort

	// perform the withdraw
	tx := r.ZETAWithdrawAndCall(
		sample.EthAddress(), // non-existing address
		amount,
		[]byte("revert"),
		evmChainID,
		gatewayzevm.RevertOptions{
			RevertAddress:    sample.EthAddress(), // non-existing address
			CallOnRevert:     true,
			RevertMessage:    []byte("revert"),
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		},
		gasLimit,
	)

	// ZETA withdraws through gateway are not supported in V2, verify no CCTX is created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)

	// // wait for the cctx to be mined
	// cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	// r.Logger.CCTX(*cctx, "withdraw")
	// utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)
	//
	// // check onAbort was called
	// aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	// require.NoError(r, err)
	// require.True(r, aborted)
	//
	// // check abort context was passed
	// abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
	// require.NoError(r, err)
	// require.EqualValues(r, common.Address{}.Hex(), abortContext.Asset.Hex())
	//
	// newBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, testAbortAddr, nil)
	// require.NoError(r, err)
	// require.True(r, newBalance.Cmp(big.NewInt(0)) > 0)
}
