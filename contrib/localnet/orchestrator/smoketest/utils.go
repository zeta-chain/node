package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// wait until cctx is mined; returns the cctxIndex
func WaitCctxMinedByInTxHash(inTxHash string, cctxClient types.QueryClient) *types.CrossChainTx {
	var cctxIndex string
	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.InTxHashToCctx(context.Background(), &types.QueryGetInTxHashToCctxRequest{InTxHash: inTxHash})
		if err != nil {
			continue
		}
		cctxIndex = res.InTxHashToCctx.CctxIndex
		fmt.Printf("Deposit receipt cctx index: %s\n", cctxIndex)
		break
	}
	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: cctxIndex})
		if err != nil || res.CrossChainTx.CctxStatus.Status != types.CctxStatus_OutboundMined {
			fmt.Printf("Deposit receipt cctx status: %s\n", res.CrossChainTx.CctxStatus.Status.String())
			continue
		}
		fmt.Printf("Deposit receipt cctx status: %+v; success\n", res.CrossChainTx.CctxStatus.Status.String())
		return res.CrossChainTx
	}

}

func LoudPrintf(format string, a ...any) {
	fmt.Println("=======================================")
	fmt.Printf(format, a...)
	fmt.Println("=======================================")
}

func CheckNonce(client *ethclient.Client, addr ethcommon.Address, expectedNonce uint64) error {
	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return err
	}
	if nonce != expectedNonce {
		return fmt.Errorf("want nonce %d; got %d", expectedNonce, nonce)
	}
	return nil
}
