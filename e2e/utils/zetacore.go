package utils

import (
	"context"
	"fmt"
	"math/big"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/constant"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

type CCTXClient = crosschaintypes.QueryClient

const (
	EmergencyPolicyName       = "emergency"
	AdminPolicyName           = "admin"
	OperationalPolicyName     = "operational"
	UserEmissionsWithdrawName = "emissions_withdraw"

	// The timeout was increased from 4 to 6 , which allows for a higher success in test runs
	// However this needs to be researched as to why the increase in timeout was needed.
	// https://github.com/zeta-chain/node/issues/2690

	DefaultCctxTimeout = 8 * time.Minute

	// nodeSyncTolerance is the time tolerance for the ZetaChain nodes behind a RPC to be synced
	nodeSyncTolerance = constant.ZetaBlockTime * 5
)

// GetCCTXByInboundHash gets cctx by inbound hash
func GetCCTXByInboundHash(
	ctx context.Context,
	client crosschaintypes.QueryClient,
	inboundHash string,
) []crosschaintypes.CrossChainTx {
	t := TestingFromContext(ctx)

	// query cctx by inbound hash
	in := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: inboundHash}
	res, err := client.InTxHashToCctxData(ctx, in)

	require.NoError(t, err)

	return res.CrossChainTxs
}

// WaitCctxMinedByInboundHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxMinedByInboundHash(
	ctx context.Context,
	inboundHash string,
	client crosschaintypes.QueryClient,
	logger infoLogger,
	timeout time.Duration,
) *crosschaintypes.CrossChainTx {
	t := TestingFromContext(ctx)
	cctxs := WaitCctxsMinedByInboundHash(ctx, inboundHash, client, 1, logger, timeout)
	require.NotEmpty(t, cctxs, "cctx not found, inboundHash: %s", inboundHash)

	return cctxs[len(cctxs)-1]
}

// EnsureNoCctxMinedByInboundHash ensures no cctx is mined by inbound hash
func EnsureNoCctxMinedByInboundHash(
	ctx context.Context,
	inboundHash string,
	client crosschaintypes.QueryClient,
) {
	t := TestingFromContext(ctx)

	// query cctx by inbound hash
	in := &crosschaintypes.QueryGetInboundHashToCctxRequest{InboundHash: inboundHash}
	_, err := client.InboundHashToCctx(ctx, in)
	require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))
}

