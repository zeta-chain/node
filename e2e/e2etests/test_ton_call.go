package e2etests

import (
	"fmt"
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONToZEVMCall(r *runner.E2ERunner, args []string) {
	fmt.Println("TestTONToZEVMCall")

	require.Len(r, args, 0)

	// ARRANGE
	ctx := r.Ctx

	// Given a gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Given payload and a ZEVM contract
	contractAddr := r.TestDAppV2ZEVMAddr
	payload := randomPayload(r)
	r.AssertTestDAppZEVMCalled(false, payload, big.NewInt(0))

	// Given an approx `call` fee
	callFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpCall)
	require.NoError(r, err)

	// ACT
	// Perform TON tx
	cctx, err := r.TONCall(gw, sender, callFee, contractAddr, []byte(payload))

	// ASSERT
	require.NoError(r, err)
	r.Logger.CCTX(*cctx, "ton_call")

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, big.NewInt(0))
}
