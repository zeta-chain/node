package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestZetaDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveZetaOnEVM(r.GatewayEVMAddr)

	payload := randomPayload(r)
	sender := r.EVMAddress().Bytes()
	receiverAddress := r.TestDAppV2ZEVMAddr

	r.AssertTestDAppZEVMCalled(false, payload, sender, amount)

	oldBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, receiverAddress, nil)
	require.NoError(r, err)

	// perform the deposit
	tx := r.ZetaDepositAndCall(
		receiverAddress,
		amount,
		[]byte(payload),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta_deposit_and_call")

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: deposit and call should succeed
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

		// check the payload was received on the contract
		r.AssertTestDAppZEVMCalled(true, payload, sender, amount)

		// check the balance was updated
		newBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, receiverAddress, nil)
		require.NoError(r, err)
		require.Equal(r, new(big.Int).Add(oldBalance, amount), newBalance)
	} else {
		// V2 ZETA flows disabled: deposit should be aborted
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)
		require.Contains(r, cctx.CctxStatus.StatusMessage, crosschaintypes.ErrZetaThroughGateway.Error())
	}
}
