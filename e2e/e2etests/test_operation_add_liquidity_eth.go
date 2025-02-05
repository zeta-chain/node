package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestOperationAddLiquidityETH is an operational test to add liquidity in gas token
func TestOperationAddLiquidityETH(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	liqZETA := utils.ParseBigInt(r, args[0])
	liqETH := utils.ParseBigInt(r, args[1])

	// perform the add liquidity
	r.AddLiquidityETH(liqZETA, liqETH)
}
