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

	// perform the withdraw
	tx := r.ZETAWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		evmChainID,
		r.EncodeERC20Call(r.ZetaEthAddr, amount, payload),
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: withdraw and call should succeed
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "zeta_withdraw_and_arbitrary_call")
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

		r.AssertTestDAppEVMCalled(true, payload, amount)
	} else {
		// V2 ZETA flows disabled: tx should revert on GatewayZEVM, no CCTX created
		utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)
	}
}
