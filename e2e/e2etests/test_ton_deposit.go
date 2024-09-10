package e2etests

import (
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/utils"

	"github.com/zeta-chain/node/e2e/runner"
)

func TestTONDeposit(r *runner.E2ERunner, _ []string) {
	ctx := r.Ctx

	deployerBalance, err := r.TONDeployer.GetBalance(ctx)
	require.NoError(r, err, "failed to get deployer balance")

	r.Logger.Print("TON deployer address %s", r.TONDeployer.Wallet().GetAddress().ToHuman(false, true))

	require.False(r, deployerBalance.IsZero(), "deployer balance is zero")

	r.Logger.Print("TON deployer balance: %s", prettyPrintTON(deployerBalance))
}

func prettyPrintTON(v math.Uint) string {
	return utils.HumanFriendlyCoinsRepr(int64(v.Uint64()))
}
