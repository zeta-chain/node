package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	oldBalance, err := r.EVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	require.NoError(r, err)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ETHWithdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check the balance was updated, we just check newBalance is greater than oldBalance because of the gas fee
	newBalance, err := r.EVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	require.NoError(r, err)
	require.Greater(r, newBalance.Uint64(), oldBalance.Uint64())
}
