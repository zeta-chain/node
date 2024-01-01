package utils

import (
	"context"
	"fmt"
	"time"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	FungibleAdminName = "fungibleadmin"

	CctxTimeout = 3 * time.Minute
)

// WaitCctxMinedByInTxHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxMinedByInTxHash(
	ctx context.Context,
	inTxHash string,
	cctxClient crosschaintypes.QueryClient,
	logger infoLogger,
) *crosschaintypes.CrossChainTx {
	cctxs := WaitCctxsMinedByInTxHash(ctx, inTxHash, cctxClient, 1, logger)
	if len(cctxs) == 0 {
		panic(fmt.Sprintf("cctx not found, inTxHash: %s", inTxHash))
	}
	return cctxs[len(cctxs)-1]
}

// WaitCctxsMinedByInTxHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxsMinedByInTxHash(
	ctx context.Context,
	inTxHash string,
	cctxClient crosschaintypes.QueryClient,
	cctxsCount int,
	logger infoLogger,
) []*crosschaintypes.CrossChainTx {
	startTime := time.Now()

	// fetch cctxs by inTxHash
	for {
		if time.Since(startTime) > CctxTimeout {
			panic(fmt.Sprintf("waiting cctx timeout, cctx not mined, inTxHash: %s", inTxHash))
		}

		time.Sleep(1 * time.Second)
		res, err := cctxClient.InTxHashToCctxData(ctx, &crosschaintypes.QueryInTxHashToCctxDataRequest{
			InTxHash: inTxHash,
		})
		if err != nil {
			logger.Info("Error getting cctx by inTxHash: %s", err.Error())
			continue
		}
		if len(res.CrossChainTxs) < cctxsCount {
			logger.Info(
				"not enough cctxs found by inTxHash: %s, expected: %d, found: %d",
				inTxHash,
				cctxsCount,
				len(res.CrossChainTxs),
			)
			continue
		}
		cctxs := make([]*crosschaintypes.CrossChainTx, 0, len(res.CrossChainTxs))
		allFound := true
		for i, cctx := range res.CrossChainTxs {
			cctx := cctx
			if !IsTerminalStatus(cctx.CctxStatus.Status) {
				logger.Info(
					"waiting for cctx index %d to be mined by inTxHash: %s, cctx status: %s, message: %s",
					i,
					inTxHash,
					cctx.CctxStatus.Status.String(),
					cctx.CctxStatus.StatusMessage,
				)
				allFound = false
				break
			}
			cctxs = append(cctxs, &cctx)
		}
		if !allFound {
			continue
		}
		return cctxs
	}
}

func IsTerminalStatus(status crosschaintypes.CctxStatus) bool {
	return status == crosschaintypes.CctxStatus_OutboundMined ||
		status == crosschaintypes.CctxStatus_Aborted ||
		status == crosschaintypes.CctxStatus_Reverted
}

// WaitForBlockHeight waits until the block height reaches the given height
func WaitForBlockHeight(
	ctx context.Context,
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
		status, err = rpc.Status(ctx)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
		logger.Info("waiting for block: %d, current height: %d\n", height, status.SyncInfo.LatestBlockHeight)
	}
}
