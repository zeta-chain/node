package main

import (
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

var (
	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
)

func main() {
	ethclient, err := ethclient.Dial("http://eth:8545")
	if err != nil {
		panic(err)
	}
	bn, err := ethclient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	chainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ChainID: %d, Current block number: %d\n", chainID, bn)
	bal, err := ethclient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Ether\n", DeployerAddress.Hex(), bal.Div(bal, big.NewInt(1e18)))

	// ==================== Deploying contracts ====================
	fmt.Printf("Step 1: Deploying a contract\n")

	// ==================== Interacting with contracts ====================

}
