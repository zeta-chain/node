package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx deposit fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDeposit)
	require.NoError(r, err)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Given zEVM sender address
	recipient := r.EVMAddress()
	r.Logger.Info("Recipient address: %s", recipient)

	// Get balance before deposit
	balanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	require.NoError(r, err)
	r.Logger.Info("Recipient's zEVM TON balance before deposit: %d (0x%x)", balanceBefore.Uint64(), balanceBefore.Uint64())

	// ACT
	cctx, err := r.TONDeposit(gw, sender, amount, recipient)

	// ASSERT
	require.NoError(r, err)

	// Check CCTX
	expectedDeposit := amount.Sub(depositFee)

	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)
	require.Equal(r, expectedDeposit.Uint64(), cctx.InboundParams.Amount.Uint64())

	// Check receiver's balance after deposit
	balanceAfter, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	r.Logger.Info("Recipient's zEVM TON balance after deposit: %d (0x%x)", balanceAfter.Uint64(), balanceAfter.Uint64())
	r.Logger.Info("Expected deposit based on calculation: %d (0x%x)", expectedDeposit.Uint64(), expectedDeposit.Uint64())
	r.Logger.Info("CCTX reported amount: %d (0x%x)", cctx.InboundParams.Amount.Uint64(), cctx.InboundParams.Amount.Uint64())
	require.NoError(r, err)

	// Calculate the actual amount deposited (balance difference)
	balanceDiff := balanceAfter.Uint64() - balanceBefore.Uint64()
	r.Logger.Info("Balance difference (actual deposit): %d (0x%x)", balanceDiff, balanceDiff)

	// Check if the balance difference matches the expected deposit amount
	require.Equal(r, expectedDeposit.Uint64(), balanceDiff, "Balance difference should match expected deposit amount")
}
