package txpool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/zeta-chain/node/rpc/backend"
	"github.com/zeta-chain/node/rpc/types"

	"cosmossdk.io/log"
)

// PublicAPI offers and API for the transaction pool. It only operates on data that is non-confidential.
// NOTE: For more info about the current status of this endpoints see https://github.com/evmos/ethermint/issues/124
type PublicAPI struct {
	logger  log.Logger
	backend backend.EVMBackend
}

// NewPublicAPI creates a new tx pool service that gives information about the transaction pool.
func NewPublicAPI(logger log.Logger, backend backend.EVMBackend) *PublicAPI {
	return &PublicAPI{
		logger:  logger.With("module", "txpool"),
		backend: backend,
	}
}

// Content returns the transactions contained within the transaction pool
func (api *PublicAPI) Content() (map[string]map[string]map[string]*types.RPCTransaction, error) {
	api.logger.Debug("txpool_content")
	return api.backend.Content()
}

// ContentFrom returns the transactions contained within the transaction pool
func (api *PublicAPI) ContentFrom(address common.Address) (map[string]map[string]*types.RPCTransaction, error) {
	api.logger.Debug("txpool_contentFrom")
	return api.backend.ContentFrom(address)
}

// Inspect returns the content of the transaction pool and flattens it into an easily inspectable list
func (api *PublicAPI) Inspect() (map[string]map[string]map[string]string, error) {
	api.logger.Debug("txpool_inspect")
	return api.backend.Inspect()
}

// Status returns the number of pending and queued transaction in the pool
func (api *PublicAPI) Status() (map[string]hexutil.Uint, error) {
	api.logger.Debug("txpool_status")
	return api.backend.Status()
}
