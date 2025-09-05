package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHWithdrawCustomGasLimit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	amount := utils.ParseBigInt(r, args[0])
	customGasLimit := utils.ParseBigInt(r, args[1])

	oldBalance, err := r.EVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	require.NoError(r, err)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw with the custom gas limit option
	tx, err := r.GatewayZEVM.Withdraw1(
		r.ZEVMAuth,
		r.EVMAddress().Bytes(),
		amount,
		r.ETHZRC20Addr,
		customGasLimit,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check gas limit used
	require.EqualValues(r, customGasLimit.Uint64(), cctx.OutboundParams[0].CallOptions.GasLimit)

	// check the balance was updated, we just check newBalance is greater than oldBalance because of the gas fee
	newBalance, err := r.EVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	require.NoError(r, err)
	require.Greater(r, newBalance.Uint64(), oldBalance.Uint64())
}
