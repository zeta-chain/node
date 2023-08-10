package zetaclient

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	chain         common.Chain
	rpcClient     *rpcclient.Client
	zetaClient    *ZetaCoreBridge
	Tss           TSSSigner
	lastBlock     int64
	BlockTime     uint64                                  // block time in seconds
	includedTx    map[string]btcjson.GetTransactionResult // key: chain-nonce
	broadcastedTx map[string]chainhash.Hash
	mu            *sync.Mutex
	utxos         []btcjson.ListUnspentResult
	db            *gorm.DB
	stop          chan struct{}
	logger        BTCLog
	cfg           *config.Config
	ts            *TelemetryServer
}

const (
	minConfirmations = 0
	chunkSize        = 1000
	maxHeightDiff    = 10000
	dustOffset       = 2000
)

func (ob *BitcoinChainClient) GetChainConfig() *config.BTCConfig {
	return ob.cfg.BitcoinConfig
}

func (ob *BitcoinChainClient) GetRPCHost() string {
	return ob.GetChainConfig().RPCHost
}

// Return configuration based on supplied target chain
func NewBitcoinClient(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics, logger zerolog.Logger, cfg *config.Config, ts *TelemetryServer) (*BitcoinChainClient, error) {
	ob := BitcoinChainClient{
		ChainMetrics: NewChainMetrics(chain.String(), metrics),
		ts:           ts,
	}
	ob.cfg = cfg
	ob.stop = make(chan struct{})
	ob.chain = chain
	if !common.IsBitcoinChain(chain.ChainId) {
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
	ob.includedTx = make(map[string]btcjson.GetTransactionResult)
	ob.broadcastedTx = make(map[string]chainhash.Hash)

	// initialize the Client
	ob.logger.ChainLogger.Info().Msgf("Chain %s endpoint %s", ob.chain.String(), ob.GetRPCHost())
	connCfg := &rpcclient.ConnConfig{
		Host:         ob.GetRPCHost(),
		User:         ob.GetChainConfig().RPCUsername,
		Pass:         ob.GetChainConfig().RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       ob.GetChainConfig().RPCParams,
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

	//Load btc chain client DB
	err = ob.loadDB(dbpath)
	if err != nil {
		return nil, err
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
	ob.ts.SetLastScannedBlockNumber((ob.chain.ChainId), (block))
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

// GetBaseGasPrice ...
// TODO: implement
// https://github.com/zeta-chain/node/issues/868
func (ob *BitcoinChainClient) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

func (ob *BitcoinChainClient) WatchInTx() {
	ticker := time.NewTicker(time.Duration(ob.GetChainConfig().CoreParams.InTxTicker) * time.Second)
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

	// "confirmed" current block number
	confirmedBlockNum := cnt - int64(ob.GetChainConfig().CoreParams.ConfCount)
	if confirmedBlockNum < 0 || confirmedBlockNum > math2.MaxInt64 {
		return fmt.Errorf("skipping observer , confirmedBlockNum is negative or too large ")
	}

	// query incoming gas asset
	if confirmedBlockNum > lastBN {
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
				//ob.rpcClient.GetTransaction(tx.Txid)
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

		// Save LastBlockHeight
		ob.SetLastBlockHeight(bn)
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeight())).Error; err != nil {
			ob.logger.WatchInTx.Error().Err(err).Msg("error writing Block to db")
		}
	}

	return nil
}

// Returns number of required Bitcoin confirmations depending on sent BTC amount.
func (ob *BitcoinChainClient) ConfirmationsThreshold(amount *big.Int) int64 {
	if amount.Cmp(big.NewInt(200000000)) >= 0 {
		return 6
	}
	return 2
}

// returns isIncluded, isConfirmed, Error
func (ob *BitcoinChainClient) IsSendOutTxProcessed(sendHash string, nonce int, _ common.CoinType, logger zerolog.Logger) (bool, bool, error) {
	outTxID := ob.GetTxID(uint64(nonce))
	logger.Info().Msgf("IsSendOutTxProcessed %s", outTxID)

	ob.mu.Lock()
	txnHash, broadcasted := ob.broadcastedTx[outTxID]
	res, included := ob.includedTx[outTxID]
	ob.mu.Unlock()

	if !included {
		if !broadcasted {
			return false, false, nil
		}
		//Query txn hash on bitcoin chain
		hash, err := chainhash.NewHashFromStr(txnHash.String())
		if err != nil {
			return false, false, nil
		}
		getTxResult, err := ob.rpcClient.GetTransactionWatchOnly(hash, true)
		if err != nil {
			ob.logger.ObserveOutTx.Warn().Err(err).Msg("IsSendOutTxProcessed: transaction not found")
			return false, false, nil
		}
		res = *getTxResult

		// Save result to avoid unnecessary query
		ob.mu.Lock()
		ob.includedTx[outTxID] = res
		ob.mu.Unlock()
	}

	var amount float64
	if res.Amount > 0 {
		ob.logger.ObserveOutTx.Warn().Msg("IsSendOutTxProcessed: res.Amount > 0")
		amount = res.Amount
	} else if res.Amount == 0 {
		ob.logger.ObserveOutTx.Error().Msg("IsSendOutTxProcessed: res.Amount == 0")
		return false, false, nil
	} else {
		amount = -res.Amount
	}

	amountInSat, _ := big.NewFloat(amount * 1e8).Int(nil)
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		return true, false, nil
	}

	logger.Debug().Msgf("Bitcoin outTx confirmed: txid %s, amount %f\n", res.TxID, res.Amount)
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

func (ob *BitcoinChainClient) WatchGasPrice() {

	gasTicker := time.NewTicker(time.Duration(ob.GetChainConfig().CoreParams.GasPriceTicker) * time.Second)
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
		_ = zetaHash
		//ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)
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
	_ = zetaHash
	//ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)
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
	for idx, tx := range txs {
		if idx == 0 {
			continue // the first tx is coinbase; we do not process coinbase tx
		}
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
					memoBytes, err := hex.DecodeString(script[4:])
					if err != nil {
						logger.Warn().Err(err).Msgf("error hex decoding memo")
						continue
					}
					if bytes.Compare(memoBytes, []byte(DonationMessage)) == 0 {
						logger.Info().Msgf("donation tx: %s; value %f", tx.Txid, value)
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

	ticker := time.NewTicker(time.Duration(ob.GetChainConfig().CoreParams.WatchUTXOTicker) * time.Second)
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
	defer func() {
		if err := recover(); err != nil {
			ob.logger.WatchUTXOS.Error().Msgf("BTC fetchUTXOS: caught panic error: %v", err)
		}
	}()
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
		//ob.logger.WatchUTXOS.Debug().Msgf("btc: fetched %d utxos", len(unspents))
		//for idx, utxo := range unspents {
		//	fmt.Printf("utxo %d\n", idx)
		//	fmt.Printf("  txid: %s\n", utxo.TxID)
		//	fmt.Printf("  address: %s\n", utxo.Address)
		//	fmt.Printf("  amount: %f\n", utxo.Amount)
		//	fmt.Printf("  confirmations: %d\n", utxo.Confirmations)
		//}
	}
	// sort by value
	sort.SliceStable(utxos, func(i, j int) bool {
		return utxos[i].Amount < utxos[j].Amount
	})

	ob.mu.Lock()
	ob.ts.SetNumberOfUTXOs(len(utxos))
	ob.utxos = utxos
	ob.mu.Unlock()
	return nil
}

func (ob *BitcoinChainClient) findNonceMarkUTXO(nonce uint64, tssAddress string) (int, error) {
	outTxID := ob.GetTxID(nonce)
	res, mined := ob.includedTx[outTxID]
	if !mined {
		return -1, fmt.Errorf("findNonceMarkUTXO: outTx %s not included yet", outTxID)
	}

	amount := NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := getSatoshis(utxo.Amount)
		if err != nil {
			ob.logger.ObserveOutTx.Error().Err(err).Msgf("findNonceMarkUTXO: error getting satoshis for utxo %v", utxo)
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == res.TxID {
			ob.logger.ObserveOutTx.Info().Msgf("findNonceMarkUTXO: found nonce-mark utxo with txid %s, amount %v", utxo.TxID, utxo.Amount)
			return i, nil
		}
	}
	return -1, fmt.Errorf("findNonceMarkUTXO: cannot find nonce-mark utxo with nonce %d", nonce)
}

// Selects a sublist of utxos to be used as inputs.
//
// Parameters:
//   - amount: The desired minimum total value of the selected UTXOs.
//   - utxoCap: The maximum number of UTXOs to be selected.
//   - nonce: The nonce of the outbound transaction.
//   - tssAddress: The TSS address.
//
// Returns: a sublist (includes previous nonce-mark) of UTXOs or an error if the qulifying sublist cannot be found.
func (ob *BitcoinChainClient) SelectUTXOs(amount float64, utxoCap uint8, nonce uint64, tssAddress string) ([]btcjson.ListUnspentResult, float64, error) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// for nonce > 0; we proceed only when we see the nonce-mark utxo
	// for nonce = 0; make exception; no need to include nonce-mark utxo
	idx := -1
	if nonce > 0 {
		index, err := ob.findNonceMarkUTXO(nonce-1, tssAddress)
		if err != nil {
			return nil, 0, err
		}
		idx = index
	}

	// select utxos
	total := 0.0
	left, right := 0, 0
	for total < amount && right < len(ob.utxos) {
		if utxoCap > 0 { // expand sublist
			total += ob.utxos[right].Amount
			right++
			utxoCap--
		} else { // pop the smallest utxo and append the current one
			total -= ob.utxos[left].Amount
			total += ob.utxos[right].Amount
			left++
			right++
		}
	}
	results := ob.utxos[left:right]

	// include nonce-mark utxo (for nonce > 0) in asending order
	if idx >= 0 {
		if idx < left {
			total += ob.utxos[idx].Amount
			results = append([]btcjson.ListUnspentResult{ob.utxos[idx]}, results...)
		}
		if idx >= right {
			total += ob.utxos[idx].Amount
			results = append(results, ob.utxos[idx])
		}
	}
	if total < amount {
		return nil, 0, fmt.Errorf("SelectUTXOs: not enough btc in reserve - available : %v , tx amount : %v", total, amount)
	}
	return results, total, nil
}

// Save successfully broadcasted transaction
func (ob *BitcoinChainClient) SaveBroadcastedTx(txHash chainhash.Hash, nonce uint64) {
	outTxID := ob.GetTxID(nonce)
	ob.mu.Lock()
	ob.broadcastedTx[outTxID] = txHash
	ob.mu.Unlock()

	broadcastEntry := clienttypes.ToTransactionHashSQLType(txHash, outTxID)
	if err := ob.db.Create(&broadcastEntry).Error; err != nil {
		ob.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: error saving broadcasted tx")
	}
}

func (ob *BitcoinChainClient) observeOutTx() {
	ticker := time.NewTicker(time.Duration(ob.GetChainConfig().CoreParams.OutTxTicker) * time.Second)
	for {
		select {
		case <-ticker.C:
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain)
			if err != nil {
				ob.logger.ObserveOutTx.Error().Err(err).Msg("error GetAllOutTxTrackerByChain")
				continue
			}
			for _, tracker := range trackers {
				outTxID := ob.GetTxID(tracker.Nonce)
				ob.logger.ObserveOutTx.Info().Msgf("tracker outTxID: %s", outTxID)
				for _, txHash := range tracker.HashList {
					hash, err := chainhash.NewHashFromStr(txHash.TxHash)
					if err != nil {
						ob.logger.ObserveOutTx.Error().Err(err).Msg("error NewHashFromStr")
						continue
					}
					// The Bitcoin node has to be configured to watch TSS address
					getTxResult, err := ob.rpcClient.GetTransaction(hash)
					if err != nil {
						ob.logger.ObserveOutTx.Warn().Err(err).Msgf("error GetTransaction: %s", txHash.TxHash)
						continue
					}
					// Check TSS outTx
					err = ob.checkTssOutTxResult(hash, getTxResult)
					if err != nil {
						ob.logger.ObserveOutTx.Warn().Err(err).Msgf("error checkTssOutTxResult: %s", txHash.TxHash)
						continue
					}
					ob.logger.ObserveOutTx.Info().Msgf("outTx %s has passed checkTssOutTxResult", txHash.TxHash)
					if getTxResult.Confirmations >= 0 {
						ob.mu.Lock()
						ob.includedTx[outTxID] = *getTxResult
						ob.mu.Unlock()
					}
				}
			}
		case <-ob.stop:
			ob.logger.ObserveOutTx.Info().Msg("observeOutTx stopped")
			return
		}
	}
}

// Basic TSS outTX checks:
//   - locate the raw tx and find the Vin
//   - check if all inputs are segwit && TSS inputs
//
// Returns: true if outTx passes basic checks.
func (ob *BitcoinChainClient) checkTssOutTxResult(hash *chainhash.Hash, res *btcjson.GetTransactionResult) error {
	if res.Confirmations == 0 {
		rawtx, err := ob.rpcClient.GetRawTransactionVerbose(hash) // for pending tx, we query the raw tx
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: error GetRawTransactionVerbose %s", res.TxID)
		}
		if !ob.isValidTSSVin(rawtx.Vin) {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid outTx with non-TSS vin %s", res.TxID)
		}
	} else if res.Confirmations > 0 {
		blkHash, err := chainhash.NewHashFromStr(res.BlockHash)
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: error NewHashFromStr %s", res.BlockHash)
		}
		block, err := ob.rpcClient.GetBlockVerboseTx(blkHash) // for confirmed tx, we query the block
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: error GetBlockVerboseTx %s", res.BlockHash)
		}
		if res.BlockIndex < 0 || res.BlockIndex >= int64(len(block.Tx)) {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid outTx with invalid block index, TxID %s, BlockIndex %d", res.TxID, res.BlockIndex)
		}
		tx := block.Tx[res.BlockIndex]
		if !ob.isValidTSSVin(tx.Vin) {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid outTx with non-TSS vin %s", res.TxID)
		}
	}
	return nil // ignore res.Confirmations < 0 (meaning not included)
}

