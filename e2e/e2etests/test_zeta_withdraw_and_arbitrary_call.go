package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	// crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestZetaWithdrawAndArbitraryCall tests that ZETA withdraw and arbitrary call through gateway
// is not supported in V2 - no CCTX should be created.
func TestZetaWithdrawAndArbitraryCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	payload := randomPayload(r)
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// r.AssertTestDAppEVMCalled(false, payload, amount)

	// perform the withdraw
	tx := r.ZETAWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		evmChainID,
		r.EncodeERC20Call(r.ZetaEthAddr, amount, payload),
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// ZETA withdraws through gateway are not supported in V2, verify no CCTX is created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)

	// // wait for the cctx to be mined
	// cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	// r.Logger.CCTX(*cctx, "withdraw")
	// utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	// r.AssertTestDAppEVMCalled(true, payload, amount)
}
