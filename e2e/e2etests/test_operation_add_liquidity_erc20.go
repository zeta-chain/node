package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

// TestOperationAddLiquidityERC20 is an operational test to add liquidity in erc20 token
func TestOperationAddLiquidityERC20(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	liqZETA := big.NewInt(0)
	_, ok := liqZETA.SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestOperationAddLiquidityERC20")

	liqERC20 := big.NewInt(0)
	_, ok = liqERC20.SetString(args[1], 10)
	require.True(r, ok, "Invalid amount specified for TestOperationAddLiquidityERC20")

	// perform the add liquidity
	r.AddLiquidityERC20(liqZETA, liqERC20)
}
