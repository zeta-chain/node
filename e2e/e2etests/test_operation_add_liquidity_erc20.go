package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
)

// TestOperationAddLiquidityETH is an operational test to add liquidity in gas token
func TestOperationAddLiquidityETH(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// #nosec G115 e2e - always in range
	liqZETA := big.NewInt(int64(parseInt(r, args[0])))
	// #nosec G115 e2e - always in range
	liqETH := big.NewInt(int64(parseInt(r, args[1])))

	// perform the add liquidity
	r.AddLiquidityETH(liqZETA, liqETH)
}
