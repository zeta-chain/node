package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

// TestOperationAddLiquidityETH is an operational test to add liquidity in gas token
func TestOperationAddLiquidityETH(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	liqZETA := big.NewInt(0)
	_, ok := liqZETA.SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestOperationAddLiquidityETH")

	liqETH := big.NewInt(0)
	_, ok = liqETH.SetString(args[1], 10)
	require.True(r, ok, "Invalid amount specified for TestOperationAddLiquidityETH")

	// perform the add liquidity
	r.AddLiquidityETH(liqZETA, liqETH)
}
