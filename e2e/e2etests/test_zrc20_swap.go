package e2etests

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestZRC20Swap(r *runner.E2ERunner) {
	// TODO: move into setup and skip it if already initialized
	// https://github.com/zeta-chain/node-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	tx, err := r.UniswapV2Factory.CreatePair(r.ZevmAuth, r.USDTZRC20Addr, r.ETHZRC20Addr)
	if err != nil {
		r.Logger.Print("ℹ️create pair error")
	} else {
		utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	}

	usdtEthPair, err := r.UniswapV2Factory.GetPair(&bind.CallOpts{}, r.USDTZRC20Addr, r.ETHZRC20Addr)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("USDT-ETH pair receipt pair addr %s", usdtEthPair.Hex())

	tx, err = r.USDTZRC20.Approve(r.ZevmAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("USDT ZRC20 approval receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	tx, err = r.ETHZRC20.Approve(r.ZevmAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("ETH ZRC20 approval receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	// temporarily increase gas limit to 400000
	previousGasLimit := r.ZevmAuth.GasLimit
	defer func() {
		r.ZevmAuth.GasLimit = previousGasLimit
	}()

	r.ZevmAuth.GasLimit = 400000
	tx, err = r.UniswapV2Router.AddLiquidity(
		r.ZevmAuth,
		r.USDTZRC20Addr,
		r.ETHZRC20Addr,
		big.NewInt(90000),
		big.NewInt(1000),
		big.NewInt(90000),
		big.NewInt(1000),
		r.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Add liquidity receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	balETHBefore, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	ethOutAmout := big.NewInt(1)
	tx, err = r.UniswapV2Router.SwapExactTokensForTokens(
		r.ZevmAuth,
		big.NewInt(1000),
		ethOutAmout,
		[]ethcommon.Address{r.USDTZRC20Addr, r.ETHZRC20Addr},
		r.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Swap USDT for ETH ZRC20 %s status %d", receipt.TxHash, receipt.Status)

	balETHAfter, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	ethDiff := big.NewInt(0).Sub(balETHAfter, balETHBefore)
	if ethDiff.Cmp(ethOutAmout) < 0 {
		panic("swap failed")
	}
}
