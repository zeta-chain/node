package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
)

func (sm *SmokeTest) TestCrosschainSwap() {
	LoudPrintf("Testing Bitcoin ERC20 crosschain swap...\n")
	// Firstly, deposit 1.15 BTC into Zeta for liquidity
	sm.DepositBTC()
	// Secondly, deposit 1000.0 USDT into Zeta for liquidity
	sm.DepositERC20(big.NewInt(1e9), []byte{})

	sm.zevmAuth.GasLimit = 20000000
	tx, err := sm.UniswapV2Factory.CreatePair(sm.zevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
	usdtBtcPair, err := sm.UniswapV2Factory.GetPair(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT-BTC pair receipt txhash %s status %d pair addr %s\n", receipt.TxHash, receipt.Status, usdtBtcPair.Hex())

	tx, err = sm.USDTZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("USDT ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	tx, err = sm.BTCZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("BTC ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	// Add 100 USDT liq and 0.001 BTC
	bal, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on BTC ZRC20: %d\n", bal)
	bal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	tx, err = sm.UniswapV2Router.AddLiquidity(sm.zevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr, big.NewInt(1e8), big.NewInt(1e5), big.NewInt(1e8), big.NewInt(1e5), DeployerAddress, big.NewInt(time.Now().Add(10*time.Minute).Unix()))
	if err != nil {
		fmt.Printf("Error liq %s", err.Error())
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	fmt.Printf("Add liquidity receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	usdtBtc, err := contracts.NewUniswapV2Pair(usdtBtcPair, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	res, err := usdtBtc.GetReserves(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reserves %s %s\n", res.Reserve0, res.Reserve1)

	btcMinOutAmount := big.NewInt(1e3)
	msg := []byte{}
	for i := 0; i < 20-len(HexToAddress(ZEVMSwapAppAddr).Bytes()); i++ {
		msg = append(msg, 0)
	}
	msg = append(msg, HexToAddress(ZEVMSwapAppAddr).Bytes()...)
	for i := 0; i < 32-len(sm.BTCZRC20Addr.Bytes()); i++ {
		msg = append(msg, 0)
	}
	msg = append(msg, sm.BTCZRC20Addr.Bytes()...)
	for i := 0; i < 32-len(DeployerAddress.Bytes()); i++ {
		msg = append(msg, 0)
	}
	msg = append(msg, DeployerAddress.Bytes()...)
	for i := 0; i < 32-len(btcMinOutAmount.Bytes()); i++ {
		msg = append(msg, 0)
	}
	msg = append(msg, btcMinOutAmount.Bytes()...)
	// Should deposit USDT for swap, swap for BTC and withdraw BTC
	fmt.Printf("gas limit %d\n", sm.zevmAuth.GasLimit)
	sm.DepositERC20(big.NewInt(11e6), msg)
}
