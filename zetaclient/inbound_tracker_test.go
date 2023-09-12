package zetaclient

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestEVMChainClient_CheckReceiptForCoinTypeGas(t *testing.T) {
	goerliClient, err := ethclient.Dial("http://eth:8545")
	if err != nil {
		panic(err)
	}
}
