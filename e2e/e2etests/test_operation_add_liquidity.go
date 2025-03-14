package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func parseArgsForAddLiquidity(r *runner.E2ERunner, args []string) (*big.Int, *big.Int) {
	require.Len(r, args, 2)

	liqZETA := utils.ParseBigInt(r, args[0])
	liqToken := utils.ParseBigInt(r, args[1])

	return liqZETA, liqToken
}

// TestOperationAddLiquidityETH is an operational test to add liquidity in the ZETA/ETH pool (evm gas token)
func TestOperationAddLiquidityETH(r *runner.E2ERunner, args []string) {
	liqZETA, liqETH := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityETH(liqZETA, liqETH)
}

// TestOperationAddLiquidityERC20 is an operational test to add liquidity in the ZETA/ERC20 pool
func TestOperationAddLiquidityERC20(r *runner.E2ERunner, args []string) {
	liqZETA, liqERC20 := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityERC20(liqZETA, liqERC20)
}

// TestOperationAddLiquidityBTC is an operational test to add liquidity in the ZETA/BTC pool
func TestOperationAddLiquidityBTC(r *runner.E2ERunner, args []string) {
	liqZETA, liqBTC := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityBTC(liqZETA, liqBTC)
}

// TestOperationAddLiquiditySOL is an operational test to add liquidity in the ZETA/SOL pool
func TestOperationAddLiquiditySOL(r *runner.E2ERunner, args []string) {
	liqZETA, liqSOL := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySOL(liqZETA, liqSOL)
}

// TestOperationAddLiquiditySPL is an operational test to add liquidity in the ZETA/SPL pool
func TestOperationAddLiquiditySPL(r *runner.E2ERunner, args []string) {
	liqZETA, liqSPL := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySPL(liqZETA, liqSPL)
}

// TestOperationAddLiquiditySUI is an operational test to add liquidity in the ZETA/SUI pool
func TestOperationAddLiquiditySUI(r *runner.E2ERunner, args []string) {
	liqZETA, liqSUI := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySUI(liqZETA, liqSUI)
}

// TestOperationAddLiquiditySuiFungibleToken is an operational test to add liquidity in the ZETA/SuiFungibleToken pool
func TestOperationAddLiquiditySuiFungibleToken(r *runner.E2ERunner, args []string) {
	liqZETA, liqSuiToken := parseArgsForAddLiquidity(r, args)
	r.AddLiquiditySuiFungibleToken(liqZETA, liqSuiToken)
}

// TestOperationAddLiquidityTON is an operational test to add liquidity in the ZETA/TON pool
func TestOperationAddLiquidityTON(r *runner.E2ERunner, args []string) {
	liqZETA, liqTON := parseArgsForAddLiquidity(r, args)
	r.AddLiquidityTON(liqZETA, liqTON)
}
