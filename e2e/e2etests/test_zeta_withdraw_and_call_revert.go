package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestZetaWithdrawAndCallRevert tests ZETA withdraw and call revert through gateway
func TestZetaWithdrawAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// use a random address to get the revert amount
	revertAddress := sample.EthAddress()

	// perform the withdraw
	tx := r.ZETAWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		evmChainID,
		r.EncodeERC20CallRevert(r.ERC20Addr, amount),
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: withdraw and call should revert
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "zeta_withdraw_and_call_revert")
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

		newBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, revertAddress, nil)
		require.NoError(r, err)
		require.True(r, newBalance.Cmp(big.NewInt(0)) > 0)
	} else {
		// V2 ZETA flows disabled: tx should revert on GatewayZEVM, no CCTX created
		utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)
	}
}
