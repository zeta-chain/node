package e2etests

import (
	"math/big"

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

	oldBalance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	// set value of the payable transactions
	previousValue := r.EVMAuth.Value
	fee, err := r.GatewayEVM.AdditionalActionFeeWei(nil)
	require.NoError(r, err)
	// add 5 fees to provided amount to pay for 6 inbounds (1st one is free)
	r.EVMAuth.Value = new(big.Int).Add(amount, new(big.Int).Mul(fee, big.NewInt(5)))
	defer func() {
		r.EVMAuth.Value = previousValue
	}()

	// send multiple deposit through contract
	tx, err := r.TestDAppV2EVM.GatewayMultipleDeposits(r.EVMAuth, r.TestDAppV2ZEVMAddr, []byte(randomPayload(r)))
	require.NoError(r, err)
	r.WaitForTxReceiptOnEVM(tx)

	// wait for the cctxs to be mined
	cctxs := utils.WaitCctxsMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, 6, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctxs[0], "deposit eth")
	r.Logger.CCTX(*cctxs[1], "deposit eth 2")
	r.Logger.CCTX(*cctxs[2], "deposit and call eth")
	r.Logger.CCTX(*cctxs[3], "deposit and call eth 2")
	r.Logger.CCTX(*cctxs[4], "call")
	r.Logger.CCTX(*cctxs[5], "call 2")

	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[0].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[1].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[2].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[3].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[4].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[5].CctxStatus.Status)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.ETHZRC20, r.TestDAppV2ZEVMAddr, oldBalance, change, r.Logger)
}
