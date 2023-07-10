//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (sm *SmokeTest) CheckZRC20ReserveAndSupply() {
	{
		tssBal, _ := sm.goerliClient.BalanceAt(context.Background(), TSSAddress, nil)
		zrc20Supply, _ := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
		if tssBal.Int64() < zrc20Supply.Int64() {
			panic(fmt.Sprintf("ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssBal, zrc20Supply))
		} else {
			fmt.Printf("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d) ", tssBal, zrc20Supply)
		}
	}

	{
		utxos, err := sm.btcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		var btcBalance float64
		for _, utxo := range utxos {
			btcBalance += utxo.Amount
		}
		zrc20Supply, _ := sm.BTCZRC20.TotalSupply(&bind.CallOpts{})
		if int64(btcBalance*1e8) < zrc20Supply.Int64() {
			panic(fmt.Sprintf("BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d) ", int64(btcBalance*1e8), zrc20Supply))
		} else {
			fmt.Printf("BTC: Balance (%d) >= ZRC20 TotalSupply (%d) ", int64(btcBalance*1e8), zrc20Supply)
		}
	}
}
