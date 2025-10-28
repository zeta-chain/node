package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHDepositAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	// use a random address to get the revert amount
	revertAddress := sample.EthAddress()
	balance, err := r.EVMClient.BalanceAt(r.Ctx, revertAddress, nil)
	require.NoError(r, err)
	require.EqualValues(r, int64(0), balance.Int64())

	// perform the deposit
	tx := r.ETHDepositAndCall(r.TestDAppV2ZEVMAddr, amount, []byte("revert"), gatewayevm.RevertOptions{
		RevertAddress:    revertAddress,
		OnRevertGasLimit: big.NewInt(0),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)

	// check the balance is more than 0
	balance, err = r.EVMClient.BalanceAt(r.Ctx, revertAddress, nil)
	require.NoError(r, err)
	require.True(r, balance.Cmp(big.NewInt(0)) > 0)
}
