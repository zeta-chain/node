package e2etests

import (
	"encoding/hex"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given a signer and an amount
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	amount := utils.ParseBigInt(r, args[0])

	// Given example contract in Sui network
	// sample withdrawAndCall payload
	// TODO: use real contract
	// https://github.com/zeta-chain/node/issues/3742
	argumentTypes := []string{
		r.SuiExample.TokenType,
	}
	objects := []string{
		r.SuiExample.GlobalConfigID,
		r.SuiExample.PoolID,
		r.SuiExample.PartnerID,
		r.SuiExample.ClockID,
	}

	// assemble payload with random Sui receiver address
	suiAddress := sample.SuiAddress(r)
	message, err := hex.DecodeString(suiAddress[2:]) // remove 0x prefix
	require.NoError(r, err)

	payload := sui.NewCallPayload(argumentTypes, objects, message)

	// ACT
	// approve SUI ZRC20 token
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// perform the withdraw and call
	tx := r.SuiWithdrawAndCallSUI(signer.Address(), amount, payload)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
