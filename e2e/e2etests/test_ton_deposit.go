package e2etests

import (
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/runner/ton"
)

// TestTONDeposit (!) This boilerplate is a demonstration of E2E capabilities for TON integration
// Actual Deposit test is not implemented yet.
func TestTONDeposit(r *runner.E2ERunner, _ []string) {
	ctx, deployer := r.Ctx, r.TONDeployer

	// Given deployer
	deployerBalance, err := deployer.GetBalance(ctx)
	require.NoError(r, err, "failed to get deployer balance")
	require.NotZero(r, deployerBalance, "deployer balance is zero")

	// Given sample wallet with a balance of 50 TON
	sender, err := deployer.CreateWallet(ctx, ton.TONCoins(50))
	require.NoError(r, err)

	// That was funded (again) but the faucet
	_, err = deployer.Fund(ctx, sender.GetAddress(), ton.TONCoins(30))
	require.NoError(r, err)

	// Check sender balance
	sb, err := sender.GetBalance(ctx)
	require.NoError(r, err)

	senderBalance := math.NewUint(sb)

	// note that it's not exactly 80 TON, but 79.99... due to gas fees
	// We'll tackle gas math later.
	r.Logger.Print(
		"Balance of sender (%s): %s",
		sender.GetAddress().ToHuman(false, true),
		ton.FormatCoins(senderBalance),
	)
}