// WaitCctxsMinedByInboundHash waits until cctx is mined; returns the cctxIndex (the last one)
func WaitCctxsMinedByInboundHash(
	ctx context.Context,
	inboundHash string,
	client crosschaintypes.QueryClient,
	cctxsCount int,
	logger infoLogger,
	timeout time.Duration,
) []*crosschaintypes.CrossChainTx {
	if timeout == 0 {
		timeout = DefaultCctxTimeout
	}

	t := TestingFromContext(ctx)

	startTime := time.Now()
	in := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: inboundHash}

	for i := 0; ; i++ {
		// declare cctxs here so we can print the last fetched one if we reach timeout
		var cctxs []*crosschaintypes.CrossChainTx

		elapsed := time.Since(startTime)
		timedOut := time.Since(startTime) > timeout
		require.False(t, timedOut, "waiting cctx timeout, cctx not mined, inbound hash: %s, elapsed: %s",
			inboundHash, elapsed)

		require.NoError(t, ctx.Err())

		time.Sleep(500 * time.Millisecond)

		// We use InTxHashToCctxData instead of InboundTrackerAllByChain to able to run these tests with the previous version
		// for the update tests
		// TODO: replace with InboundHashToCctxData once removed
		// https://github.com/zeta-chain/node/issues/2200
		res, err := client.InTxHashToCctxData(ctx, in)
		if err != nil {
			// prevent spamming logs
			if i%20 == 0 {
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
			if !cctx.CctxStatus.Status.IsTerminal() {
				// prevent spamming logs
				if i%20 == 0 {
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
	client crosschaintypes.QueryClient,
	logger infoLogger,
	timeout time.Duration,
) *crosschaintypes.CrossChainTx {
	if timeout == 0 {
		timeout = DefaultCctxTimeout
	}

	t := TestingFromContext(ctx)
	startTime := time.Now()
	in := &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex}

	for i := 0; ; i++ {
		require.False(t, time.Since(startTime) > timeout, "waiting cctx timeout, cctx not mined, cctx: %s", cctxIndex)
		require.NoError(t, ctx.Err())

		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		// fetch cctx by index
		res, err := client.Cctx(ctx, in)
		if err != nil {
			// prevent spamming logs
			if i%10 == 0 {
				logger.Info("Error getting cctx by inboundHash: %s", err.Error())
			}
			continue
		}

		cctx := res.CrossChainTx
		if !cctx.CctxStatus.Status.IsTerminal() {
			// prevent spamming logs
			if i%20 == 0 {
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

// WaitOutboundTracker wait for outbound tracker to be filled with 'hashCount' hashes
func WaitOutboundTracker(
	ctx context.Context,
	client crosschaintypes.QueryClient,
	chainID int64,
	nonce uint64,
	hashCount int,
	logger infoLogger,
	timeout time.Duration,
) []string {
	if timeout == 0 {
		timeout = DefaultCctxTimeout
	}

	t := TestingFromContext(ctx)
	startTime := time.Now()
	in := &crosschaintypes.QueryAllOutboundTrackerByChainRequest{Chain: chainID}

	for {
		require.False(
			t,
			time.Since(startTime) > timeout,
			fmt.Sprintf("waiting outbound tracker timeout, chainID: %d, nonce: %d", chainID, nonce),
		)

		// wait for a Zeta block before querying outbound tracker
		time.Sleep(constant.ZetaBlockTime)

		outboundTracker, err := client.OutboundTrackerAllByChain(ctx, in)
		require.NoError(t, err)

		// loop through all outbound trackers
		for i, tracker := range outboundTracker.OutboundTracker {
			if tracker.Nonce == nonce {
				logger.Info("Tracker[%d]:\n", i)
				logger.Info("  ChainId: %d\n", tracker.ChainId)
				logger.Info("  Nonce: %d\n", tracker.Nonce)
				logger.Info("  HashList:\n")

				hashes := []string{}
				for j, hash := range tracker.HashList {
					hashes = append(hashes, hash.TxHash)
					logger.Info("    hash[%d]: %s\n", j, hash.TxHash)
				}
				if len(hashes) >= hashCount {
					return hashes
				}
			}
		}
	}
}

type WaitOpts func(c *waitConfig)

// MatchStatus is the WaitOpts that matches CCTX with the given status.
func MatchStatus(s crosschaintypes.CctxStatus) WaitOpts {
	return Matches(func(tx crosschaintypes.CrossChainTx) bool {
		return tx.CctxStatus != nil && tx.CctxStatus.Status == s
	})
}

// MatchReverted is the WaitOpts that matches reverted CCTX.
func MatchReverted() WaitOpts {
	return Matches(func(tx crosschaintypes.CrossChainTx) bool {
		return tx.GetCctxStatus().Status == crosschaintypes.CctxStatus_Reverted &&
			len(tx.OutboundParams) == 2 &&
			tx.OutboundParams[1].Hash != ""
	})
}

// HasOutboundTxHash returns true when the CCTX has been assigned an outbound hash.
// This now happens when the first tracker is written.
func HasOutboundTxHash() WaitOpts {
	return Matches(func(tx crosschaintypes.CrossChainTx) bool {
		return len(tx.OutboundParams) > 0 &&
			tx.OutboundParams[0].Hash != ""
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

// WaitCctxRevertedByInboundHash waits until cctx is reverted by inbound hash.
func WaitCctxRevertedByInboundHash(
	ctx context.Context,
	t require.TestingT,
	hash string,
	c CCTXClient,
) crosschaintypes.CrossChainTx {
	// wait for cctx to be reverted
	cctx := WaitCctxByInboundHash(ctx, t, hash, c, MatchReverted())

	return cctx
}

// WaitCctxAbortedByInboundHash waits until cctx is aborted by inbound hash.
func WaitCctxAbortedByInboundHash(
	ctx context.Context,
	t require.TestingT,
	hash string,
	c CCTXClient,
) crosschaintypes.CrossChainTx {
	// wait for cctx to be aborted
	cctx := WaitCctxByInboundHash(ctx, t, hash, c, MatchStatus(crosschaintypes.CctxStatus_Aborted))

	return cctx
}

// WaitCctxByInboundHash waits until cctx appears by inbound hash.
func WaitCctxByInboundHash(
	ctx context.Context,
	t require.TestingT,
	hash string,
	c CCTXClient,
	opts ...WaitOpts,
) crosschaintypes.CrossChainTx {
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
			return out.CrossChainTxs[0]
		case ctx.Err() == nil:
			require.NoError(t, err, "failed to get cctx by inbound hash (ctx error): %s", hash)
		}

		time.Sleep(tick)
	}
}

// WaitForBlockHeight waits until the block height reaches the given height
func WaitForBlockHeight(
	ctx context.Context,
	desiredHeight int64,
	rpcURL string,
	logger infoLogger,
) error {
	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcURL, "/websocket")
	if err != nil {
		return errors.Wrap(err, "unable to create rpc client")
	}

	var currentHeight int64
	for i := 0; currentHeight < desiredHeight; i++ {
		time.Sleep(1 * time.Second)
		s, err := rpc.Status(ctx)
		if err != nil {
			continue
		}
		currentHeight = s.SyncInfo.LatestBlockHeight

		// prevent spamming logs
		if i%10 == 0 {
			logger.Info("waiting for block: %d, current height: %d\n", desiredHeight, currentHeight)
		}
		if i > 30 {
			return errors.Wrapf(err, "unable to get status after %d attempts (current height: %d)", i, currentHeight)
		}
	}

	return nil
}

// WaitForZetaBlocks waits for the given number of Zeta blocks
func WaitForZetaBlocks(
	ctx context.Context,
	t require.TestingT,
	zevmClient *ethclient.Client,
	waitBlocks uint64,
	timeout time.Duration,
) {
	oldHeight, err := zevmClient.BlockNumber(ctx)
	require.NoError(t, err)

	// wait for given number of Zeta blocks
	newHeight := oldHeight
	startTime := time.Now()
	checkInterval := constant.ZetaBlockTime / 2
	for {
		time.Sleep(checkInterval)
		require.False(
			t,
			time.Since(startTime) > timeout,
			"timeout waiting for Zeta blocks, oldHeight: %d, currentHeight: %d, waitBlocks: %d, elapsed: %v",
			oldHeight,
			newHeight,
			waitBlocks,
			time.Since(startTime),
		)

		// check how many blocks elapsed
		newHeight, err = zevmClient.BlockNumber(ctx)
		require.NoError(t, err)
		if newHeight >= oldHeight+waitBlocks {
			return
		}
	}
}

// WaitAndVerifyZRC20BalanceChange waits for the zrc20 balance of the given address to change by the given delta amount
// This function is to tolerate the fact that the balance update may not be synced across all nodes behind a RPC.
func WaitAndVerifyZRC20BalanceChange(
	t require.TestingT,
	zrc20 *zrc20.ZRC20,
	address common.Address,
	oldBalance *big.Int,
	change BalanceChange,
	logger infoLogger,
) {
	// wait until the expected balance is reached or timeout
	startTime := time.Now()
	checkInterval := 2 * time.Second
	for {
		time.Sleep(checkInterval)
		require.False(t, time.Since(startTime) > nodeSyncTolerance, "timeout waiting for balance change")

		symbol, err := zrc20.Symbol(&bind.CallOpts{})
		if err != nil {
			logger.Info("unable to get symbol: %s", err.Error())
			continue
		}

		newBalance, err := zrc20.BalanceOf(&bind.CallOpts{}, address)
		if err != nil {
			logger.Info("unable to get balance: %s", err.Error())
			continue
		}

		if oldBalance.Cmp(newBalance) == 0 {
			logger.Info("balance has not changed yet")
			continue
		}
		logger.Info("%s balance changed from %d to %d on address %s", symbol, oldBalance, newBalance, address.Hex())

		change.Verify(t, oldBalance, newBalance)

		return
	}
}
