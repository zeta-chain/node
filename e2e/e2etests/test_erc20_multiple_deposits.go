package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20MultipleDeposits(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.Logger.Info("starting erc20 multiple deposits test")

	oldBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	// send erc20 tokens to test dapp to be deposited to gateway
	tx := r.SendERC20OnEVM(r.TestDAppV2EVMAddr, new(big.Int).Div(amount, big.NewInt(1e18)).Int64())
	r.WaitForTxReceiptOnEVM(tx)

	// set value of the payable transactions
	previousValue := r.EVMAuth.Value
	fee, err := r.GatewayEVM.AdditionalActionFeeWei(nil)
	require.NoError(r, err)

	// add 1 fee to provided amount to pay for 2 inbounds (1st one is free)
	r.EVMAuth.Value = new(big.Int).Add(amount, fee)

	// send multiple deposit through contract
	r.Logger.Print("🏃test multiple erc20 deposits through contract")
	tx, err = r.TestDAppV2EVM.GatewayMultipleERC20Deposits(
		r.EVMAuth,
		r.TestDAppV2ZEVMAddr,
		r.ERC20Addr,
		amount,
		[]byte(randomPayload(r)),
	)
	require.NoError(r, err)
	r.WaitForTxReceiptOnEVM(tx)
	r.Logger.Print("🍾 multiple erc20 deposits through contract observed")

	r.EVMAuth.Value = previousValue

	// wait for the cctxs to be mined
	cctxs := utils.WaitCctxsMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, 2, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctxs[0], "deposit erc20")
	r.Logger.CCTX(*cctxs[1], "deposit erc20 and call")

	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[0].CctxStatus.Status)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxs[1].CctxStatus.Status)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.ERC20ZRC20, r.TestDAppV2ZEVMAddr, oldBalance, change, r.Logger)
}
