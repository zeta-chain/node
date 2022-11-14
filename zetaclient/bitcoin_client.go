package zetaclient

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"
	"math/big"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var _ ChainClient = &BitcoinChainClient{}

// Chain configuration struct
// Filled with above constants depending on chain
type BitcoinChainClient struct {
	chain       common.Chain
	endpoint    string
	ticker      *time.Ticker
	rpcClient   *rpcclient.Client
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

	connCfg := &rpcclient.ConnConfig{
		Host:         ob.endpoint,
		User:         "user",
		Pass:         "pass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating rpc client: %s", err)
	}
	ob.rpcClient = client
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return nil, fmt.Errorf("error getting block count: %s", err)
	}

	ob.logger.Info().Msgf("%s: start scanning from block %d", chain, bn)

	envvar := ob.chain.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		ob.logger.Info().Msgf("envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == "latest" {
			ob.SetLastBlockHeight(uint64(bn))
		} else {
			scanFromBlockInt, err := strconv.ParseInt(scanFromBlock, 10, 64)
			if err != nil {
				return nil, err
			}
			ob.SetLastBlockHeight(uint64(scanFromBlockInt))
		}
	}

	return &ob, nil
}

func (ob *BitcoinChainClient) Start() {
	go ob.WatchInTx()
}

func (ob *BitcoinChainClient) Stop() {
	ob.logger.Info().Msgf("ob %s is stopping", ob.chain)
	close(ob.stop) // this notifies all goroutines to stop
	//
	//ob.logger.Info().Msg("closing ob.db")
	//err := ob.db.Close()
	//if err != nil {
	//	ob.logger.Error().Err(err).Msg("error closing db")
	//}
	//
	ob.logger.Info().Msgf("%s observer stopped", ob.chain)
}

func (ob *BitcoinChainClient) SetLastBlockHeight(block uint64) {
	atomic.StoreUint64(&ob.lastBlock, block)
}

func (ob *BitcoinChainClient) GetLastBlockHeight() uint64 {
	return atomic.LoadUint64(&ob.lastBlock)
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
	lastBN := ob.GetLastBlockHeight()
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("error getting block count: %s", err)
	}
	bn := uint64(cnt)
	// query incoming gas asset
	if bn > lastBN {
		hash, err := ob.rpcClient.GetBlockHash(int64(bn))
		if err != nil {
			return err
		}

		block, err := ob.rpcClient.GetBlockVerboseTx(hash)
		if err != nil {
			return err
		}

		tssAddress := ob.Tss.BTCAddress()
		inTxs := FilterAndParseIncomingTx(block.Tx, uint64(block.Height), tssAddress)

		for _, inTx := range inTxs {
			//ob.logger.Info().Msgf("incoming tx %v", inTx)
			amount := big.NewFloat(inTx.Value)
			amount = amount.Mul(amount, big.NewFloat(1e8))
			amountInt, _ := amount.Int(nil)
			message := hex.EncodeToString(inTx.MemoBytes)
			zetaHash, err := ob.zetaClient.PostSend(
				inTx.FromAddress,
				ob.chain.String(),
				inTx.FromAddress,
				"ZETA",
				amountInt.String(),
				amountInt.String(),
				message,
				inTx.TxHash,
				inTx.BlockNumber,
				0,
				common.CoinType_Gas,
			)
			if err != nil {
				ob.logger.Error().Err(err).Msg("error posting to zeta core")
				continue
			}
			ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
		}

		ob.SetLastBlockHeight((bn))
	}

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

func (ob *BitcoinChainClient) WatchGasPrice() {
	gasTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-gasTicker.C:
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.Error().Err(err).Msg("PostGasPrice error on " + ob.chain.String())
				continue
			}
		case <-ob.stop:
			ob.logger.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) PostGasPrice() error {
	//
	//
	//_, err = ob.zetaClient.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	//if err != nil {
	//	ob.logger.Err(err).Msg("PostGasPrice:")
	//	return err
	//}

	return nil
}

type BTCInTxEvnet struct {
	FromAddress string  // the first input address
	ToAddress   string  // some TSS address
	Value       float64 // in BTC, not satoshi
	MemoBytes   []byte
	BlockNumber uint64
	TxHash      string
}

// given txs list returned by the "getblock 2" RPC command, return the txs that are relevant to us
// relevant tx must have the following vouts as the first two vouts:
// vout0: p2wpkh to the TSS address (targetAddress)
// vout1: OP_RETURN memo, base64 encoded
func FilterAndParseIncomingTx(txs []btcjson.TxRawResult, blockNumber uint64, targetAddress string) []*BTCInTxEvnet {
	inTxs := make([]*BTCInTxEvnet, 0)
	for _, tx := range txs {
		found := false
		var value float64 = 0
		var memo []byte
		if len(tx.Vout) >= 2 {
			// first vout must to addressed to the targetAddress with p2wpkh scriptPubKey
			out := tx.Vout[0]
			script := out.ScriptPubKey.Hex
			if len(script) == 44 && script[:4] == "0014" { // segwit output: 0x00 + 20 bytes of pubkey hash
				hash, err := hex.DecodeString(script[4:])
				if err != nil {
					continue
				}
				wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(hash, &chaincfg.TestNet3Params)
				if err != nil {
					continue
				}
				if wpkhAddress.EncodeAddress() != targetAddress {
					continue
				}
				value = out.Value
			}
			out = tx.Vout[1]
			script = out.ScriptPubKey.Hex
			if len(script) >= 4 && script[:2] == "6a" { // OP_RETURN
				memoSize, err := strconv.ParseInt(script[2:4], 16, 32)
				if err != nil {
					continue
				}
				if int(memoSize) != len(script)/2-2 {
					continue
				}
				memoStr, err := hex.DecodeString(script[4:])
				memoBytes, err := base64.StdEncoding.DecodeString(string(memoStr))
				if err != nil {
					continue
				}
				memo = memoBytes
				found = true
			}
		}
		if found {
			var fromAddress string
			if len(tx.Vin) > 0 {
				vin := tx.Vin[0]
				//log.Info().Msgf("vin: %v", vin.Witness)
				if len(vin.Witness) == 2 {
					pk := vin.Witness[1]
					pkBytes, err := hex.DecodeString(pk)
					if err != nil {
						break
					}
					hash := btcutil.Hash160(pkBytes)
					addr, err := btcutil.NewAddressWitnessPubKeyHash(hash, &chaincfg.TestNet3Params)
					if err != nil {
						break
					}
					fromAddress = addr.EncodeAddress()
				}
			}
			inTxs = append(inTxs, &BTCInTxEvnet{
				FromAddress: fromAddress,
				ToAddress:   targetAddress,
				Value:       value,
				MemoBytes:   memo,
				BlockNumber: blockNumber,
				TxHash:      tx.Txid,
			})
		}
	}
	return inTxs
}
