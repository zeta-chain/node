package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestZetaWithdrawAndArbitraryCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	payload := randomPayload(r)
	//payload := strings.ToLower(r.ZetaEthAddr.String())
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.AssertTestDAppEVMCalled(false, payload, amount)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// use a contract as revert recipient so the post-revert balance check
	// isn't perturbed by gas costs (the revert tx is signed by TSS).
	revertAddress := r.TestDAppV2ZEVMAddr

	// perform the withdraw
	tx := r.ZETAWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		evmChainID,
		r.EncodeERC20Call(r.ZetaEthAddr, amount, payload),
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: ZetaConnector.withdrawAndCall arbitrary calls
		// route to GatewayEVM.executeWithERC20's sender==0 arbitrary-call branch
		// — the same drain shape as GatewayEVM.execute. The signer cancels these
		// CCTXs (TSS self-transfer); the CCTX terminates as Reverted via the
		// standard V2 revert flow; the destination dApp must not be invoked.

		// wait for the withdraw tx to mine on ZEVM, then capture the
		// post-payment ZEVM-coin balance of the revert recipient.
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)
		revertBalanceBefore, err := r.ZEVMClient.BalanceAt(r.Ctx, revertAddress, nil)
		require.NoError(r, err)

		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "zeta_withdraw_and_arbitrary_call")
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

		// destination dApp must not have been called
		r.AssertTestDAppEVMCalled(false, payload, amount)

		// revert recipient should receive at least the principal back
		revertBalanceAfter, err := r.ZEVMClient.BalanceAt(r.Ctx, revertAddress, nil)
		require.NoError(r, err)
		require.EqualValues(r, new(big.Int).Add(revertBalanceBefore, amount), revertBalanceAfter)
	} else {
		// V2 ZETA flows disabled: tx should revert on GatewayZEVM, no CCTX created
		utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)
	}
}
