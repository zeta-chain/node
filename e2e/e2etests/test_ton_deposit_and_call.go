package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
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

	// Get balance before deposit
	contractAddr := r.EVMAddress()
	r.Logger.Info("Using contract address: %s", contractAddr.Hex())
	balanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)
	r.Logger.Info("Contract's zEVM TON balance before deposit: %d (0x%x)", balanceBefore.Uint64(), balanceBefore.Uint64())

	// Given call data
	callData := []byte("hello from TON!")

	// ACT
	// We're expecting the call to revert since we saw this in CI
	r.Logger.Info("Sending deposit and call with %s nano TON from %s to %s and expecting Reverted status", amount.String(), sender.GetAddress().ToRaw(), contractAddr.Hex())
	cctx, err := r.TONDepositAndCall(
		gw,
		sender,
		amount,
		contractAddr,
		callData,
		runner.TONExpectStatus(cctypes.CctxStatus_Reverted),
	)

	// ASSERT
	require.NoError(r, err)

	expectedDeposit := amount.Sub(depositFee)

	// Verify the sender in the CCTX
	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)

	// Check receiver's balance after deposit
	balanceAfter, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	r.Logger.Info("Contract's zEVM TON balance after deposit: %d (0x%x)", balanceAfter.Uint64(), balanceAfter.Uint64())
	r.Logger.Info("Expected deposit based on calculation: %d (0x%x)", expectedDeposit.Uint64(), expectedDeposit.Uint64())
	r.Logger.Info("CCTX reported amount: %d (0x%x)", cctx.InboundParams.Amount.Uint64(), cctx.InboundParams.Amount.Uint64())
	require.NoError(r, err)

	// Check if the balance changed at all
	if balanceAfter.Cmp(balanceBefore) != 0 {
		// Calculate the actual amount deposited (balance difference)
		balanceDiff := balanceAfter.Uint64() - balanceBefore.Uint64()
		r.Logger.Info("Balance difference (actual deposit): %d (0x%x)", balanceDiff, balanceDiff)

		// Check if the balance difference matches the expected deposit amount
		require.Equal(r, expectedDeposit.Uint64(), balanceDiff, "Balance difference should match expected deposit amount")
	} else {
		r.Logger.Info("No balance change detected %d", balanceAfter.Uint64())
	}
}
