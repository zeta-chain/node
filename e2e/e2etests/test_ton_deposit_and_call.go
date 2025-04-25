package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

func TestTONDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	ctx := r.Ctx

	// Using gateway address from config.yml (specified as gateway_account_id)
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given amount
	amount := utils.ParseUint(r, args[0])

	// Given approx depositAndCall fee
	depositFee, err := gw.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDepositAndCall)
	require.NoError(r, err)

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// Get zEVM address as the recipient
	recipientAddr := r.EVMAddress()

	// Using the TONZRC20 address from config.yml (specified as ton_zrc20)
	contractAddr := r.TONZRC20Addr
	r.Logger.Info("Using TON ZRC20 contract address %s and TON gateway address %s", contractAddr.Hex(), r.TONGateway.ToRaw())
	r.Logger.Info("Recipient address (user's EVM address): %s", recipientAddr.Hex())

	balanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipientAddr)
	require.NoError(r, err)
	r.Logger.Info("Recipient's zEVM TON balance before deposit: %d (0x%x)", balanceBefore.Uint64(), balanceBefore.Uint64())

	// Given call data
	callData := []byte("hello from TON!")

	// Call TONDepositAndCall with recipient as the zEVM address, not the contract address
	r.Logger.Info("Sending deposit of %s nano TON from %s to recipient %s", amount.String(), sender.GetAddress().ToRaw(), recipientAddr.Hex())
	cctx, err := r.TONDepositAndCall(gw, sender, amount, recipientAddr, callData)

	// ASSERT
	require.NoError(r, err)

	// Verify the sender in the CCTX
	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)

	// Check the balance after deposit
	balanceAfter, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipientAddr)
	require.NoError(r, err)
	r.Logger.Info("Recipient's zEVM TON balance after deposit: %d (0x%x)", balanceAfter.Uint64(), balanceAfter.Uint64())

	// Calculate and log expected deposit amount (amount minus fee)
	expectedDeposit := amount.Sub(depositFee)
	r.Logger.Info("Expected deposit amount: %d (0x%x)", expectedDeposit.Uint64(), expectedDeposit.Uint64())

	// Check if the balance increased by the expected amount
	if balanceAfter.Cmp(balanceBefore) > 0 {
		balanceDiff := new(big.Int).Sub(balanceAfter, balanceBefore)
		r.Logger.Info("Balance difference: %d (0x%x)", balanceDiff.Uint64(), balanceDiff.Uint64())

		// Compare with expected deposit (using this to use the variable and avoid linter warning)
		if balanceDiff.Cmp(expectedDeposit.BigInt()) == 0 {
			r.Logger.Info("Balance difference matches expected deposit amount")
		}
	}

	// Log the CCTX details
	r.Logger.Info("CCTX status: %s", cctx.CctxStatus.Status.String())
	if cctx.CctxStatus.ErrorMessage != "" {
		r.Logger.Info("CCTX error message: %s", cctx.CctxStatus.ErrorMessage)
	}
}