// Returns true only if all inputs are TSS vins
func (ob *BitcoinChainClient) isValidTSSVin(vins []btcjson.Vin) bool {
	if len(vins) == 0 {
		return false
	}
	pubKeyTss := hex.EncodeToString(ob.Tss.PubKeyCompressedBytes())
	for _, vin := range vins {
		// The length of the Witness should be always 2 for P2WPKH SegWit inputs.
		if len(vin.Witness) != 2 {
			return false
		}
		if vin.Witness[1] != pubKeyTss {
			return false
		}
	}
	return true
}

func (ob *BitcoinChainClient) BuildBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.TransactionHashSQLType
	if err := ob.db.Find(&broadcastedTransactions).Error; err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, entry := range broadcastedTransactions {
		ob.broadcastedTx[entry.Key] = entry.Hash
	}
	return nil
}

func (ob *BitcoinChainClient) LoadLastBlock() error {
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}

	//Load persisted block number
	var lastBlockNum clienttypes.LastBlockSQLType
	if err := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumID).Error; err != nil {
		ob.logger.ChainLogger.Info().Msg("LastBlockNum not found in DB, scan from latest")
		ob.SetLastBlockHeight(bn)
	} else {
		ob.SetLastBlockHeight(lastBlockNum.Num)

		//If persisted block number is too low, use the latest height
		if (bn - lastBlockNum.Num) > maxHeightDiff {
			ob.logger.ChainLogger.Info().Msgf("LastBlockNum too low: %d, scan from latest", lastBlockNum.Num)
			ob.SetLastBlockHeight(bn)
		}
	}

	if ob.chain.ChainId == 18444 { // bitcoin regtest: start from block 100
		ob.SetLastBlockHeight(100)
	}
	ob.logger.ChainLogger.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.lastBlock)

	return nil
}

func (ob *BitcoinChainClient) loadDB(dbpath string) error {
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		err := os.MkdirAll(dbpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	path := fmt.Sprintf("%s/btc_chain_client", dbpath)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic("failed to connect database")
	}
	ob.db = db

	err = db.AutoMigrate(&clienttypes.TransactionResultSQLType{},
		&clienttypes.TransactionHashSQLType{},
		&clienttypes.LastBlockSQLType{})
	if err != nil {
		return err
	}

	//Load last block
	err = ob.LoadLastBlock()
	if err != nil {
		return err
	}

	//Load broadcasted transactions
	err = ob.BuildBroadcastedTxMap()

	return err
}

func (ob *BitcoinChainClient) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.BTCAddress()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

// A very special value to mark current nonce in UTXO
func NonceMarkAmount(nonce uint64) int64 {
	return int64(nonce) + dustOffset // +2000 to avoid being a dust rejection
}
