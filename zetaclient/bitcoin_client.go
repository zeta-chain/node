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
	go ob.WatchInTx()
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

func (ob *BitcoinChainClient) WatchInTx() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := ob.observeInTx()
			if err != nil {
				ob.logger.Error().Err(err).Msg("error observing in tx")
				continue
			}
		case <-ob.stop:
			ob.logger.Info().Msg("WatchInTx stopped")
			return
		}
	}
}

// TODO
func (ob *BitcoinChainClient) observeInTx() error {
	bn := ob.GetBlockHeight()
	if bn == 0 {
		return fmt.Errorf("error getting block height")
	}

	blockHash, err := ob.chainClient.GetBlockHash(int64(bn))
	if err != nil {
		return err
	}

	rawEvents, err := ob.chainClient.GetEventsByHash(blockHash)

	// ============= query the incoming tx to TSS address ==============
	//tssAddress := ob.Tss.EVMAddress()
	tssAddress := ""

	_ = rawEvents
	_ = tssAddress
	// query incoming gas asset
	//for bn := startBlock; bn <= toBlock; bn++ {
	//	block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
	//	if err != nil {
	//		ob.logger.Error().Err(err).Msg("error getting block")
	//		continue
	//	}
	//	for _, tx := range block.Transactions() {
	//		if tx.To() == nil {
	//			continue
	//		}
	//		if *tx.To() == tssAddress {
	//			receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash())
	//			if receipt.Status != 1 { // 1: successful, 0: failed
	//				ob.logger.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
	//				continue
	//			}
	//			if err != nil {
	//				ob.logger.Err(err).Msg("TransactionReceipt")
	//				continue
	//			}
	//			from, err := ob.EvmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
	//			if err != nil {
	//				ob.logger.Err(err).Msg("TransactionSender")
	//				continue
	//			}
	//			ob.logger.Info().Msgf("TSS inTx detected: %s, blocknum %d", tx.Hash().Hex(), receipt.BlockNumber)
	//			ob.logger.Info().Msgf("TSS inTx value: %s", tx.Value().String())
	//			ob.logger.Info().Msgf("TSS inTx from: %s", from.Hex())
	//			message := ""
	//			if len(tx.Data()) != 0 {
	//				message = hex.EncodeToString(tx.Data())
	//			}
	//			zetaHash, err := ob.zetaClient.PostSend(
	//				from.Hex(),
	//				ob.chain.String(),
	//				from.Hex(),
	//				"ZETA",
	//				tx.Value().String(),
	//				tx.Value().String(),
	//				message,
	//				tx.Hash().Hex(),
	//				receipt.BlockNumber.Uint64(),
	//				90_000,
	//				common.CoinType_Gas,
	//			)
	//			if err != nil {
	//				ob.logger.Error().Err(err).Msg("error posting to zeta core")
	//				continue
	//			}
	//			ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	//		}
	//	}
	//}
	//// ============= end of query the incoming tx to TSS address ==============
	//
	////ob.LastBlock = toBlock
	//ob.SetBlockHeight(toBlock)
	//buf := make([]byte, binary.MaxVarintLen64)
	//n := binary.PutUvarint(buf, toBlock)
	//err = ob.db.Put([]byte(PosKey), buf[:n], nil)
	//if err != nil {
	//	ob.logger.Error().Err(err).Msg("error writing toBlock to db")
	//}
	return nil

}

// TODO
func (ob *BitcoinChainClient) IsSendOutTxProcessed(sendHash string, nonce int, fromOrToZeta bool) (bool, bool, error) {
	return false, false, nil
}

func (ob *BitcoinChainClient) PostNonceIfNotRecorded() error {
	return nil
}
