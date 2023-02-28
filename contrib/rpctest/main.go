package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <blocknum>\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("Start testing the zEVM ETH JSON-RPC for all txs...\n")
	fmt.Printf("Test1: simple gas voter tx\n")

	bn, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	zevmClient, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	block, err := zevmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Block number: %d, num of txs %d (should be 1)\n", block.Number(), len(block.Transactions()))

}
