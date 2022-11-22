package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

func ExploreKlaytnGetBlock() {
	endpoint := "https://klaytn-baobab-rpc.allthatnode.com:8551/Y8ZQrEReHIKMtMLfXneeLgQHrQc47jzn"
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(107513951))
	if err != nil {
		panic(err)
	}
	fmt.Printf("block number: %d; %d txs\n", block.Number().Int64(), block.Transactions().Len())

	errEmptyBlock := fmt.Errorf("server returned empty transaction list but block header indicates transactions")

	_, err = client.BlockByNumber(context.Background(), big.NewInt(107514095))
	if err != nil {
		fmt.Printf("error is errEmptyBlock: %v\n", errors.Is(err, errEmptyBlock))
		fmt.Printf("%v\n", strings.Contains(err.Error(), errEmptyBlock.Error()))
		fmt.Printf("error blockByNumber: %v\n", err)
	}
}
