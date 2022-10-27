package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/unlocker"

	"github.com/zeta-chain/zetacore/zetaclient/btc/model"
)

func main() {
	tx, err := getTx()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
}

func getTx() (string, error) {
	var txStr string

	BTCAddress := os.Getenv("BTC_ADDRESS")
	BTCWalletPK := os.Getenv("BTC_PK")
	if BTCAddress == "" || BTCWalletPK == "" {
		return txStr, fmt.Errorf("Empty BTCAddress or BTCWalletPK ")
	}
	tx := bt.NewTx()

	err := tx.From(
		"b7b0650a7c3a1bd4716369783876348b59f5404784970192cec1996e86950576",
		0,
		"76a9149cbe9f5e72fa286ac8a38052d1d5337aa363ea7f88ac",
		1000,
	)
	if err != nil {
		return txStr, fmt.Errorf("creating tx: %v\n", err.Error())
	}

	err = tx.PayToAddress(BTCAddress, 900)
	if err != nil {
		return txStr, fmt.Errorf("pay address : %v\n", err.Error())
	}

	evt := &model.ConnectorEvent{
		DestChainID: 1337,
		DestAddress: common.HexToAddress("0x1234"),
		Amount:      10,
	}

	err = tx.AddOpReturnOutput([]byte(evt.ToBTCOP()))
	if err != nil {
		return txStr, fmt.Errorf("op return : %v\n", err.Error())
	}

	decodedWif, err := wif.DecodeWIF(BTCWalletPK)
	if err != nil {
		return txStr, fmt.Errorf("invalid wallet pk: %v\n", err.Error())
	}

	err = tx.FillAllInputs(context.Background(), &unlocker.Getter{PrivateKey: decodedWif.PrivKey})
	if err != nil {
		return txStr, fmt.Errorf("fill inputs : %v\n", err.Error())
	}
	return fmt.Sprint(tx), nil
}
