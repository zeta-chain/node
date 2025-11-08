package backend

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/node/rpc/types"
)

const (
	StatusPending = "pending"
	StatusQueued  = "queued"
)

// The code style for this API is based off of the Go-Ethereum implementation:

// Content returns the transactions contained within the transaction pool.
func (b *Backend) Content() (map[string]map[string]map[string]*types.RPCTransaction, error) {
	content := map[string]map[string]map[string]*types.RPCTransaction{
		StatusPending: make(map[string]map[string]*types.RPCTransaction),
		StatusQueued:  make(map[string]map[string]*types.RPCTransaction),
	}

	// Get current block header
	curHeader, err := b.CurrentHeader()
	if err != nil {
		return content, fmt.Errorf("failed to get current header: %w", err)
	}

	// Get the global mempool instance
	evmMempool := b.Mempool
	if evmMempool == nil {
		return content, nil
	}

	// Get pending (runnable) and queued (blocked) transactions from the mempool
	pending, queued := evmMempool.GetTxPool().Content()

	// Convert pending (pending) transactions
	for addr, txList := range pending {
		addrStr := addr.Hex()
		if content[StatusPending][addrStr] == nil {
			content[StatusPending][addrStr] = make(map[string]*types.RPCTransaction)
		}

		for _, tx := range txList {
			rpcTx := types.NewRPCPendingTransaction(tx, curHeader, b.ChainConfig())
			content[StatusPending][addrStr][strconv.FormatUint(tx.Nonce(), 10)] = rpcTx
		}
	}

	// Convert queued (queued) transactions
	for addr, txList := range queued {
		addrStr := addr.Hex()
		if content[StatusQueued][addrStr] == nil {
			content[StatusQueued][addrStr] = make(map[string]*types.RPCTransaction)
		}

		for _, tx := range txList {
			rpcTx := types.NewRPCPendingTransaction(tx, curHeader, b.ChainConfig())
			content[StatusQueued][addrStr][strconv.FormatUint(tx.Nonce(), 10)] = rpcTx
		}
	}

	return content, nil
}

// ContentFrom returns the transactions contained within the transaction pool
func (b *Backend) ContentFrom(addr common.Address) (map[string]map[string]*types.RPCTransaction, error) {
	content := make(map[string]map[string]*types.RPCTransaction, 2)

	// Get current block header
	curHeader, err := b.CurrentHeader()
	if err != nil {
		return content, fmt.Errorf("failed to get current header: %w", err)
	}

	// Get the global mempool instance
	evmMempool := b.Mempool
	if evmMempool == nil {
		return content, nil
	}

	// Get transactions for the specific address
	pending, queue := evmMempool.GetTxPool().ContentFrom(addr)

	// Build the pending transactions
	dump := make(map[string]*types.RPCTransaction, len(pending)) // variable name comes from go-ethereum: https://github.com/ethereum/go-ethereum/blob/0dacfef8ac42e7be5db26c2956f2b238ba7c75e8/internal/ethapi/api.go#L221
	for _, tx := range pending {
		rpcTx := types.NewRPCPendingTransaction(tx, curHeader, b.ChainConfig())
		dump[fmt.Sprintf("%d", tx.Nonce())] = rpcTx
	}
	content[StatusPending] = dump

	// Build the queued transactions
	dump = make(map[string]*types.RPCTransaction, len(queue)) // variable name comes from go-ethereum: https://github.com/ethereum/go-ethereum/blob/0dacfef8ac42e7be5db26c2956f2b238ba7c75e8/internal/ethapi/api.go#L221
	for _, tx := range queue {
		rpcTx := types.NewRPCPendingTransaction(tx, curHeader, b.ChainConfig())
		dump[fmt.Sprintf("%d", tx.Nonce())] = rpcTx
	}
	content[StatusQueued] = dump

	return content, nil
}

// Inspect returns the content of the transaction pool and flattens it into an easily inspectable list.
func (b *Backend) Inspect() (map[string]map[string]map[string]string, error) {
	inspect := map[string]map[string]map[string]string{
		StatusPending: make(map[string]map[string]string),
		StatusQueued:  make(map[string]map[string]string),
	}

	// Get the global mempool instance
	evmMempool := b.Mempool
	if evmMempool == nil {
		return inspect, nil
	}

	// Get pending (runnable) and queued (blocked) transactions from the mempool
	pending, queued := evmMempool.GetTxPool().Content()

	// Helper function to format transaction for inspection
	format := func(tx *ethtypes.Transaction) string {
		if to := tx.To(); to != nil {
			return fmt.Sprintf("%s: %v wei + %v gas × %v wei",
				tx.To().Hex(), tx.Value(), tx.Gas(), tx.GasPrice())
		}
		return fmt.Sprintf("contract creation: %v wei + %v gas × %v wei",
			tx.Value(), tx.Gas(), tx.GasPrice())
	}

	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		inspect[StatusPending][account.Hex()] = dump
	}

	// Flatten the queued transactions
	for account, txs := range queued {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		inspect[StatusQueued][account.Hex()] = dump
	}

	return inspect, nil
}

// Status returns the number of pending and queued transaction in the pool.
func (b *Backend) Status() (map[string]hexutil.Uint, error) {
	// Get the global mempool instance
	evmMempool := b.Mempool
	if evmMempool == nil {
		return map[string]hexutil.Uint{
			StatusPending: hexutil.Uint(0),
			StatusQueued:  hexutil.Uint(0),
		}, nil
	}

	pending, queued := evmMempool.GetTxPool().Stats()
	return map[string]hexutil.Uint{
		StatusPending: hexutil.Uint(pending), // #nosec G115 -- overflow not a concern for tx counts, as the mempool will limit far before this number is hit. This is taken directly from Geth.
		StatusQueued:  hexutil.Uint(queued),  // #nosec G115 -- overflow not a concern for tx counts, as the mempool will limit far before this number is hit. This is taken directly from Geth.
	}, nil
}
