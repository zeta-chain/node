package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestOperationAddLiquidityERC20 is an operational test to add liquidity in erc20 token
func TestOperationAddLiquidityERC20(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	liqZETA := utils.ParseBigInt(r, args[0])
	liqERC20 := utils.ParseBigInt(r, args[1])

	// perform the add liquidity
	r.AddLiquidityERC20(liqZETA, liqERC20)
}
