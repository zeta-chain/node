package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHMultipleDeposits(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.Logger.Info("starting eth multiple deposits test")

	oldBalance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// set value of the payable transactions
	previousValue := r.EVMAuth.Value
	r.EVMAuth.Value = amount

	// send two deposit through contract
	r.Logger.Print("🏃test two deposits through contract")
	tx, err := r.TestDAppV2EVM.GatewayTwoDeposits(r.EVMAuth, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Print("🍾 multiple deposits through contract observed")

	r.EVMAuth.Value = previousValue

	// wait for the cctxs to be mined
	cctx := utils.WaitCctxsMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, 2, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx[0], "deposit 1")
	r.Logger.CCTX(*cctx[1], "deposit 2")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx[0].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx[1].CctxStatus.Status)

	// wait for the zrc20 balance to be updated
	fee, err := r.GatewayEVM.AdditionalActionFeeWei(nil)
	require.NoError(r, err)
	change := utils.NewExactChange(amount.Sub(amount, fee))
	utils.WaitAndVerifyZRC20BalanceChange(r, r.ETHZRC20, r.EVMAddress(), oldBalance, change, r.Logger)
}
