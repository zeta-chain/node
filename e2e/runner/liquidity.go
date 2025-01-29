package runner

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

// AddLiquidityETH adds liquidity token to the uniswap pool ZETA/ETH
func (r *E2ERunner) AddLiquidityETH(amountZETA, amountETH *big.Int) {
	r.ApproveETHZRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.ETHZRC20Addr, amountZETA, amountETH)
}

// AddLiquidityERC20 adds liquidity token to the uniswap pool ZETA/ERC20
func (r *E2ERunner) AddLiquidityERC20(amountZETA, amountERC20 *big.Int) {
	r.ApproveERC20ZRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.ERC20ZRC20Addr, amountZETA, amountERC20)
}

// AddLiquiditySPL adds liquidity token to the uniswap pool ZETA/SPL
func (r *E2ERunner) AddLiquiditySPL(amountZETA, amountSPL *big.Int) {
	r.ApproveSPLZRC20(r.UniswapV2RouterAddr)
	r.addLiquidity(r.SPLZRC20Addr, amountZETA, amountSPL)
}

// addLiquidity adds liquidity token to the uniswap pool ZETA/token
// we use the provided amount of ZETA and token to add liquidity as wanted amount
// 0 is used for the minimum amount of ZETA and token
func (r *E2ERunner) addLiquidity(tokenAddr ethcommon.Address, amountZETA, amountToken *big.Int) {
	previousValue := r.ZEVMAuth.Value
	r.ZEVMAuth.Value = amountZETA
	defer func() {
		r.ZEVMAuth.Value = previousValue
	}()

	r.Logger.Info("Adding liquidity to ZETA/token pool")
	tx, err := r.UniswapV2Router.AddLiquidityETH(
		r.ZEVMAuth,
		tokenAddr,
		amountToken,
		big.NewInt(0),
		big.NewInt(0),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == types.ReceiptStatusFailed {
		r.Logger.Error("Add liquidity failed for ZETA/token")
	}

	// get the pair address
	pairAddress, err := r.UniswapV2Factory.GetPair(&bind.CallOpts{}, r.WZetaAddr, tokenAddr)
	require.NoError(r, err)

	r.Logger.Info("ZETA/token pair address: %s", pairAddress.Hex())
}
