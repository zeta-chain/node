package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20Withdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20ZRC20(r.GatewayZEVMAddr)
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ERC20Withdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
