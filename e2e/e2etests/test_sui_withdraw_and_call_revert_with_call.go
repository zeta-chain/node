package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSuiWithdrawAndCallRevertWithCall executes withdrawAndCall on zevm gateway with SUI token.
// The outbound authenticated call is rejected by the connected module due to unauthorized sender address,
// and the 'onRevert' method is called in the ZEVM to handle the revert.
func TestSuiWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given target package ID (example package), SUI amount and gas limit
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	// use the deployer address as on_call payload message
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	suiAddress := signer.Address()

	// given receiver and TSS balances in Sui network
	receiverBalanceBefore := r.SuiGetSUIBalance(suiAddress)
	tssBalanceBefore := r.SuiGetSUIBalance(r.SuiTSSAddress)

	// given ZEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2ZEVMAddr
	dAppBalanceBefore, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// create a 'on_call' payload that gives authorization to a random address,
	// such that the 'r.EVMAddress()' becomes an unauthorized sender in the 'on_call'
	authorizedSender := sample.EthAddress()
	payloadOnCall := r.SuiCreateExampleWACPayload(authorizedSender, suiAddress)
	message, err := payloadOnCall.PackABI()
	require.NoError(r, err)

	// given random payload for 'onRevert'
	payloadOnRevert := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payloadOnRevert, amount)

	// ACT
	// approve SUI ZRC20 token
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and authenticated call with revert options
	tx := r.SuiWithdrawAndCall(
		targetPackageID,
		amount,
		r.SUIZRC20Addr,
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

	// the receiver balance in Sui network should remain the same
	receiverBalanceAfter := r.SuiGetSUIBalance(suiAddress)
	require.Equal(r, receiverBalanceBefore, receiverBalanceAfter)

	// the TSS balance in Sui network should be higher or equal to the balance before
	// reason is that the max budget is refunded to the TSS
	tssBalanceAfter := r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore)

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
	dAppBalanceAfter, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())
}
