package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	FungibleAdminName = "fungibleadmin"
)

// WaitCctxMinedByInTxHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxMinedByInTxHash(inTxHash string, cctxClient crosschaintypes.QueryClient) *crosschaintypes.CrossChainTx {
	cctxs := WaitCctxsMinedByInTxHash(inTxHash, cctxClient, 1)
	return cctxs[len(cctxs)-1]
}

// WaitCctxsMinedByInTxHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxsMinedByInTxHash(inTxHash string, cctxClient crosschaintypes.QueryClient, cctxsCount int) []*crosschaintypes.CrossChainTx {
	var cctxIndexes []string
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("Waiting for cctx to be mined by inTxHash: %s\n", inTxHash)
		res, err := cctxClient.InTxHashToCctx(context.Background(), &crosschaintypes.QueryGetInTxHashToCctxRequest{InTxHash: inTxHash})
		if err != nil {
			fmt.Println("Error getting cctx by inTxHash: ", err.Error())
			continue
		}
		if len(res.InTxHashToCctx.CctxIndex) < cctxsCount {
			fmt.Printf("Waiting for %d cctxs to be mined; %d cctxs are mined\n", cctxsCount, len(res.InTxHashToCctx.CctxIndex))
			continue
		}
		cctxIndexes = res.InTxHashToCctx.CctxIndex
		fmt.Printf("Deposit receipt cctx index: %v\n", cctxIndexes)
		break
	}
	var wg sync.WaitGroup
	var cctxs []*crosschaintypes.CrossChainTx
	for _, cctxIndex := range cctxIndexes {
		cctxIndex := cctxIndex
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				time.Sleep(3 * time.Second)
				res, err := cctxClient.Cctx(context.Background(), &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
				if err == nil && IsTerminalStatus(res.CrossChainTx.CctxStatus.Status) {
					fmt.Printf("Deposit receipt cctx status: %+v; The cctx is processed\n", res.CrossChainTx.CctxStatus.Status.String())
					cctxs = append(cctxs, res.CrossChainTx)
					break
				} else if err != nil {
					fmt.Println("Error getting cctx by index: ", err.Error())
				} else {
					cctxStatus := res.CrossChainTx.CctxStatus
					fmt.Printf(
						"Deposit receipt cctx status: %s; Message: %s; Waiting for the cctx to be processed\n",
						cctxStatus.Status.String(),
						cctxStatus.StatusMessage,
					)
				}
			}
		}()
	}
	wg.Wait()
	return cctxs
}

func IsTerminalStatus(status crosschaintypes.CctxStatus) bool {
	return status == crosschaintypes.CctxStatus_OutboundMined || status == crosschaintypes.CctxStatus_Aborted || status == crosschaintypes.CctxStatus_Reverted
}

// WaitForBlockHeight waits until the block height reaches the given height
func WaitForBlockHeight(height int64, rpcURL string) {
	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcURL, "/websocket")
	if err != nil {
		panic(err)
	}
	status := &coretypes.ResultStatus{}
	for status.SyncInfo.LatestBlockHeight < height {
		status, err = rpc.Status(context.Background())
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
		fmt.Printf("waiting for block: %d, current height: %d\n", height, status.SyncInfo.LatestBlockHeight)
	}
}
