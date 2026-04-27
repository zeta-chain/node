package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20WithdrawAndArbitraryCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	payload := randomPayload(r)

	r.AssertTestDAppEVMCalled(false, payload, amount)

	r.ApproveERC20ZRC20(r.GatewayZEVMAddr)
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ERC20WithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		r.EncodeERC20Call(r.ERC20Addr, amount, payload),
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")

	// V2 arbitrary calls are cancelled by the signer (TSS self-transfer).
	// The observer reports the outbound as failed, which routes the CCTX
	// through the standard V2 revert flow and terminates as Reverted; the
	// withdrawn ERC20 ZRC20 amount is refunded to the sender on ZEVM.
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// destination dApp must not have been called
	r.AssertTestDAppEVMCalled(false, payload, amount)
}
