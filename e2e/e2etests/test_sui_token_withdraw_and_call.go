package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiTokenWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given target package ID (example package), token amount and gas limit
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	// use the deployer address as on_call payload message
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	suiAddress := signer.Address()

	// Given initial balance and called_count
	balanceBefore := r.SuiGetFungibleTokenBalance(suiAddress)
	calledCountBefore := r.SuiGetConnectedCalledCount()

	// create the on_call payload
	authorizedSender := r.EVMAddress()
	payloadOnCall := r.SuiCreateExampleWACPayload(authorizedSender, suiAddress)
	message, err := payloadOnCall.PackABI()
	require.NoError(r, err)

	// ACT
	// approve both SUI gas budget token and fungible token ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)
	r.ApproveFungibleTokenZRC20(r.GatewayZEVMAddr)

	// perform the fungible token withdraw and authenticated call
	tx := r.SuiWithdrawAndCall(
		targetPackageID,
		amount,
		r.SuiTokenZRC20Addr,
		message,
		gasLimit,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	r.Logger.EVMTransaction(tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw_and_call")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the balance after the withdraw
	balanceAfter := r.SuiGetFungibleTokenBalance(signer.Address())
	require.EqualValues(r, balanceBefore+amount.Uint64(), balanceAfter)

	// verify the called_count increased by 1
	calledCountAfter := r.SuiGetConnectedCalledCount()
	require.Equal(r, calledCountBefore+1, calledCountAfter)
}
