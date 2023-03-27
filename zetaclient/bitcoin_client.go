package zetaclient

import (
	"cosmossdk.io/math"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"

	math2 "math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

var _ ChainClient = &BitcoinChainClient{}

type BTCLog struct {
	ChainLogger   zerolog.Logger
	WatchInTx     zerolog.Logger
	ObserveOutTx  zerolog.Logger
	WatchUTXOS    zerolog.Logger
	WatchGasPrice zerolog.Logger
}

// Chain configuration struct
// Filled with above constants depending on chain
type BitcoinChainClient struct {
	*ChainMetrics

	chain        common.Chain
	endpoint     string
	rpcClient    *rpcclient.Client
	zetaClient   *ZetaCoreBridge
	Tss          TSSSigner
	lastBlock    int64
	confCount    int64                                   // must wait this many blocks to be considered "confirmed"
	BlockTime    uint64                                  // block time in seconds
	submittedTx  map[string]btcjson.GetTransactionResult // key: chain-nonce
	mu           *sync.Mutex
	utxos        []btcjson.ListUnspentResult
	pendingUtxos *leveldb.DB // key is txid_outpoint, value is ListUnspentResult
	stop         chan struct{}
	logger       BTCLog
}

const (
	minConfirmations = 1
	chunkSize        = 500
)

// Return configuration based on supplied target chain
func NewBitcoinClient(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics, logger zerolog.Logger) (*BitcoinChainClient, error) {
	ob := BitcoinChainClient{
		ChainMetrics: NewChainMetrics(chain.String(), metrics),
	}
	ob.stop = make(chan struct{})
	ob.chain = chain
	if !chain.IsBitcoinChain() {
		return nil, fmt.Errorf("chain %s is not a Bitcoin chain", chain.ChainName)
	}
	ob.mu = &sync.Mutex{}
	chainLogger := logger.With().Str("chain", chain.ChainName.String()).Logger()
	ob.logger = BTCLog{
		ChainLogger:   chainLogger,
		WatchInTx:     chainLogger.With().Str("module", "WatchInTx").Logger(),
		ObserveOutTx:  chainLogger.With().Str("module", "observeOutTx").Logger(),
		WatchUTXOS:    chainLogger.With().Str("module", "WatchUTXOS").Logger(),
		WatchGasPrice: chainLogger.With().Str("module", "WatchGasPrice").Logger(),
	}

	ob.zetaClient = bridge
	ob.Tss = tss
	ob.confCount = 0
	ob.submittedTx = make(map[string]btcjson.GetTransactionResult)

	path := fmt.Sprintf("%s/btc_utxos.pendingUtxos", dbpath)
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	ob.pendingUtxos = db
	ob.endpoint = config.BitcoinConfig.RPCEndpoint

	// initialize the Client
	ob.logger.ChainLogger.Info().Msgf("Chain %s endpoint %s", ob.chain.String(), ob.endpoint)
	connCfg := &rpcclient.ConnConfig{
		Host:         ob.endpoint,
		User:         config.BitcoinConfig.RPCUsername,
		Pass:         config.BitcoinConfig.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       config.BitcoinConfig.RPCParams,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating rpc client: %s", err)
	}
	ob.rpcClient = client
	err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("error ping the bitcoin server: %s", err)
	}
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return nil, fmt.Errorf("error getting block count: %s", err)
	}

	ob.logger.ChainLogger.Info().Msgf("%s: start scanning from block %d", chain.String(), bn)

	envvar := ob.chain.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		ob.logger.ChainLogger.Info().Msgf("envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == clienttypes.EnvVarLatest {
			ob.SetLastBlockHeight(bn)
		} else {
			scanFromBlockInt, err := strconv.ParseInt(scanFromBlock, 10, 64)
			if err != nil {
				return nil, err
			}
			ob.SetLastBlockHeight(scanFromBlockInt)
		}
	}
	if ob.chain.ChainId == 18444 { // bitcoin regtest: start from block 100
		ob.SetLastBlockHeight(100)
	}

	return &ob, nil
}

func (ob *BitcoinChainClient) Start() {
	ob.logger.ChainLogger.Info().Msgf("BitcoinChainClient is starting")
	go ob.WatchInTx()
	go ob.observeOutTx()
	go ob.WatchUTXOS()
	go ob.WatchGasPrice()
}

func (ob *BitcoinChainClient) Stop() {
	ob.logger.ChainLogger.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop
	//
	//ob.Logger.Info().Msg("closing ob.pendingUtxos")
	//err := ob.pendingUtxos.Close()
	//if err != nil {
	//	ob.Logger.Error().Err(err).Msg("error closing pendingUtxos")
	//}
	//
	ob.logger.ChainLogger.Info().Msgf("%s observer stopped", ob.chain.String())
}

