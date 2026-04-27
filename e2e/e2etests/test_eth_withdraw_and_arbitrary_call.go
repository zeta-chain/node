package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHWithdrawAndArbitraryCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	payload := randomPayload(r)

	r.AssertTestDAppEVMCalled(false, payload, amount)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	revertAddress := r.EVMAddress()

	// perform the withdraw
	tx := r.ETHWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		r.EncodeGasCall(payload),
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the withdraw tx to mine on ZEVM (user has now paid amount +
	// gasFee + protocolFlatFee out of their ETH-ZRC20 balance), then capture
	// the post-payment balance — the revert refund will add on top of this.
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	revertBalanceBefore, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")

	// V2 arbitrary calls are cancelled by the signer (TSS self-transfer).
	// The observer reports the outbound as failed, which routes the CCTX
	// through the standard V2 revert flow and terminates as Reverted; the
	// withdrawn ZRC20 amount is refunded to the revert address on ZEVM.
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// the cancelled outbound is a TSS self-transfer with zero value
	r.EVMVerifyOutboundTransferAmount(cctx.OutboundParams[0].Hash, 0)

	// destination dApp must not have been called
	r.AssertTestDAppEVMCalled(false, payload, amount)

	// Compared to the post-withdraw balance, the revert refund must add at
	// least the principal (zetacore's ProcessRevert refunds amount in full).
	// The exact value also includes a fractional gas-fee leftover refund from
	// RefundUnusedGasFee, but the precise math is fragile across percentage
	// settings, so we assert the lower bound only.
	revertBalanceAfter, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	wantMin := new(big.Int).Add(revertBalanceBefore, amount)
	require.GreaterOrEqual(
		r,
		revertBalanceAfter.Cmp(wantMin),
		0,
		"revert address should receive at least the principal: got %s, want >= %s",
		revertBalanceAfter.String(),
		wantMin.String(),
	)
}
