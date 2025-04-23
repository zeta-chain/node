package e2etests

import (
	"encoding/hex"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiWithdrawAndCall(r *runner.E2ERunner, args []string) {
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

	// define a deterministic address and use it for on_call payload message
	// the example contract will just forward the withdrawn SUI token to this address
	suiAddress := "0x34a30aaee833d649d7313ddfe4ff5b6a9bac48803236b919369e6636fe93392e"
	message, err := hex.DecodeString(suiAddress[2:]) // remove 0x prefix
	require.NoError(r, err)
	balanceBefore := r.SuiGetSUIBalance(suiAddress)

	// create the payload
	payload := sui.NewCallPayload(argumentTypes, objects, message)

	// ACT
	// approve SUI ZRC20 token
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and call
	tx := r.SuiWithdrawAndCallSUI(targetPackageID, amount, payload)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// balance after
	balanceAfter := r.SuiGetSUIBalance(suiAddress)
	require.Equal(r, balanceBefore+amount.Uint64(), balanceAfter)
}
