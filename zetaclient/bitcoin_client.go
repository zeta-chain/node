package zetaclient

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/btc/infra"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"
	"math/big"
	"sync"
	"time"
)

var _ ChainClient = &BitcoinChainClient{}

// Chain configuration struct
// Filled with above constants depending on chain
type BitcoinChainClient struct {
	chain       common.Chain
	endpoint    string
	ticker      *time.Ticker
	chainClient *infra.JSONRpcClient
	zetaClient  *ZetaCoreBridge
	Tss         TSSSigner
	lastBlock   uint64
	confCount   uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime   uint64 // block time in seconds
	txWatchList map[ethcommon.Hash]string
	mu          *sync.Mutex
	db          *leveldb.DB
	metrics     *metricsPkg.Metrics
	stop        chan struct{}
	logger      zerolog.Logger
}

// Return configuration based on supplied target chain
func NewBitcoinClient(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics) (*BitcoinChainClient, error) {
	ob := BitcoinChainClient{}
	ob.stop = make(chan struct{})
	ob.chain = chain
	if !chain.IsBitcoinChain() {
		return nil, fmt.Errorf("chain %s is not a Bitcoin chain", chain)
	}
	ob.mu = &sync.Mutex{}
	ob.logger = log.With().Str("chain", chain.String()).Logger()
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.metrics = metrics

	ob.endpoint = config.Chains[chain.String()].Endpoint

	// initialize the Client
	ob.logger.Info().Msgf("Chain %s endpoint %s", ob.chain, ob.endpoint)
	ob.chainClient = infra.NewJSONRpcClient(ob.endpoint, "tb1q9dlnu5dr254s8xvtzlhk5ttu0c923u623qup39")

	ob.logger.Info().Msgf("%s: start scanning from block %d", chain, ob.GetBlockHeight())

	return &ob, nil
}

func (ob *BitcoinChainClient) Start() {
	//go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	//go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	//go ob.WatchExchangeRate()    // Observers ZetaPriceQuerier for Zeta prices and posts to core
	//go ob.observeOutTx()
	go ob.observeInTx()
}

func (ob *BitcoinChainClient) Stop() {
	ob.logger.Info().Msgf("ob %s is stopping", ob.chain)
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.Info().Msg("closing ob.db")
	err := ob.db.Close()
	if err != nil {
		ob.logger.Error().Err(err).Msg("error closing db")
	}

	ob.logger.Info().Msgf("%s observer stopped", ob.chain)
}

func (ob *BitcoinChainClient) GetBlockHeight() uint64 {
	bn, err := ob.chainClient.GetBlockHeight()
	if err != nil {
		ob.logger.Error().Err(err).Msg("error getting block height")
		return 0
	}
	return uint64(bn)
}

// TODO
func (ob *BitcoinChainClient) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

// TODO
func (ob *BitcoinChainClient) observeInTx() {
	return
}

// TODO
func (ob *BitcoinChainClient) IsSendOutTxProcessed(sendHash string, nonce int, fromOrToZeta bool) (bool, bool, error) {
	return false, false, nil
}

func (ob *BitcoinChainClient) PostNonceIfNotRecorded() error {
	return nil
}
