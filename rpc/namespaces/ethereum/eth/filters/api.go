package filters

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/zeta-chain/node/rpc/stream"
	"github.com/zeta-chain/node/rpc/types"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/client"
)

var (
	errInvalidBlockRange      = errors.New("invalid block range params")
	errPendingLogsUnsupported = errors.New("pending logs are not supported")
)

// FilterAPI gathers
type FilterAPI interface {
	NewPendingTransactionFilter() rpc.ID
	NewBlockFilter() rpc.ID
	NewFilter(criteria filters.FilterCriteria) (rpc.ID, error)
	GetFilterChanges(id rpc.ID) (interface{}, error)
	GetFilterLogs(ctx context.Context, id rpc.ID) ([]*ethtypes.Log, error)
	UninstallFilter(id rpc.ID) bool
	GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*ethtypes.Log, error)
}

// Backend defines the methods requided by the PublicFilterAPI backend
type Backend interface {
	GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (map[string]interface{}, error)
	HeaderByNumber(blockNum types.BlockNumber) (*ethtypes.Header, error)
	HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error)
	CometBlockByHash(hash common.Hash) (*coretypes.ResultBlock, error)
	CometBlockResultByNumber(height *int64) (*coretypes.ResultBlockResults, error)
	GetLogs(blockHash common.Hash) ([][]*ethtypes.Log, error)
	GetLogsByHeight(*int64) ([][]*ethtypes.Log, error)
	BlockBloomFromCometBlock(blockRes *coretypes.ResultBlockResults) (ethtypes.Bloom, error)

	BloomStatus() (uint64, uint64)

	RPCFilterCap() int32
	RPCLogsCap() int32
	RPCBlockRangeCap() int32
}

// consider a filter inactive if it has not been polled for within deadline
const defaultDeadline = 5 * time.Minute

// filter is a helper struct that holds meta information over the filter type
// and associated subscription in the event system.
type filter struct {
	typ      filters.Type
	deadline *time.Timer // filter is inactive when deadline triggers
	crit     filters.FilterCriteria
	offset   int // offset for stream subscription
}

// PublicFilterAPI offers support to create and manage filters. This will allow external clients to retrieve various
// information related to the Ethereum protocol such as blocks, transactions and logs.
type PublicFilterAPI struct {
	logger    log.Logger
	clientCtx client.Context
	backend   Backend
	events    *stream.RPCStream
	filtersMu sync.Mutex
	filters   map[rpc.ID]*filter
	deadline  time.Duration
}

// NewPublicAPI returns a new PublicFilterAPI instance.
func NewPublicAPI(
	logger log.Logger,
	clientCtx client.Context,
	stream *stream.RPCStream,
	backend Backend,
) *PublicFilterAPI {
	api := NewPublicAPIWithDeadline(logger, clientCtx, stream, backend, defaultDeadline)
	return api
}

// NewPublicAPIWithDeadline returns a new PublicFilterAPI instance with the given deadline.
func NewPublicAPIWithDeadline(
	logger log.Logger,
	clientCtx client.Context,
	stream *stream.RPCStream,
	backend Backend,
	deadline time.Duration,
) *PublicFilterAPI {
	logger = logger.With("api", "filter")
	api := &PublicFilterAPI{
		logger:    logger,
		clientCtx: clientCtx,
		backend:   backend,
		filters:   make(map[rpc.ID]*filter),
		events:    stream,
		deadline:  deadline,
	}

	go api.timeoutLoop()

	return api
}

// timeoutLoop runs every 5 minutes and deletes filters that have not been recently used.
// Tt is started when the api is created.
func (api *PublicFilterAPI) timeoutLoop() {
	ticker := time.NewTicker(api.deadline)
	defer ticker.Stop()

	for {
		<-ticker.C
		api.filtersMu.Lock()
		// #nosec G705
		for id, f := range api.filters {
			select {
			case <-f.deadline.C:
				delete(api.filters, id)
			default:
				continue
			}
		}
		api.filtersMu.Unlock()
	}
}

// NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
// as transactions enter the pending state.
//
// It is part of the filter package because this filter can be used through the
// `eth_getFilterChanges` polling method that is also used for log filters.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newPendingTransactionFilter
func (api *PublicFilterAPI) NewPendingTransactionFilter() rpc.ID {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if len(api.filters) >= int(api.backend.RPCFilterCap()) {
		return rpc.ID("error creating pending tx filter: max limit reached")
	}

	id := rpc.NewID()
	_, offset := api.events.PendingTxStream().ReadNonBlocking(-1)
	api.filters[id] = &filter{
		typ:      filters.PendingTransactionsSubscription,
		deadline: time.NewTimer(api.deadline),
		offset:   offset,
	}

	return id
}

// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
// It is part of the filter package since polling goes with eth_getFilterChanges.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter
func (api *PublicFilterAPI) NewBlockFilter() rpc.ID {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if len(api.filters) >= int(api.backend.RPCFilterCap()) {
		return rpc.ID("error creating block filter: max limit reached")
	}

	id := rpc.NewID()
	_, offset := api.events.HeaderStream().ReadNonBlocking(-1)
	api.filters[id] = &filter{
		typ:      filters.BlocksSubscription,
		deadline: time.NewTimer(api.deadline),
		offset:   offset,
	}

	return id
}

// NewFilter creates a new filter and returns the filter id. It can be
// used to retrieve logs when the state changes. This method cannot be
// used to fetch logs that are already stored in the state.
//
// Default criteria for the from and to block are "latest".
// Using "latest" as block number will return logs for mined blocks.
// Using "pending" as block number returns logs for not yet mined (pending) blocks.
// In case logs are removed (chain reorg) previously returned logs are returned
// again but with the removed property set to true.
//
// In case "fromBlock" > "toBlock" an error is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter
func (api *PublicFilterAPI) NewFilter(criteria filters.FilterCriteria) (rpc.ID, error) {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if len(api.filters) >= int(api.backend.RPCFilterCap()) {
		return "", fmt.Errorf("error creating filter: max limit reached")
	}

	id := rpc.NewID()
	_, offset := api.events.LogStream().ReadNonBlocking(-1)
	api.filters[id] = &filter{
		typ:      filters.LogsSubscription,
		deadline: time.NewTimer(api.deadline),
		crit:     criteria,
		offset:   offset,
	}

	return id, nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs
func (api *PublicFilterAPI) GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*ethtypes.Log, error) {
	var filter *Filter
	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.logger, api.backend, crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		// Block numbers below 0 are special cases.
		// for more info, https://github.com/ethereum/go-ethereum/blob/v1.15.11/eth/filters/api.go#L360
		if begin > 0 && end > 0 && begin > end {
			return nil, errInvalidBlockRange
		}
		// Construct the range filter
		filter = NewRangeFilter(api.logger, api.backend, begin, end, crit.Addresses, crit.Topics)
	}

	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx, int(api.backend.RPCLogsCap()), int64(api.backend.RPCBlockRangeCap()))
	if err != nil {
		return nil, err
	}

	return returnLogs(logs), err
}

// UninstallFilter removes the filter with the given filter id.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_uninstallfilter
func (api *PublicFilterAPI) UninstallFilter(id rpc.ID) bool {
	api.filtersMu.Lock()
	_, found := api.filters[id]
	if found {
		delete(api.filters, id)
	}
	api.filtersMu.Unlock()

	return found
}

// GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterlogs
func (api *PublicFilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([]*ethtypes.Log, error) {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	api.filtersMu.Unlock()

	if !found {
		return returnLogs(nil), fmt.Errorf("filter %s not found", id)
	}

	if f.typ != filters.LogsSubscription {
		return returnLogs(nil), fmt.Errorf("filter %s doesn't have a LogsSubscription type: got %d", id, f.typ)
	}

	var filter *Filter
	if f.crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.logger, api.backend, f.crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if f.crit.FromBlock != nil {
			begin = f.crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if f.crit.ToBlock != nil {
			end = f.crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = NewRangeFilter(api.logger, api.backend, begin, end, f.crit.Addresses, f.crit.Topics)
	}
	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx, int(api.backend.RPCLogsCap()), int64(api.backend.RPCBlockRangeCap()))
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), nil
}

// GetFilterChanges returns the logs for the filter with the given id since
// last time it was called. This can be used for polling.
//
// For pending transaction and block filters the result is []common.Hash.
// (pending)Log filters return []Log.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterchanges
func (api *PublicFilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	f, found := api.filters[id]
	if !found {
		return nil, fmt.Errorf("filter %s not found", id)
	}

	if !f.deadline.Stop() {
		// timer expired but filter is not yet removed in timeout loop
		// receive timer value and reset timer
		<-f.deadline.C
	}
	f.deadline.Reset(api.deadline)

	switch f.typ {
	case filters.PendingTransactionsSubscription:
		var hashes []common.Hash
		hashes, f.offset = api.events.PendingTxStream().ReadAllNonBlocking(f.offset)
		return returnHashes(hashes), nil
	case filters.BlocksSubscription:
		var headers []stream.RPCHeader
		headers, f.offset = api.events.HeaderStream().ReadAllNonBlocking(f.offset)
		hashes := make([]common.Hash, len(headers))
		for i, header := range headers {
			hashes[i] = header.Hash
		}
		return hashes, nil
	case filters.LogsSubscription:
		var (
			logs  []*ethtypes.Log
			chunk []*ethtypes.Log
		)
		for {
			chunk, f.offset = api.events.LogStream().ReadNonBlocking(f.offset)
			if len(chunk) == 0 {
				break
			}
			chunk = FilterLogs(chunk, f.crit.FromBlock, f.crit.ToBlock, f.crit.Addresses, f.crit.Topics)
			logs = append(logs, chunk...)
		}
		return returnLogs(logs), nil
	default:
		return nil, fmt.Errorf("invalid filter %s type %d", id, f.typ)
	}
}