func (ob *BitcoinChainClient) SetLastBlockHeight(block int64) {
	if block < 0 {
		panic("lastBlock is negative")
	}
	if block >= math2.MaxInt64 {
		panic("lastBlock is too large")
	}
	atomic.StoreInt64(&ob.lastBlock, block)
}

func (ob *BitcoinChainClient) GetLastBlockHeight() int64 {
	height := atomic.LoadInt64(&ob.lastBlock)
	if height < 0 {
		panic("lastBlock is negative")
	}
	if height >= math2.MaxInt64 {
		panic("lastBlock is too large")
	}
	return height
}

// TODO
func (ob *BitcoinChainClient) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

func (ob *BitcoinChainClient) WatchInTx() {
	//ob.logger = ob.logger.With().Str("module", "WatchInTx").Logger()
	ticker := time.NewTicker(time.Duration(config.BitcoinConfig.WatchInTxPeriod) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := ob.observeInTx()
			if err != nil {
				ob.logger.WatchInTx.Error().Err(err).Msg("error observing in tx")
				continue
			}
		case <-ob.stop:
			ob.logger.WatchInTx.Info().Msg("WatchInTx stopped")
			return
		}
	}
}

// TODO
func (ob *BitcoinChainClient) observeInTx() error {
	permssions, err := ob.zetaClient.GetInboundPermissions()
	if err != nil {
		return err
	}
	if !permssions.IsInboundEnabled {
		return errors.New("inbound TXS / Send has been disabled by the protocol")
	}
	lastBN := ob.GetLastBlockHeight()
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("error getting block count: %s", err)
	}
	if cnt < 0 || cnt >= math2.MaxInt64 {
		return fmt.Errorf("block count is out of range: %d", cnt)
	}
	// query incoming gas asset
	if cnt > lastBN {
		bn := lastBN + 1
		ob.logger.WatchInTx.Info().Msgf("filtering block %d, current block %d, last block %d", bn, cnt, lastBN)
		hash, err := ob.rpcClient.GetBlockHash(bn)
		if err != nil {
			return err
		}

		block, err := ob.rpcClient.GetBlockVerboseTx(hash)
		if err != nil {
			return err
		}
		ob.logger.WatchInTx.Info().Msgf("block %d has %d txs", bn, len(block.Tx))
		if len(block.Tx) > 1 {
			for idx, tx := range block.Tx {
				ob.logger.WatchInTx.Info().Msgf("BTC InTX |  %d: %s\n", idx, tx.Txid)
				for vidx, vout := range tx.Vout {
					ob.logger.WatchInTx.Debug().Msgf("vout %d \n value: %v\n scriptPubKey: %v\n", vidx, vout.Value, vout.ScriptPubKey.Hex)
				}
			}
		}

		tssAddress := ob.Tss.BTCAddress()
		inTxs := FilterAndParseIncomingTx(block.Tx, uint64(block.Height), tssAddress, &ob.logger.WatchInTx)

		for _, inTx := range inTxs {
			ob.logger.WatchInTx.Debug().Msgf("Processing inTx: %s", inTx.TxHash)
			amount := big.NewFloat(inTx.Value)
			amount = amount.Mul(amount, big.NewFloat(1e8))
			amountInt, _ := amount.Int(nil)
			message := hex.EncodeToString(inTx.MemoBytes)
			zetaHash, err := ob.zetaClient.PostSend(
				inTx.FromAddress,
				ob.chain.ChainId,
				inTx.FromAddress,
				inTx.FromAddress,
				common.ZetaChain().ChainId,
				math.NewUintFromBigInt(amountInt),
				message,
				inTx.TxHash,
				inTx.BlockNumber,
				0,
				common.CoinType_Gas,
				PostSendEVMGasLimit,
				"",
			)
			if err != nil {
				ob.logger.WatchInTx.Error().Err(err).Msg("error posting to zeta core")
				continue
			}
			ob.logger.WatchInTx.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
		}

		ob.SetLastBlockHeight(bn)
	}

	return nil
}

// returns isIncluded, isConfirmed, Error
func (ob *BitcoinChainClient) IsSendOutTxProcessed(sendHash string, nonce int, _ common.CoinType, logger zerolog.Logger) (bool, bool, error) {
	chain := ob.chain.ChainId
	outTxID := fmt.Sprintf("%d-%d", chain, nonce)
	logger.Info().Msgf("IsSendOutTxProcessed %s", outTxID)

	res, found := ob.submittedTx[outTxID]
	if !found {
		return false, false, nil
	}
	if res.Confirmations == 0 {
		return true, false, nil
	} else if res.Confirmations > 0 { // FIXME: use configured block confirmation
		amountInSat, _ := big.NewFloat(res.Amount * 1e8).Int(nil)
		zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
			sendHash,
			res.TxID,
			uint64(res.BlockIndex),
			amountInSat,
			common.ReceiveStatus_Success,
			ob.chain,
			nonce,
			common.CoinType_Gas,
		)
		if err != nil {
			logger.Error().Err(err).Msgf("error posting to zeta core")
		} else {
			logger.Info().Msgf("Bitcoin outTx confirmed: PostReceiveConfirmation zeta tx: %s", zetaHash)
		}
		return true, true, nil
	}
	return false, false, nil
}

