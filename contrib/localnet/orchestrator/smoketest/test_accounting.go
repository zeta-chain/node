package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (sm *SmokeTest) TestAccounting() {
	LoudPrintf("Accouting: inventory check\n")

	{
		fmt.Printf("ZRC20 <-> Ether\n")
		tssBal, _ := sm.goerliClient.BalanceAt(context.Background(), TSSAddress, nil)
		zrc20Supply, _ := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
		fmt.Printf("Ether TSS Balance:  %d\n", tssBal)
		fmt.Printf("ZRC20 Total Supply: %d\n", zrc20Supply)
	}

	{
		fmt.Printf("ZRC20 <-> BTC\n")
		utxos, err := sm.btcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		var btcBalance float64
		for _, utxo := range utxos {
			btcBalance += utxo.Amount
		}
		zrc20Supply, _ := sm.BTCZRC20.TotalSupply(&bind.CallOpts{})
		fmt.Printf("BTC Balance:        %d\n", int64(btcBalance*1e8))
		fmt.Printf("ZRC20 Total Supply: %d\n", zrc20Supply)
	}
}
