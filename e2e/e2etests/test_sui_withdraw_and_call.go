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

	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// sample withdrawAndCall payload
	// TODO: use real contract
	// https://github.com/zeta-chain/node/issues/3742
	argumentTypes := []string{
		"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
	}
	objects := []string{
		"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
		"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
		"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3",
		"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
	}
	message, err := hex.DecodeString("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06")
	require.NoError(r, err)

	payload := sui.NewCallPayload(argumentTypes, objects, message)

	// perform the withdraw and call
	tx := r.SuiWithdrawAndCallSUI(signer.Address(), amount, payload)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
