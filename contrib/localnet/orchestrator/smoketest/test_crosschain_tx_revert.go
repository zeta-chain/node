//go:build PRIVNET
// +build PRIVNET

package main

import (
	"fmt"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (sm *SmokeTest) TestCrosschainRevert() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Testing Crosschain swap zEVM contract call failed; user deposit should be refunded...\n")
	sm.zevmAuth.GasLimit = 10000000

	btcMinOutAmount := big.NewInt(0)
	msg := []byte{}
	for i := 0; i < 20-len(sm.ZEVMSwapAppAddr.Bytes()); i++ {
		msg = append(msg, 0)
	}
	msg = append(msg, sm.ZEVMSwapAppAddr.Bytes()...)
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
	cctx := WaitCctxMinedByInTxHash(txhash.Hex(), sm.cctxClient)
	fmt.Printf("cctx status: %s\n", cctx.CctxStatus.Status)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}

	newBal, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	diff := oldBal.Sub(oldBal, newBal)
	fmt.Printf("bal diff: %d\n", diff)

	LoudPrintf("Passed\n")
}
