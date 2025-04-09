package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawRestrictedAddress tests a withdrawal to a restricted address that reverts to a revert address
func TestSuiWithdrawRestrictedAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseBigInt(r, args[0])

	// ARRANGE
	// Given receiver, revert address
	receiver := sample.RestrictedSuiAddressTest
	revertAddress := sample.EthAddress()

	// balances before
	receiverBalanceBefore := r.SuiGetSUIBalance(receiver)
	revertBalanceBefore, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// approve the ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// ACT
	// perform the withdraw to restricted receiver
	tx := r.SuiWithdrawSUI(
		receiver,
		amount,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(*tx, "withdraw to restricted sui address")

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
}
