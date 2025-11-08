package stream

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	cmtquery "github.com/cometbft/cometbft/libs/pubsub/query"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	cmttypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/evm/utils"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/types"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	streamSubscriberName = "evm-json-rpc"
	subscribBufferSize   = 1024

	headerStreamSegmentSize = 128
	headerStreamCapacity    = 128 * 32
	txStreamSegmentSize     = 1024
	txStreamCapacity        = 1024 * 32
	logStreamSegmentSize    = 2048
	logStreamCapacity       = 2048 * 32
)

var (
	evmEvents = cmtquery.MustCompile(fmt.Sprintf("%s='%s' AND %s.%s='%s'",
		cmttypes.EventTypeKey,
		cmttypes.EventTx,
		sdk.EventTypeMessage,
		sdk.AttributeKeyModule, evmtypes.ModuleName)).String()
	blockEvents          = cmttypes.QueryForEvent(cmttypes.EventNewBlock).String()
	evmTxHashKey         = fmt.Sprintf("%s.%s", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash)
	NewBlockHeaderEvents = cmtquery.MustCompile(fmt.Sprintf("%s='%s'", cmttypes.EventTypeKey, cmttypes.EventNewBlockHeader))
)

type RPCHeader struct {
	EthHeader *ethtypes.Header
	Hash      common.Hash
}

// RPCStream provides data streams for newHeads, logs, and pendingTransactions.
type RPCStream struct {
	evtClient rpcclient.EventsClient
	logger    log.Logger
	txDecoder sdk.TxDecoder

	// headerStream/logStream are backed by cometbft event subscription
	headerStream *Stream[RPCHeader]
	logStream    *Stream[*ethtypes.Log]

	// pendingTxStream is backed by check-tx ante handler
	pendingTxStream *Stream[common.Hash]

	wg sync.WaitGroup
}

func NewRPCStreams(evtClient rpcclient.EventsClient, logger log.Logger, txDecoder sdk.TxDecoder) *RPCStream {
	return &RPCStream{
		evtClient:       evtClient,
		logger:          logger,
		txDecoder:       txDecoder,
		pendingTxStream: NewStream[common.Hash](txStreamSegmentSize, txStreamCapacity),
	}
}

func (s *RPCStream) initSubscriptions() {
	if s.headerStream != nil {
		// already initialized
		return
	}

	s.headerStream = NewStream[RPCHeader](headerStreamSegmentSize, headerStreamCapacity)
	s.logStream = NewStream[*ethtypes.Log](logStreamSegmentSize, logStreamCapacity)

	ctx := context.Background()

	chBlocks, err := s.evtClient.Subscribe(ctx, streamSubscriberName, blockEvents, subscribBufferSize)
	if err != nil {
		panic(err)
	}

	chLogs, err := s.evtClient.Subscribe(ctx, streamSubscriberName, evmEvents, subscribBufferSize)
	if err != nil {
		if err := s.evtClient.UnsubscribeAll(context.Background(), streamSubscriberName); err != nil {
			s.logger.Error("failed to unsubscribe", "err", err)
		}
		panic(err)
	}

	go s.start(&s.wg, chBlocks, chLogs)
}

func (s *RPCStream) Close() error {
	if s.headerStream == nil {
		// not initialized
		return nil
	}

	if err := s.evtClient.UnsubscribeAll(context.Background(), streamSubscriberName); err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}

func (s *RPCStream) HeaderStream() *Stream[RPCHeader] {
	s.initSubscriptions()
	return s.headerStream
}

func (s *RPCStream) PendingTxStream() *Stream[common.Hash] {
	return s.pendingTxStream
}

func (s *RPCStream) LogStream() *Stream[*ethtypes.Log] {
	s.initSubscriptions()
	return s.logStream
}

// ListenPendingTx is a callback passed to application to listen for pending transactions in CheckTx.
func (s *RPCStream) ListenPendingTx(hash common.Hash) {
	s.PendingTxStream().Add(hash)
}

func (s *RPCStream) start(
	wg *sync.WaitGroup,
	chBlocks <-chan coretypes.ResultEvent,
	chLogs <-chan coretypes.ResultEvent,
) {
	wg.Add(1)
	defer func() {
		wg.Done()
		if err := s.evtClient.UnsubscribeAll(context.Background(), streamSubscriberName); err != nil {
			s.logger.Error("failed to unsubscribe", "err", err)
		}
	}()

	for {
		select {
		case ev, ok := <-chBlocks:
			if !ok {
				chBlocks = nil
				break
			}

			data, ok := ev.Data.(cmttypes.EventDataNewBlock)
			if !ok {
				s.logger.Error("event data type mismatch", "type", fmt.Sprintf("%T", ev.Data))
				continue
			}

			baseFee := types.BaseFeeFromEvents(data.ResultFinalizeBlock.Events)
			// TODO: After indexer improvement, we should get eth header event from indexer
			// Currently, many fields are missing or incorrect (e.g. bloom, receiptsRoot, ...)
			header := types.EthHeaderFromComet(data.Block.Header, ethtypes.Bloom{}, baseFee)
			s.headerStream.Add(RPCHeader{EthHeader: header, Hash: common.BytesToHash(data.BlockID.Hash)})

		case ev, ok := <-chLogs:
			if !ok {
				chLogs = nil
				break
			}

			if _, ok := ev.Events[evmTxHashKey]; !ok {
				// ignore transaction as it's not from the evm module
				continue
			}

			// get transaction result data
			dataTx, ok := ev.Data.(cmttypes.EventDataTx)
			if !ok {
				s.logger.Error("event data type mismatch", "type", fmt.Sprintf("%T", ev.Data))
				continue
			}
			height, err := utils.SafeUint64(dataTx.Height)
			if err != nil {
				continue
			}
			txLogs, err := evmtypes.DecodeTxLogs(dataTx.Result.Data, height)
			if err != nil {
				s.logger.Error("fail to decode evm tx response", "error", err.Error())
				continue
			}

			s.logStream.Add(txLogs...)
		}

		if chBlocks == nil && chLogs == nil {
			break
		}
	}
}
