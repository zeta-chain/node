//go:build PRIVNET
// +build PRIVNET

package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (sm *SmokeTest) TestCrosschainRevert() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Testing Crosschain swap failed tx revert...\n")
	sm.zevmAuth.GasLimit = 10000000

	btcMinOutAmount := big.NewInt(0)
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

	oldBal, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	// Should deposit USDT for swap, swap for BTC and withdraw BTC
	txhash := sm.DepositERC20(big.NewInt(8e7), msg)
	_ = WaitCctxMinedByInTxHash(txhash.Hex(), sm.cctxClient)

	start := time.Now()
	for {
		if time.Since(start) > 30*time.Second {
			panic("waiting tx balance revert timeout")
		}
		newBal, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := oldBal.Sub(oldBal, newBal)
		if diff.Cmp(big.NewInt(1e5)) < 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	LoudPrintf("Passed\n")
}
