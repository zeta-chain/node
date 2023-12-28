package utils

import (
	"context"
	"sync"
	"time"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	FungibleAdminName = "fungibleadmin"

	CctxTimeout = 120 * time.Second
)

// WaitCctxMinedByInTxHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxMinedByInTxHash(
	inTxHash string,
	cctxClient crosschaintypes.QueryClient,
	logger infoLogger,
) *crosschaintypes.CrossChainTx {
	cctxs := WaitCctxsMinedByInTxHash(inTxHash, cctxClient, 1, logger)
	return cctxs[len(cctxs)-1]
}

// WaitCctxsMinedByInTxHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxsMinedByInTxHash(
	inTxHash string,
	cctxClient crosschaintypes.QueryClient,
	cctxsCount int,
	logger infoLogger,
) []*crosschaintypes.CrossChainTx {
	// start a go routine that will send a signal to the channel when the timeout is reached
	// this is to prevent the go routine from leaking
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(CctxTimeout)
		timeout <- true
	}()

	// wait for cctx to be mined
	var cctxIndexes []string
	for {
		// check if timeout is reached
		select {
		case <-timeout:
			panic("waiting cctx timeout")
		default:
		}

		time.Sleep(1 * time.Second)
		logger.Info("Waiting for cctx to be mined by inTxHash: %s", inTxHash)
		res, err := cctxClient.InTxHashToCctx(
			context.Background(),
			&crosschaintypes.QueryGetInTxHashToCctxRequest{InTxHash: inTxHash},
		)
		if err != nil {
			logger.Info("Error getting cctx by inTxHash: %s", err.Error())
			continue
		}
		if len(res.InTxHashToCctx.CctxIndex) < cctxsCount {
			logger.Info(
				"Waiting for %d cctxs to be mined; %d cctxs are mined",
				cctxsCount,
				len(res.InTxHashToCctx.CctxIndex),
			)
			continue
		}
		cctxIndexes = res.InTxHashToCctx.CctxIndex
		logger.Info("Deposit receipt cctx index: %v", cctxIndexes)
		break
	}

	// cctxs have been mined, retrieve all data
	var wg sync.WaitGroup
	var cctxs []*crosschaintypes.CrossChainTx
	var cctxsMutex sync.Mutex

	for _, cctxIndex := range cctxIndexes {
		cctxIndex := cctxIndex
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				time.Sleep(1 * time.Second)
				res, err := cctxClient.Cctx(context.Background(), &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
				if err == nil && IsTerminalStatus(res.CrossChainTx.CctxStatus.Status) {
					logger.Info("Deposit receipt cctx status: %+v; The cctx is processed", res.CrossChainTx.CctxStatus.Status.String())
					cctxsMutex.Lock()
					cctxs = append(cctxs, res.CrossChainTx)
					cctxsMutex.Unlock()
					break
				} else if err != nil {
					logger.Info("Error getting cctx by index: ", err.Error())
				} else {
					cctxStatus := res.CrossChainTx.CctxStatus
					logger.Info(
						"Deposit receipt cctx status: %s; Message: %s; Waiting for the cctx to be processed",
						cctxStatus.Status.String(),
						cctxStatus.StatusMessage,
					)
				}
			}
		}()
	}

	// go routine to wait for all go routines to finish
	allMined := make(chan bool, 1)
	go func() {
		wg.Wait()
		allMined <- true
	}()

	// wait for all cctxs to be mined
	select {
	case <-allMined:
		logger.Info("All cctxs are mined")
	case <-timeout:
		panic("waiting cctx timeout")
	}

	return cctxs
}

func IsTerminalStatus(status crosschaintypes.CctxStatus) bool {
	return status == crosschaintypes.CctxStatus_OutboundMined ||
		status == crosschaintypes.CctxStatus_Aborted ||
		status == crosschaintypes.CctxStatus_Reverted
}

// WaitForBlockHeight waits until the block height reaches the given height
func WaitForBlockHeight(
	height int64,
	rpcURL string,
	logger infoLogger,
) {
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
		time.Sleep(1 * time.Second)
		logger.Info("waiting for block: %d, current height: %d\n", height, status.SyncInfo.LatestBlockHeight)
	}
}
