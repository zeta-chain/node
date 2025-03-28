package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawAndCallRevertWithCall executes withdrawAndCall on zevm and calls a smart contract on Sui.
// The execution is rejected in Sui smart contract 'on_call' function, and 'nonce_increase' is called instead.
//
// Note: this test is faked as we don't have the support for cross-chain call support for Sui yet
// and it uses simple Gas coin withdrawal as a workaround to test the failed (rejected) outbound.
func TestSuiWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseBigInt(r, args[0])

	// ARRANGE
	// given signer
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	signerBalanceBefore := r.SuiGetSUIBalance(signer.Address())

	// given ZEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2ZEVMAddr
	dAppBalanceBefore, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// given random payload
	payload := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payload, amount)

	// ACT
	// approve the ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and call
	tx := r.SuiWithdrawSUI(signer.Address(), amount)
	r.Logger.EVMTransaction(*tx, "withdraw")

	// ASSERT
	// wait for the CCTX to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// should have called 'onRevert'
	r.AssertTestDAppZEVMCalled(true, payload, big.NewInt(0))

	// sender and message should match
	sender, err := r.TestDAppV2ZEVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payload),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, sender)

	// signer balance should remain unchanged in Sui chain
	signerBalanceAfter := r.SuiGetSUIBalance(signer.Address())
	require.Equal(r, signerBalanceBefore, signerBalanceAfter)

	// the dApp address should get reverted amount
	dAppBalanceAfter, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())
}
