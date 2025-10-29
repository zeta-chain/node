package server

import (
	"context"
	"time"

	"github.com/cometbft/cometbft/libs/service"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"

	servertypes "github.com/cosmos/evm/server/types"
)

const (
	ServiceName = "EVMIndexerService"

	NewBlockWaitTimeout = 60 * time.Second
)

// EVMIndexerService indexes transactions for json-rpc service.
type EVMIndexerService struct {
	service.BaseService

	txIdxr servertypes.EVMTxIndexer
	client rpcclient.Client
}

// NewEVMIndexerService returns a new service instance.
func NewEVMIndexerService(
	txIdxr servertypes.EVMTxIndexer,
	client rpcclient.Client,
) *EVMIndexerService {
	is := &EVMIndexerService{txIdxr: txIdxr, client: client}
	is.BaseService = *service.NewBaseService(nil, ServiceName, is)
	return is
}

// OnStart implements service.Service by subscribing for new blocks
// and indexing them by events.
func (eis *EVMIndexerService) OnStart() error {
	ctx := context.Background()
	status, err := eis.client.Status(ctx)
	if err != nil {
		return err
	}
	latestBlock := status.SyncInfo.LatestBlockHeight
	newBlockSignal := make(chan struct{}, 1)

	// Use SubscribeUnbuffered here to ensure both subscriptions does not get
	// canceled due to not pulling messages fast enough. Cause this might
	// sometimes happen when there are no other subscribers.
	blockHeadersChan, err := eis.client.Subscribe(
		ctx,
		ServiceName,
		types.QueryForEvent(types.EventNewBlockHeader).String(),
		0)
	if err != nil {
		return err
	}

	go func() {
		for {
			msg := <-blockHeadersChan
			eventDataHeader := msg.Data.(types.EventDataNewBlockHeader)
			if eventDataHeader.Header.Height > latestBlock {
				latestBlock = eventDataHeader.Header.Height
				// notify
				select {
				case newBlockSignal <- struct{}{}:
				default:
				}
			}
		}
	}()

	lastBlock, err := eis.txIdxr.LastIndexedBlock()
	if err != nil {
		return err
	}
	if lastBlock == -1 {
		lastBlock = latestBlock
	}

	// blockErr indicates an error fetching an expected block or its results
	var blockErr error

	for {
		if latestBlock <= lastBlock || blockErr != nil {
			// two cases to enter this block:
			// 1. nothing to index (indexer is caught up). wait for signal of new block.
			// 2. previous attempt to index errored (failed to fetch the Block or BlockResults).
			//    in this case, wait before retrying the data fetching, rather than infinite looping
			//    a failing fetch. this can occur due to drive latency between the block existing and its
			//    block_results getting saved.
			select {
			case <-newBlockSignal:
			case <-time.After(NewBlockWaitTimeout):
			}
			continue
		}
		for i := lastBlock + 1; i <= latestBlock; i++ {
			var (
				block       *coretypes.ResultBlock
				blockResult *coretypes.ResultBlockResults
			)

			block, blockErr = eis.client.Block(ctx, &i)
			if blockErr != nil {
				eis.Logger.Error("failed to fetch block", "height", i, "err", blockErr)
				break
			}
			blockResult, blockErr = eis.client.BlockResults(ctx, &i)
			if blockErr != nil {
				eis.Logger.Error("failed to fetch block result", "height", i, "err", blockErr)
				break
			}
			if err := eis.txIdxr.IndexBlock(block.Block, blockResult.TxsResults); err != nil {
				eis.Logger.Error("failed to index block", "height", i, "err", err)
			}
			lastBlock = blockResult.Height
		}
	}
}