// FIXME: bitcoin tx does not have nonce; however, nonce can be maintained
// by the client to easily identify the cctx outbound command
func (ob *BitcoinChainClient) PostNonceIfNotRecorded(logger zerolog.Logger) error {
	zetaHash, err := ob.zetaClient.PostNonce(ob.chain, 0)
	if err != nil {
		return errors.Wrap(err, "error posting nonce to zeta core")
	}
	logger.Info().Msgf("PostNonce zeta tx %s , signer %s , nonce %d", zetaHash, ob.zetaClient.keys.GetOperatorAddress(), 0)
	return nil
}

func (ob *BitcoinChainClient) WatchGasPrice() {
	//ob.logger = ob.logger.With().Str("module", "WatchGasPrice").Logger()

	gasTicker := time.NewTicker(time.Duration(config.BitcoinConfig.WatchGasPricePeriod) * time.Second)
	for {
		select {
		case <-gasTicker.C:
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.WatchGasPrice.Error().Err(err).Msg("PostGasPrice error on " + ob.chain.String())
				continue
			}
		case <-ob.stop:
			ob.logger.WatchGasPrice.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) PostGasPrice() error {
	if ob.chain.ChainId == 18444 { //bitcoin regtest
		bn, err := ob.rpcClient.GetBlockCount()
		if err != nil {
			return err
		}
		zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, 1000, "100", uint64(bn))
		if err != nil {
			ob.logger.WatchGasPrice.Err(err).Msg("PostGasPrice:")
			return err
		}
		ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)
		return nil
	}
	// EstimateSmartFee returns the fees per kilobyte (BTC/kb) targeting given block confirmation
	feeResult, err := ob.rpcClient.EstimateSmartFee(1, &btcjson.EstimateModeConservative)
	if err != nil {
		return err
	}
	if feeResult.Errors != nil || feeResult.FeeRate == nil {
		return fmt.Errorf("error getting gas price: %s", feeResult.Errors)
	}
	gasPrice := big.NewFloat(0)
	gasPriceU64, _ := gasPrice.Mul(big.NewFloat(*feeResult.FeeRate), big.NewFloat(1e8)).Uint64()
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}
	zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, gasPriceU64, "100", uint64(bn))
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("PostGasPrice:")
		return err
	}
	ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)
	_ = feeResult
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
func FilterAndParseIncomingTx(txs []btcjson.TxRawResult, blockNumber uint64, targetAddress string, logger *zerolog.Logger) []*BTCInTxEvnet {
	inTxs := make([]*BTCInTxEvnet, 0)
	for _, tx := range txs {
		found := false
		var value float64
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
				wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(hash, config.BitconNetParams)
				if err != nil {
					continue
				}
				if wpkhAddress.EncodeAddress() != targetAddress {
					continue
				}
				value = out.Value
				out = tx.Vout[1]
				script = out.ScriptPubKey.Hex
				if len(script) >= 4 && script[:2] == "6a" { // OP_RETURN
					memoSize, err := strconv.ParseInt(script[2:4], 16, 32)
					if err != nil {
						logger.Warn().Err(err).Msgf("error decoding pubkey hash")
						continue
					}
					if int(memoSize) != (len(script)-4)/2 {
						logger.Warn().Msgf("memo size mismatch: %d != %d", memoSize, (len(script)-4)/2)
						continue
					}
					memoStr, err := hex.DecodeString(script[4:])
					if err != nil {
						logger.Warn().Err(err).Msgf("error hex decoding memo")
						continue
					}
					memoBytes, err := base64.StdEncoding.DecodeString(string(memoStr))
					if err != nil {
						logger.Warn().Err(err).Msgf("error b64 decoding memoStr %x", memoStr)
						continue
					}
					memo = memoBytes
					found = true

				}
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
						logger.Warn().Msgf("error decoding pubkey: %s", err)
						break
					}
					hash := btcutil.Hash160(pkBytes)
					addr, err := btcutil.NewAddressWitnessPubKeyHash(hash, config.BitconNetParams)
					if err != nil {
						logger.Warn().Msgf("error decoding pubkey hash: %s", err)
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

func (ob *BitcoinChainClient) WatchUTXOS() {
	//ob.logger = ob.logger.With().Str("module", "WatchUTXOS").Logger()

	ticker := time.NewTicker(time.Duration(config.BitcoinConfig.WatchUTXOSPeriod) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := ob.fetchUTXOS()
			if err != nil {
				ob.logger.WatchUTXOS.Error().Err(err).Msg("error fetching btc utxos")
				continue
			}
		case <-ob.stop:
			ob.logger.WatchUTXOS.Info().Msg("WatchUTXOS stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) fetchUTXOS() error {
	// get the current block height.
	bh, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("btc: error getting block height : %v", err)
	}
	maxConfirmations := int(bh)

	// List unspent.
	tssAddr := ob.Tss.BTCAddress()
	address, err := btcutil.DecodeAddress(tssAddr, config.BitconNetParams)
	if err != nil {
		return fmt.Errorf("btc: error decoding wallet address (%s) : %s", tssAddr, err.Error())
	}
	addresses := []btcutil.Address{address}
	var utxos []btcjson.ListUnspentResult

	// populate utxos array
	for i := minConfirmations; i < maxConfirmations; i += chunkSize {
		unspents, err := ob.rpcClient.ListUnspentMinMaxAddresses(i, i+chunkSize, addresses)
		if err != nil {
			return err
		}
		utxos = append(utxos, unspents...)
		ob.logger.WatchUTXOS.Debug().Msgf("btc: fetched %d utxos", len(unspents))
		//for idx, utxo := range unspents {
		//	fmt.Printf("utxo %d\n", idx)
		//	fmt.Printf("  txid: %s\n", utxo.TxID)
		//	fmt.Printf("  address: %s\n", utxo.Address)
		//	fmt.Printf("  amount: %f\n", utxo.Amount)
		//	fmt.Printf("  confirmations: %d\n", utxo.Confirmations)
		//}
	}
	// filter pending
	var filtered []btcjson.ListUnspentResult
	for _, utxo := range utxos {
		pending, err := ob.isPending(utxoKey(utxo))
		if err != nil {
			return fmt.Errorf("btc: error accessing pending utxos pendingUtxos: %v", err.Error())
		}
		if !pending {
			filtered = append(filtered, utxo)
		}
	}
	// sort by value
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Amount < filtered[j].Amount
	})
	ob.utxos = filtered
	// remove completed from pending pendingUtxos
	ob.housekeepPending()
	return nil
}

func (ob *BitcoinChainClient) housekeepPending() {
	// create map with utxos
	utxosMap := make(map[string]bool, len(ob.utxos))
	for _, utxo := range ob.utxos {
		utxosMap[utxoKey(utxo)] = true
	}

	// traverse pending pendingUtxos
	var removed int64
	iter := ob.pendingUtxos.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		// if key not in utxos map, remove from pendingUtxos
		if !utxosMap[string(key)] {
			if err := ob.pendingUtxos.Delete(key, nil); err != nil {
				ob.logger.WatchUTXOS.Warn().Err(err).Msgf("btc: error removing key [%s] from pending utxos pendingUtxos", key)
			}
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		ob.logger.WatchUTXOS.Warn().Err(err).Msgf("btc: pending utxos housekeeping")
	}
	if removed > 0 {
		ob.logger.WatchUTXOS.Info().Msgf("btc : %d txs purged from pending pendingUtxos", removed)
	}
}

func (ob *BitcoinChainClient) isPending(utxoKey string) (bool, error) {
	if _, err := ob.pendingUtxos.Get([]byte(utxoKey), nil); err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ob *BitcoinChainClient) observeOutTx() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain)
			if err != nil {
				ob.logger.ObserveOutTx.Error().Err(err).Msg("error GetAllOutTxTrackerByChain")
				continue
			}
			for _, tracker := range trackers {
				outTxID := fmt.Sprintf("%d-%d", tracker.ChainId, tracker.Nonce)
				ob.logger.ObserveOutTx.Info().Msgf("tracker outTxID: %s", outTxID)
				for _, txHash := range tracker.HashList {
					hash, err := chainhash.NewHashFromStr(txHash.TxHash)
					if err != nil {
						ob.logger.ObserveOutTx.Error().Err(err).Msg("error NewHashFromStr")
						continue
					}
					getTxResult, err := ob.rpcClient.GetTransaction(hash)
					if err != nil {
						ob.logger.ObserveOutTx.Warn().Err(err).Msg("error GetTransaction")
						continue
					}
					if getTxResult.Confirmations >= 0 {
						ob.submittedTx[outTxID] = *getTxResult
					}
				}
			}
		case <-ob.stop:
			ob.logger.ObserveOutTx.Info().Msg("observeOutTx stopped")
			return
		}
	}
}
