package main

import (
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

func DebugGoerliOutboundStuckInMempool() {
	endpoint := "https://nd-411-320-015.p2pify.com/375d5d4e0ce5ab6d4fb1c4d24498febc"
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		panic(err)
	}

	txhash := ethcommon.HexToHash("0xcee3b1492ad3a14d93303192c1e27f79b2ffd27782d8c20b1d73109cf6073fc8")
	tx, pending, err := client.TransactionByHash(context.Background(), txhash)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" pending: %v\n", pending)
	if !pending {
		return
	}

	endpoint2 := "https://eth-goerli.g.alchemy.com/v2/J-W7M8JtqtQI3ckka76fz9kxX-Sa_CSK"
	client2, err := ethclient.Dial(endpoint2)
	if err != nil {
		panic(err)
	}
	err = client2.SendTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("successfully broadcasted tx %s\n", tx.Hash().Hex())

	time.Sleep(30 * time.Second)
	receipt, err := client2.TransactionReceipt(context.Background(), txhash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("tx receipt mined in block %d\n", receipt.BlockNumber)
}
