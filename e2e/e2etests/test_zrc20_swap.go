package e2etests

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestZRC20Swap(r *runner.E2ERunner, _ []string) {
	// TODO: move into setup and skip it if already initialized
	// https://github.com/zeta-chain/node-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	tx, err := r.UniswapV2Factory.CreatePair(r.ZEVMAuth, r.ERC20ZRC20Addr, r.ETHZRC20Addr)
	if err != nil {
		r.Logger.Print("ℹ️ create pair error")
	} else {
		utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	}

	zrc20EthPair, err := r.UniswapV2Factory.GetPair(&bind.CallOpts{}, r.ERC20ZRC20Addr, r.ETHZRC20Addr)
	require.NoError(r, err)

	r.Logger.Info("ZRC20-ETH pair receipt pair addr %s", zrc20EthPair.Hex())

	tx, err = r.ERC20ZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("ERC20 ZRC20 approval receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	tx, err = r.ETHZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("ETH ZRC20 approval receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	// temporarily increase gas limit to 400000
	previousGasLimit := r.ZEVMAuth.GasLimit
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	r.ZEVMAuth.GasLimit = 400000
	tx, err = r.UniswapV2Router.AddLiquidity(
		r.ZEVMAuth,
		r.ERC20ZRC20Addr,
		r.ETHZRC20Addr,
		big.NewInt(90000),
		big.NewInt(1000),
		big.NewInt(90000),
		big.NewInt(1000),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Add liquidity receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	balETHBefore, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	ethOutAmout := big.NewInt(1)
	tx, err = r.UniswapV2Router.SwapExactTokensForTokens(
		r.ZEVMAuth,
		big.NewInt(1000),
		ethOutAmout,
		[]ethcommon.Address{r.ERC20ZRC20Addr, r.ETHZRC20Addr},
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Swap ERC20 ZRC20 for ETH ZRC20 %s status %d", receipt.TxHash, receipt.Status)

	balETHAfter, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	ethDiff := big.NewInt(0).Sub(balETHAfter, balETHBefore)
	require.NotEqual(r, -1, ethDiff.Cmp(ethOutAmout), "swap failed")
}
