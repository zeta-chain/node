package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestOperationAddLiquiditySPL is an operational test to add liquidity in spl token
func TestOperationAddLiquiditySPL(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	liqZETA := utils.ParseBigInt(r, args[0])
	liqSPL := utils.ParseBigInt(r, args[1])

	// perform the add liquidity
	r.AddLiquiditySPL(liqZETA, liqSPL)
}
