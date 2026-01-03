package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestZetaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])
	receiverAddress := r.EVMAddress()

	r.ApproveZetaOnEVM(r.GatewayEVMAddr)
	// perform the deposit
	tx := r.ZETADeposit(receiverAddress, amount, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta_deposit")

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: deposit should succeed
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
		oldBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, receiverAddress, nil)
		require.NoError(r, err)
		newBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, receiverAddress, nil)
		require.NoError(r, err)
		require.True(r, newBalance.Cmp(oldBalance) >= 0)
	} else {
		// V2 ZETA flows disabled: deposit should be aborted
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)
		require.Equal(r, cctx.CctxStatus.StatusMessage, crosschaintypes.ErrZetaThroughGateway.Error())
	}
}
