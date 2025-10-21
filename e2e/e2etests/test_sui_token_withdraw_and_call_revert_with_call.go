package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiTokenWithdrawAndCallRevertWithCall executes withdrawAndCall on zevm gateway with fungible token.
// The outbound is rejected by the connected module due to the special payload message "revert" and the
// 'onRevert' method is called in the ZEVM to handle the revert.
func TestSuiTokenWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given target package ID (example package), token amount and gas limit
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	// create the special revert payload for 'on_call'
	revertPayloadOnCall, err := r.SuiCreateExampleWACPayloadForRevert()
	require.NoError(r, err)
	message, err := revertPayloadOnCall.PackABI()
	require.NoError(r, err)

	// given ZEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2ZEVMAddr
	dAppBalanceBefore, err := r.SuiTokenZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// given random payload for 'onRevert'
	payloadOnRevert := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payloadOnRevert, amount)

	// ACT
	// approve both SUI gas budget token and fungible token ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)
	r.ApproveFungibleTokenZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and authenticated call with revert options
	tx := r.SuiWithdrawAndCall(
		targetPackageID,
		amount,
		r.SuiTokenZRC20Addr,
		message,
		gasLimit,
		gatewayzevm.RevertOptions{
			CallOnRevert:     true,
			RevertAddress:    dAppAddress,
			RevertMessage:    []byte(payloadOnRevert),
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// should have called 'onRevert'
	r.AssertTestDAppZEVMCalled(true, payloadOnRevert, nil, big.NewInt(0))

	// sender and message should match
	sender, err := r.TestDAppV2ZEVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payloadOnRevert),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, sender)

	// the dApp address should get reverted amount
	dAppBalanceAfter, err := r.SuiTokenZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())
}
