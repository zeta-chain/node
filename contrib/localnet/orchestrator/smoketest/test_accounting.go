//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"io/ioutil"
	"math/big"
	"net/http"
)

func (sm *SmokeTest) CheckZRC20ReserveAndSupply() {
	{
		tssBal, _ := sm.goerliClient.BalanceAt(context.Background(), TSSAddress, nil)
		zrc20Supply, _ := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
		if tssBal.Cmp(zrc20Supply) < 0 {
			panic(fmt.Sprintf("ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssBal, zrc20Supply))
		} else {
			fmt.Printf("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d)\n", tssBal, zrc20Supply)
		}
	}

	{
		utxos, err := sm.btcRPCClient.ListUnspent()
		if err != nil {
			panic(err)
		}
		var btcBalance float64
		for _, utxo := range utxos {
			if utxo.Address == BTCTSSAddress.EncodeAddress() {
				btcBalance += utxo.Amount
			}
		}
		zrc20Supply, _ := sm.BTCZRC20.TotalSupply(&bind.CallOpts{})
		if int64(btcBalance*1e8) < zrc20Supply.Int64() {
			panic(fmt.Sprintf("BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d) ", int64(btcBalance*1e8), zrc20Supply))
		} else {
			fmt.Printf("BTC: Balance (%d) >= ZRC20 TotalSupply (%d)\n", int64(btcBalance*1e8), zrc20Supply)
		}
	}

	{
		usdtBal, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.ERC20CustodyAddr)
		if err != nil {
			panic(err)
		}
		zrc20Supply, _ := sm.USDTZRC20.TotalSupply(&bind.CallOpts{})
		if usdtBal.Cmp(zrc20Supply) < 0 {
			panic(fmt.Sprintf("USDT: TSS balance (%d) < ZRC20 TotalSupply (%d) ", usdtBal, zrc20Supply))
		} else {
			fmt.Printf("USDT: TSS balance (%d) >= ZRC20 TotalSupply (%d)\n", usdtBal, zrc20Supply)
		}
	}

	{
		type Amount struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		}

		type Response struct {
			Amount Amount `json:"amount"`
		}

		zetaLocked, err := sm.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		resp, err := http.Get("http://zetacore0:1317/cosmos/bank/v1beta1/supply/by_denom?denom=azeta")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var result Response
		err = json.Unmarshal(body, &result)
		if err != nil {
			panic(err)
		}
		zetaSupply, _ := big.NewInt(0).SetString(result.Amount.Amount, 10)
		if zetaLocked.Cmp(zetaSupply) < 0 {
			fmt.Printf(fmt.Sprintf("ZETA: TSS balance (%d) < ZRC20 TotalSupply (%d) ", zetaLocked, zetaSupply))
		} else {
			fmt.Printf("ZETA: TSS balance (%d) >= ZRC20 TotalSupply (%d)\n", zetaLocked, zetaSupply)
		}
	}
}
