package utils

import (
	"context"
	"fmt"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

type CCTXClient = crosschaintypes.QueryClient

const (
	FungibleAdminName = "fungibleadmin"

	DefaultCctxTimeout = 4 * time.Minute
)

// WaitCctxMinedByInboundHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxMinedByInboundHash(
	ctx context.Context,
	inboundHash string,
	cctxClient crosschaintypes.QueryClient,
	logger infoLogger,
	cctxTimeout time.Duration,
) *crosschaintypes.CrossChainTx {
	cctxs := WaitCctxsMinedByInboundHash(ctx, inboundHash, cctxClient, 1, logger, cctxTimeout)
	if len(cctxs) == 0 {
		panic(fmt.Sprintf("cctx not found, inboundHash: %s", inboundHash))
	}
	return cctxs[len(cctxs)-1]
}

// WaitCctxsMinedByInboundHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxsMinedByInboundHash(
	ctx context.Context,
	inboundHash string,
	cctxClient crosschaintypes.QueryClient,
	cctxsCount int,
	logger infoLogger,
	cctxTimeout time.Duration,
) []*crosschaintypes.CrossChainTx {
	startTime := time.Now()

	timeout := DefaultCctxTimeout
	if cctxTimeout != 0 {
		timeout = cctxTimeout
	}

	// fetch cctxs by inboundHash
	for i := 0; ; i++ {
		// declare cctxs here so we can print the last fetched one if we reach timeout
		var cctxs []*crosschaintypes.CrossChainTx

		if time.Since(startTime) > timeout {
			cctxMessage := ""
			if len(cctxs) > 0 {
				cctxMessage = fmt.Sprintf(", last cctx: %v", cctxs[0].String())
			}

			panic(fmt.Sprintf("waiting cctx timeout, cctx not mined, inboundHash: %s%s", inboundHash, cctxMessage))
		}
		time.Sleep(1 * time.Second)

		// We use InTxHashToCctxData instead of InboundTrackerAllByChain to able to run these tests with the previous version
		// for the update tests
		// TODO: replace with InboundHashToCctxData once removed
		// https://github.com/zeta-chain/node/issues/2200
		res, err := cctxClient.InTxHashToCctxData(ctx, &crosschaintypes.QueryInboundHashToCctxDataRequest{
			InboundHash: inboundHash,
		})

		if err != nil {
			// prevent spamming logs
			if i%10 == 0 {
				logger.Info("Error getting cctx by inboundHash: %s", err.Error())
			}
			continue
		}
		if len(res.CrossChainTxs) < cctxsCount {
			// prevent spamming logs
			if i%10 == 0 {
				logger.Info(
					"not enough cctxs found by inboundHash: %s, expected: %d, found: %d",
					inboundHash,
					cctxsCount,
					len(res.CrossChainTxs),
				)
			}
			continue
		}
		cctxs = make([]*crosschaintypes.CrossChainTx, 0, len(res.CrossChainTxs))
		allFound := true
		for j, cctx := range res.CrossChainTxs {
			cctx := cctx
			if !IsTerminalStatus(cctx.CctxStatus.Status) {
				// prevent spamming logs
				if i%10 == 0 {
					logger.Info(
						"waiting for cctx index %d to be mined by inboundHash: %s, cctx status: %s, message: %s",
						j,
						inboundHash,
						cctx.CctxStatus.Status.String(),
						cctx.CctxStatus.StatusMessage,
					)
				}
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

// WaitCCTXMinedByIndex waits until cctx is mined; returns the cctxIndex
func WaitCCTXMinedByIndex(
	ctx context.Context,
	cctxIndex string,
	cctxClient crosschaintypes.QueryClient,
	logger infoLogger,
	cctxTimeout time.Duration,
) *crosschaintypes.CrossChainTx {
	startTime := time.Now()

	timeout := DefaultCctxTimeout
	if cctxTimeout != 0 {
		timeout = cctxTimeout
	}

	for i := 0; ; i++ {
		if time.Since(startTime) > timeout {
			panic(fmt.Sprintf(
				"waiting cctx timeout, cctx not mined, cctx: %s",
				cctxIndex,
			))
		}
		time.Sleep(1 * time.Second)

		// fetch cctx by index
		res, err := cctxClient.Cctx(ctx, &crosschaintypes.QueryGetCctxRequest{
			Index: cctxIndex,
		})
		if err != nil {
			// prevent spamming logs
			if i%10 == 0 {
				logger.Info("Error getting cctx by inboundHash: %s", err.Error())
			}
			continue
		}
		cctx := res.CrossChainTx
		if !IsTerminalStatus(cctx.CctxStatus.Status) {
			// prevent spamming logs
			if i%10 == 0 {
				logger.Info(
					"waiting for cctx to be mined from index: %s, cctx status: %s, message: %s",
					cctxIndex,
					cctx.CctxStatus.Status.String(),
					cctx.CctxStatus.StatusMessage,
				)
			}
			continue
		}

		return cctx
	}
}

type WaitOpts func(c *waitConfig)

// MatchStatus waits for a specific CCTX status.
func MatchStatus(s crosschaintypes.CctxStatus) WaitOpts {
	return Matches(func(tx crosschaintypes.CrossChainTx) bool {
		return tx.CctxStatus != nil && tx.CctxStatus.Status == s
	})
}

// Matches adds a filter to WaitCctxByInboundHash that checks cctxs match provided callback.
// ALL cctxs should match this filter.
func Matches(fn func(tx crosschaintypes.CrossChainTx) bool) WaitOpts {
	return func(c *waitConfig) { c.matchFunction = fn }
}

type waitConfig struct {
	matchFunction func(tx crosschaintypes.CrossChainTx) bool
}

// WaitCctxByInboundHash waits until cctx appears by inbound hash.
func WaitCctxByInboundHash(
	ctx context.Context,
	t require.TestingT,
	hash string,
	c CCTXClient,
	opts ...WaitOpts,
) []crosschaintypes.CrossChainTx {
	const tick = time.Millisecond * 200

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, DefaultCctxTimeout)
		defer cancel()
	}

	in := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: hash}

	var cfg waitConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	matches := func(txs []crosschaintypes.CrossChainTx) bool {
		if cfg.matchFunction == nil {
			return true
		}

		for _, tx := range txs {
			if ok := cfg.matchFunction(tx); !ok {
				return false
			}
		}

		return true
	}

	for {
		out, err := c.InTxHashToCctxData(ctx, in)
		statusCode, _ := status.FromError(err)

		switch {
		case statusCode.Code() == codes.NotFound:
			// expected; let's retry
		case err != nil:
			require.NoError(t, err, "failed to get cctx by inbound hash: %s", hash)
		case len(out.CrossChainTxs) > 0 && matches(out.CrossChainTxs):
			return out.CrossChainTxs
		case ctx.Err() == nil:
			require.NoError(t, err, "failed to get cctx by inbound hash (ctx error): %s", hash)
		}

		time.Sleep(tick)
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
	for i := 0; status.SyncInfo.LatestBlockHeight < height; i++ {
		status, err = rpc.Status(ctx)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)

		// prevent spamming logs
		if i%10 == 0 {
			logger.Info("waiting for block: %d, current height: %d\n", height, status.SyncInfo.LatestBlockHeight)
		}
	}
}
