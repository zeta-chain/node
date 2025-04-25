package e2etests

import (
	"encoding/hex"
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiTokenWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package) and a token amount
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])

	// Given example contract on_call function arguments
	// only the CCTX's coinType (0x***::fake_usdc::FAKE_USDC) is needed, no additional arguments
	argumentTypes := []string{}
	objects := []string{
		r.SuiExample.GlobalConfigID.String(),
		r.SuiExample.PartnerID.String(),
		r.SuiExample.ClockID.String(),
	}

	// Given sui address
	// the example contract will just forward the withdrawn token to this address
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	suiAddress := signer.Address()

	// Given initial balance and called_count
	balanceBefore := r.SuiGetFungibleTokenBalance(suiAddress)
	calledCountBefore := r.SuiGetConnectedCalledCount()

	// create the payload message
	message, err := hex.DecodeString(suiAddress[2:]) // remove 0x prefix
	require.NoError(r, err)
	payload := sui.NewCallPayload(argumentTypes, objects, message)

	// ACT
	// approve both SUI gas budget token and fungible token ZRC20
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)
	r.ApproveFungibleTokenZRC20(r.GatewayZEVMAddr)

	// perform the fungible token withdraw and call
	tx := r.SuiWithdrawAndCallFungibleToken(
		targetPackageID,
		amount,
		payload,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the balance after the withdraw
	balanceAfter := r.SuiGetFungibleTokenBalance(signer.Address())
	require.EqualValues(r, balanceBefore+amount.Uint64(), balanceAfter)

	// verify the called_count increased by 1
	calledCountAfter := r.SuiGetConnectedCalledCount()
	require.Equal(r, calledCountBefore+1, calledCountAfter)
}
