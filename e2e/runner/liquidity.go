package runner

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

// AddLiquidityETH adds liquidity token to the uniswap pool ZETA/ETH
// we use the provided amount of ZETA and ETH to add liquidity as wanted amount
// 0 is used for the minimum amount of ZETA and ETH
func (r *E2ERunner) AddLiquidityETH(amountZETA, amountETH *big.Int) {
	// approve uni router
	r.ApproveETHZRC20(r.UniswapV2RouterAddr)

	previousValue := r.ZEVMAuth.Value
	r.ZEVMAuth.Value = amountZETA
	defer func() {
		r.ZEVMAuth.Value = previousValue
	}()

	r.Logger.Info("Adding liquidity to ZETA/ETH pool")
	tx, err := r.UniswapV2Router.AddLiquidityETH(
		r.ZEVMAuth,
		r.ETHZRC20Addr,
		amountETH,
		big.NewInt(0),
		big.NewInt(0),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, types.ReceiptStatusSuccessful, receipt.Status, "add liquidity failed")

	// get the pair address
	pairAddress, err := r.UniswapV2Factory.GetPair(&bind.CallOpts{}, r.WZetaAddr, r.ETHZRC20Addr)
	require.NoError(r, err)

	r.Logger.Info("ZETA/ETH pair address: %s", pairAddress.Hex())
}

// AddLiquidityERC20 adds liquidity token to the uniswap pool ZETA/ERC20
// we use the provided amount of ZETA and ERC20 to add liquidity as wanted amount
// 0 is used for the minimum amount of ZETA and ERC20
func (r *E2ERunner) AddLiquidityERC20(amountZETA, amountERC20 *big.Int) {
	// approve uni router
	r.ApproveERC20ZRC20(r.UniswapV2RouterAddr)

	previousValue := r.ZEVMAuth.Value
	r.ZEVMAuth.Value = amountZETA
	defer func() {
		r.ZEVMAuth.Value = previousValue
	}()

	r.Logger.Info("Adding liquidity to ZETA/ERC20 pool")
	tx, err := r.UniswapV2Router.AddLiquidityETH(
		r.ZEVMAuth,
		r.ERC20ZRC20Addr,
		amountERC20,
		big.NewInt(0),
		big.NewInt(0),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, types.ReceiptStatusSuccessful, receipt.Status, "add liquidity failed")

	// get the pair address
	pairAddress, err := r.UniswapV2Factory.GetPair(&bind.CallOpts{}, r.WZetaAddr, r.ERC20ZRC20Addr)
	require.NoError(r, err)

	r.Logger.Info("ZETA/ERC20 pair address: %s", pairAddress.Hex())
}
