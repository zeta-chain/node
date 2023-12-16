package smoketests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestZRC20Swap(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	utils.LoudPrintf("Swap ZRC20 USDT for ZRC20 ETH\n")

	// TODO: move into setup and skip it if already initialized
	// https://github.com/zeta-chain/node-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	tx, err := sm.UniswapV2Factory.CreatePair(sm.ZevmAuth, sm.USDTZRC20Addr, sm.ETHZRC20Addr)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	//fmt.Printf("USDT-ETH pair receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	usdtEthPair, err := sm.UniswapV2Factory.GetPair(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.ETHZRC20Addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT-ETH pair receipt pair addr %s\n", usdtEthPair.Hex())

	tx, err = sm.USDTZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("USDT ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	tx, err = sm.ETHZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("ETH ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	// temporarily increase gas limit to 400000
	previousGasLimit := sm.ZevmAuth.GasLimit
	defer func() {
		sm.ZevmAuth.GasLimit = previousGasLimit
	}()

	sm.ZevmAuth.GasLimit = 400000
	tx, err = sm.UniswapV2Router.AddLiquidity(
		sm.ZevmAuth,
		sm.USDTZRC20Addr,
		sm.ETHZRC20Addr,
		big.NewInt(90000),
		big.NewInt(1000),
		big.NewInt(90000),
		big.NewInt(1000),
		sm.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("Add liquidity receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	balETHBefore, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	ethOutAmout := big.NewInt(1)
	tx, err = sm.UniswapV2Router.SwapExactTokensForTokens(
		sm.ZevmAuth,
		big.NewInt(1000),
		ethOutAmout,
		[]ethcommon.Address{sm.USDTZRC20Addr, sm.ETHZRC20Addr},
		sm.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("Swap USDT for ETH ZRC20 %s status %d\n", receipt.TxHash, receipt.Status)

	balETHAfter, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	ethDiff := big.NewInt(0).Sub(balETHAfter, balETHBefore)
	if ethDiff.Cmp(ethOutAmout) < 0 {
		panic("swap failed")
	}
}
