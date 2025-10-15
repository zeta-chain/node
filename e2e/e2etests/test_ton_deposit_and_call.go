package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx depositAndCall fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDepositAndCall)
	require.NoError(r, err)

	// Given a senderWallet
	_, senderWallet, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)
	sender := []byte(senderWallet.GetAddress().String())

	// Given payload and a ZEVM contract
	contractAddr := r.TestDAppV2ZEVMAddr
	payload := randomPayload(r)
	r.AssertTestDAppZEVMCalled(false, payload, sender, big.NewInt(0))

	// ACT
	_, err = r.TONDepositAndCall(gw, senderWallet, amount, contractAddr, []byte(payload))

	// ASSERT
	require.NoError(r, err)

	expectedDeposit := amount.Sub(depositFee)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(expectedDeposit.BigInt())
	utils.WaitAndVerifyZRC20BalanceChange(r, r.TONZRC20, contractAddr, big.NewInt(0), change, r.Logger)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, sender, expectedDeposit.BigInt())
}
