package backend

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/zeta-chain/node/rpc/types"
)

// Content returns the transactions contained within the transaction pool
func (b *Backend) Content() (map[string]map[string]map[string]*types.RPCTransaction, error) {
	content := map[string]map[string]map[string]*types.RPCTransaction{
		"pending": make(map[string]map[string]*types.RPCTransaction),
		"queued":  make(map[string]map[string]*types.RPCTransaction),
	}
	return content, nil
}

// ContentFrom returns the transactions contained within the transaction pool
func (b *Backend) ContentFrom(_ common.Address) (map[string]map[string]map[string]*types.RPCTransaction, error) {
	content := map[string]map[string]map[string]*types.RPCTransaction{
		"pending": make(map[string]map[string]*types.RPCTransaction),
		"queued":  make(map[string]map[string]*types.RPCTransaction),
	}
	return content, nil
}

// Inspect returns the content of the transaction pool and flattens it into an easily inspectable list.
func (b *Backend) Inspect() (map[string]map[string]map[string]string, error) {
	inspect := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}
	return inspect, nil
}

// Status returns the number of pending and queued transaction in the pool.
func (b *Backend) Status() (map[string]hexutil.Uint, error) {
	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(0),
		"queued":  hexutil.Uint(0),
	}, nil
}
