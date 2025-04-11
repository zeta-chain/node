package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
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

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Check sender balance
	senderBalance, err := r.Clients.TON.GetBalanceOf(ctx, sender.GetAddress(), false)
	if err != nil {
		r.Logger.Print("Failed to get sender balance: %v", err)
		require.NoError(r, err)
	}

	r.Logger.Print("Sender balance: %s", toncontracts.FormatCoins(senderBalance))

	// Calculate total required amount (deposit + fee)
	totalRequired := amount.Add(depositFee)

	// Check if sender has enough balance
	if senderBalance.LT(totalRequired) {
		r.Logger.Print("⚠️ WARNING: Sender doesn't have enough TON to complete the deposit and call!")
		r.Logger.Print("Required: %s, Available: %s",
			toncontracts.FormatCoins(totalRequired),
			toncontracts.FormatCoins(senderBalance))
		r.Logger.Print("❓ This is expected when running without a faucet URL (ton_faucet: \"\")")
		r.Logger.Print("⏩ SKIPPING TEST: pre-conditions aren't met (insufficient balance).")
		return // Skip test instead of failing
	}

	// Given sample zEVM contract deployed by userTON account
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err, "unable to deploy example contract")
	r.Logger.Info("Example zevm contract deployed at: %s", contractAddr.String())

	// Given call data
	callData := []byte("hello from TON!")

	// ACT
	_, err = r.TONDepositAndCall(gw, sender, amount, contractAddr, callData)

	// ASSERT
	require.NoError(r, err)

	expectedDeposit := amount.Sub(depositFee)

	// check if example contract has been called, bar value should be set to amount
	utils.MustHaveCalledExampleContract(r, contract, expectedDeposit.BigInt(), []byte(sender.GetAddress().ToRaw()))

	// Check receiver's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)

	r.Logger.Info("Contract's zEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
