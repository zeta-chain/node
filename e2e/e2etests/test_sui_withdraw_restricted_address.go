package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawRestrictedAddress tests that a withdrawal to a restricted address reverts to a revert address
// the test also add a case to check withdrawal to Sui invalid address immediately fail, we don't add a dedicated test as this is a small logic
func TestSuiWithdrawRestrictedAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given amount, receiver, revert address
	receiver := args[0]
	amount := utils.ParseBigInt(r, args[1])
	revertAddress := r.EVMAddress()

	// receiver balance before
	receiverBalanceBefore := r.SuiGetSUIBalance(receiver)

	// approve the ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// ACT
	// perform the withdraw to restricted receiver
	tx := r.SuiWithdraw(
		receiver,
		amount,
		r.SUIZRC20Addr,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(tx, "withdraw to restricted sui address")

	// wait for the withdraw tx to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// revert address balance before
	revertBalanceBefore, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// receiver balance should not change
	receiverBalanceAfter := r.SuiGetSUIBalance(receiver)
	require.EqualValues(r, receiverBalanceBefore, receiverBalanceAfter)

	// revert address should receive the amount
	revertBalanceAfter, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.EqualValues(r, new(big.Int).Add(revertBalanceBefore, amount), revertBalanceAfter)

	// Invalid address format
	tx, err = r.GatewayZEVM.Withdraw0(
		r.ZEVMAuth,
		[]byte("0x25db16c3ca555f6702c07860503107bb73cce9f6c1d6df00464529db15d5a5abaa"),
		amount,
		r.SUIZRC20Addr,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	if err != nil {
		// if error is not directly returned, check that the tx failed
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequiredTxFailed(r, receipt)
	}
}
