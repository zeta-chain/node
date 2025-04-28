package e2etests

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawAndCallRevertWithCall executes withdraw and call on zevm gateway.
// The outbound is rejected by the connected module due to invalid payload (invalid address),
// and 'onRevert' is called instead to handle the revert.
func TestSuiWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package) and a SUI amount
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])

	// Given example contract on_call function arguments
	argumentTypes := []string{
		r.SuiExample.TokenType.String(),
	}
	objects := []string{
		r.SuiExample.GlobalConfigID.String(),
		r.SuiExample.PoolID.String(),
		r.SuiExample.PartnerID.String(),
		r.SuiExample.ClockID.String(),
	}

	// define an invalid address to cause 'on_call' failure
	invalidAddress := "8f569597ebca884b784d32678a6f"
	message, err := hex.DecodeString(invalidAddress)
	require.NoError(r, err)

	// given ZEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2ZEVMAddr
	dAppBalanceBefore, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// given random payload for 'onRevert'
	payloadOnRevert := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payloadOnRevert, amount)

	// create the payload for 'on_call'
	payloadOnCall := sui.NewCallPayload(argumentTypes, objects, message)

	// ACT
	// approve SUI ZRC20 token
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and call with revert options
	tx := r.SuiWithdrawAndCallSUI(
		targetPackageID,
		amount,
		payloadOnCall,
		gatewayzevm.RevertOptions{
			CallOnRevert:     true,
			RevertAddress:    dAppAddress,
			RevertMessage:    []byte(payloadOnRevert),
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// should have called 'onRevert'
	r.AssertTestDAppZEVMCalled(true, payloadOnRevert, big.NewInt(0))

	// sender and message should match
	sender, err := r.TestDAppV2ZEVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payloadOnRevert),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, sender)

	// the dApp address should get reverted amount
	dAppBalanceAfter, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())
}
