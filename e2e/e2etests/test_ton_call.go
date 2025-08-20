package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// ARRANGE
	ctx := r.Ctx

	// Given a gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Given a zEVM contract
	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Example contract deployed at: %s", contractAddr.String())

	// Given a payload
	payload := randomPayloadBytes(r)

	// Given an approx `call` fee
	callFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpCall)
	require.NoError(r, err)

	// ACT
	// Perform TON tx
	cctx, err := r.TONCall(gw, sender, callFee, contractAddr, payload)

	// ASSERT
	require.NoError(r, err)
	r.Logger.CCTX(*cctx, "ton_call")

	// Ensure the example contract has been called, bar value should be set to amount
	utils.WaitAndVerifyExampleContractCallWithMsg(
		r,
		contract,
		big.NewInt(0),
		payload,
		[]byte(sender.GetAddress().ToRaw()),
	)
}
