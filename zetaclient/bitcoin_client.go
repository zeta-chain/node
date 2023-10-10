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

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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

	chain            common.Chain
	rpcClient        *rpcclient.Client
	zetaClient       *ZetaCoreBridge
	Tss              TSSSigner
	lastBlock        int64
	lastBlockScanned int64
	BlockTime        uint64 // block time in seconds

	mu                *sync.Mutex                             // lock for pending nonce, all the maps, utxos and core params
	pendingNonce      uint64                                  // the artificial pending nonce (next nonce to process) for outTx
	includedTxHashes  map[string]uint64                       // key: tx hash
	includedTxResults map[string]btcjson.GetTransactionResult // key: chain-tss-nonce
	broadcastedTx     map[string]string                       // key: chain-tss-nonce, value: outTx hash
	utxos             []btcjson.ListUnspentResult
	params            observertypes.CoreParams

	db     *gorm.DB
	stop   chan struct{}
	logger BTCLog
	ts     *TelemetryServer
}

const (
	minConfirmations = 0
	maxHeightDiff    = 10000
	dustOffset       = 2000
)

func (ob *BitcoinChainClient) SetCoreParams(params observertypes.CoreParams) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.params = params
}

func (ob *BitcoinChainClient) GetCoreParams() observertypes.CoreParams {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.params
}

// Return configuration based on supplied target chain
func NewBitcoinClient(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics, logger zerolog.Logger, btcCfg config.BTCConfig, ts *TelemetryServer) (*BitcoinChainClient, error) {
	ob := BitcoinChainClient{
		ChainMetrics: NewChainMetrics(chain.ChainName.String(), metrics),
		ts:           ts,
	}
	ob.stop = make(chan struct{})
	ob.chain = chain
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
	ob.includedTxHashes = make(map[string]uint64)
	ob.includedTxResults = make(map[string]btcjson.GetTransactionResult)
	ob.broadcastedTx = make(map[string]string)
	ob.params = btcCfg.CoreParams

	// initialize the Client
	ob.logger.ChainLogger.Info().Msgf("Chain %s endpoint %s", ob.chain.String(), btcCfg.RPCHost)
	connCfg := &rpcclient.ConnConfig{
		Host:         btcCfg.RPCHost,
		User:         btcCfg.RPCUsername,
		Pass:         btcCfg.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       btcCfg.RPCParams,
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

	err = ob.RegisterPromGauge(metricsPkg.PendingTxs, "Number of pending transactions")
	if err != nil {
		return nil, err
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

func (ob *BitcoinChainClient) SetLastBlockHeightScanned(block int64) {
	if block < 0 {
		panic("lastBlockScanned is negative")
	}
	if block >= math2.MaxInt64 {
		panic("lastBlockScanned is too large")
	}
	atomic.StoreInt64(&ob.lastBlockScanned, block)
	ob.ts.SetLastScannedBlockNumber((ob.chain.ChainId), (block))
}

func (ob *BitcoinChainClient) GetLastBlockHeightScanned() int64 {
	height := atomic.LoadInt64(&ob.lastBlockScanned)
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	if height >= math2.MaxInt64 {
		panic("lastBlockScanned is too large")
	}
	return height
}

func (ob *BitcoinChainClient) GetPendingNonce() uint64 {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.pendingNonce
}

// GetBaseGasPrice ...
// TODO: implement
// https://github.com/zeta-chain/node/issues/868
func (ob *BitcoinChainClient) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

func (ob *BitcoinChainClient) WatchInTx() {
	ticker := NewDynamicTicker("Bitcoin_WatchInTx", ob.GetCoreParams().InTxTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.observeInTx()
			if err != nil {
				ob.logger.WatchInTx.Error().Err(err).Msg("error observing in tx")
			}
			ticker.UpdateInterval(ob.GetCoreParams().InTxTicker, ob.logger.WatchInTx)
		case <-ob.stop:
			ob.logger.WatchInTx.Info().Msg("WatchInTx stopped")
			return
		}
	}
}

// TODO
func (ob *BitcoinChainClient) observeInTx() error {
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("error getting block count: %s", err)
	}
	if cnt < 0 || cnt >= math2.MaxInt64 {
		return fmt.Errorf("block count is out of range: %d", cnt)
	}

	// "confirmed" current block number
	// #nosec G701 always in range
	confirmedBlockNum := cnt - int64(ob.GetCoreParams().ConfirmationCount)
	if confirmedBlockNum < 0 || confirmedBlockNum > math2.MaxInt64 {
		return fmt.Errorf("skipping observer , confirmedBlockNum is negative or too large ")
	}
	ob.SetLastBlockHeight(confirmedBlockNum)

	flags, err := ob.zetaClient.GetCrosschainFlags()
	if err != nil {
		return err
	}
	if !flags.IsInboundEnabled {
		return errors.New("inbound TXS / Send has been disabled by the protocol")
	}

	lastBN := ob.GetLastBlockHeightScanned()

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
		// #nosec G701 always positive
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
		ob.SetLastBlockHeightScanned(bn)
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeightScanned())).Error; err != nil {
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

// returns isIncluded(or inMempool), isConfirmed, Error
func (ob *BitcoinChainClient) IsSendOutTxProcessed(sendHash string, nonce uint64, _ common.CoinType, logger zerolog.Logger) (bool, bool, error) {
	outTxID := ob.GetTxID(nonce)
	logger.Info().Msgf("IsSendOutTxProcessed %s", outTxID)

	ob.mu.Lock()
	txnHash, broadcasted := ob.broadcastedTx[outTxID]
	res, included := ob.includedTxResults[outTxID]
	ob.mu.Unlock()

	// Get original cctx parameters
	params, err := ob.GetPendingCctxParams(nonce)
	if err != nil {
		ob.logger.ObserveOutTx.Info().Msgf("IsSendOutTxProcessed: can't find pending cctx for nonce %d", nonce)
		return false, false, err
	}

	if !included {
		if !broadcasted {
			return false, false, nil
		}
		// If the broadcasted outTx is nonce 0, just wait for inclusion and don't schedule more keysign
		// Schedule more than one keysign for nonce 0 can lead to duplicate payments.
		// One purpose of nonce mark UTXO is to avoid duplicate payment based on the fact that Bitcoin
		// prevents double spending of same UTXO. However, for nonce 0, we don't have a prior nonce (e.g., -1)
		// for the signer to check against when making the payment. Signer treats nonce 0 as a special case in downstream code.
		if nonce == 0 {
			return true, false, nil
		}

		// Try including this outTx broadcasted by myself
		inMempool, err := ob.checkNSaveIncludedTx(txnHash, params)
		if err != nil {
			ob.logger.ObserveOutTx.Error().Err(err).Msg("IsSendOutTxProcessed: checkNSaveIncludedTx failed")
			return false, false, nil
		}
		if inMempool { // to avoid unnecessary Tss keysign
			ob.logger.ObserveOutTx.Info().Msgf("IsSendOutTxProcessed: outTx %s is still in mempool", outTxID)
			return true, false, nil
		}

		// Get tx result again in case it is just included
		ob.mu.Lock()
		res, included = ob.includedTxResults[outTxID]
		ob.mu.Unlock()
		if !included {
			return false, false, nil
		}
		ob.logger.ObserveOutTx.Info().Msgf("IsSendOutTxProcessed: checkNSaveIncludedTx succeeded for outTx %s", outTxID)
	}

	// It's safe to use cctx's amount to post confirmation because it has already been verified in observeOutTx()
	amountInSat := params.Amount.BigInt()
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		return true, false, nil
	}

	logger.Debug().Msgf("Bitcoin outTx confirmed: txid %s, amount %s\n", res.TxID, amountInSat.String())
	zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
		sendHash,
		res.TxID,
		// #nosec G701 always positive
		uint64(res.BlockIndex),
		0,   // gas used not used with Bitcoin
		nil, // gas price not used with Bitcoin
		0,   // gas limit not used with Bitcoin
		amountInSat,
		common.ReceiveStatus_Success,
		ob.chain,
		nonce,
		common.CoinType_Gas,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("error posting to zeta core")
	} else {
		logger.Info().Msgf("Bitcoin outTx %s confirmed: PostReceiveConfirmation zeta tx: %s", res.TxID, zetaHash)
	}
	return true, true, nil
}

func (ob *BitcoinChainClient) WatchGasPrice() {
	ticker := NewDynamicTicker("Bitcoin_WatchGasPrice", ob.GetCoreParams().GasPriceTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.WatchGasPrice.Error().Err(err).Msg("PostGasPrice error on " + ob.chain.String())
			}
			ticker.UpdateInterval(ob.GetCoreParams().GasPriceTicker, ob.logger.WatchGasPrice)
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
		// #nosec G701 always in range
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
	// feerate from RPC is BTC/KB, convert it to satoshi/byte
	// FIXME: in zetacore the gaslimit(vsize in BTC) is 100 which is too low for a typical outbound tx
	// until we fix the gaslimit in BTC, we need to multiply the feerate by 20 to make sure the tx is confirmed
	gasPriceU64, _ := gasPrice.Mul(big.NewFloat(*feeResult.FeeRate), big.NewFloat(20*1e5)).Uint64()
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}
	// #nosec G701 always positive
	zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, gasPriceU64, "100", uint64(bn))
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("PostGasPrice:")
		return err
	}
	_ = zetaHash
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
	ticker := NewDynamicTicker("Bitcoin_WatchUTXOS", ob.GetCoreParams().WatchUtxoTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.FetchUTXOS()
			if err != nil {
				ob.logger.WatchUTXOS.Error().Err(err).Msg("error fetching btc utxos")
			}
			ticker.UpdateInterval(ob.GetCoreParams().WatchUtxoTicker, ob.logger.WatchUTXOS)
		case <-ob.stop:
			ob.logger.WatchUTXOS.Info().Msg("WatchUTXOS stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) FetchUTXOS() error {
	defer func() {
		if err := recover(); err != nil {
			ob.logger.WatchUTXOS.Error().Msgf("BTC fetchUTXOS: caught panic error: %v", err)
		}
	}()

	// This is useful when a zetaclient's pending nonce lagged behind for whatever reason.
	ob.refreshPendingNonce()

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

	// fetching all TSS utxos takes 160ms
	utxos, err := ob.rpcClient.ListUnspentMinMaxAddresses(0, maxConfirmations, addresses)
	if err != nil {
		return err
	}
	//ob.logger.WatchUTXOS.Debug().Msgf("btc: fetched %d utxos in confirmation range [0, %d]", len(unspents), maxConfirmations)

	// rigid sort to make utxo list deterministic
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Amount == utxos[j].Amount {
			if utxos[i].TxID == utxos[j].TxID {
				return utxos[i].Vout < utxos[j].Vout
			}
			return utxos[i].TxID < utxos[j].TxID
		}
		return utxos[i].Amount < utxos[j].Amount
	})

	ob.mu.Lock()
	ob.ts.SetNumberOfUTXOs(len(utxos))
	ob.utxos = utxos
	ob.mu.Unlock()
	return nil
}

// refreshPendingNonce tries increasing the artificial pending nonce of outTx (if lagged behind).
// There could be many (unpredictable) reasons for a pending nonce lagging behind, for example:
// 1. The zetaclient gets restarted.
// 2. The tracker is missing in zetacore.
func (ob *BitcoinChainClient) refreshPendingNonce() {
	// get pending nonces from zetacore
	p, err := ob.zetaClient.GetPendingNoncesByChain(ob.chain.ChainId)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("refreshPendingNonce: error getting pending nonces")
	}

	// increase pending nonce if lagged behind
	ob.mu.Lock()
	pendingNonce := ob.pendingNonce
	ob.mu.Unlock()

	// #nosec G701 always positive
	nonceLow := uint64(p.NonceLow)

	if nonceLow > 0 && nonceLow >= pendingNonce {
		// get the last included outTx hash
		txid, err := ob.getOutTxidByNonce(nonceLow-1, false)
		if err != nil {
			ob.logger.ChainLogger.Error().Err(err).Msg("refreshPendingNonce: error getting last outTx txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.mu.Lock()
		defer ob.mu.Unlock()
		ob.pendingNonce = nonceLow
		ob.logger.ChainLogger.Info().Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", ob.pendingNonce, txid)
	}
}

// Set `test` flag to true in unit test to bypass query to zetacore
func (ob *BitcoinChainClient) getOutTxidByNonce(nonce uint64, test bool) (string, error) {
	ob.mu.Lock()
	res, included := ob.includedTxResults[ob.GetTxID(nonce)]
	ob.mu.Unlock()

	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetacore based on majority vote.
	if included {
		return res.TxID, nil
	}
	if !test { // if not unit test, get cctx from zetacore
		send, err := ob.zetaClient.GetCctxByNonce(ob.chain.ChainId, nonce)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting cctx for nonce %d", nonce)
		}
		txid := send.GetCurrentOutTxParam().OutboundTxHash
		if txid == "" {
			return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
		}
		// make sure it's a real Bitcoin txid
		_, getTxResult, err := ob.GetTxResultByHash(txid)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting outTx result for nonce %d hash %s", nonce, txid)
		}
		if getTxResult.Confirmations <= 0 { // just a double check
			return "", fmt.Errorf("getOutTxidByNonce: outTx txid %s for nonce %d is not included", txid, nonce)
		}
		return txid, nil
	}
	return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
}

func (ob *BitcoinChainClient) findNonceMarkUTXO(nonce uint64, txid string) (int, error) {
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	amount := NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := getSatoshis(utxo.Amount)
		if err != nil {
			ob.logger.ObserveOutTx.Error().Err(err).Msgf("findNonceMarkUTXO: error getting satoshis for utxo %v", utxo)
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == txid {
			ob.logger.ObserveOutTx.Info().Msgf("findNonceMarkUTXO: found nonce-mark utxo with txid %s, amount %d satoshi", utxo.TxID, sats)
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
//   - test: true for unit test only.
//
// Returns: a sublist (includes previous nonce-mark) of UTXOs or an error if the qulifying sublist cannot be found.
func (ob *BitcoinChainClient) SelectUTXOs(amount float64, utxoCap uint8, nonce uint64, test bool) ([]btcjson.ListUnspentResult, float64, error) {
	idx := -1
	if nonce == 0 {
		// for nonce = 0; make exception; no need to include nonce-mark utxo
		ob.mu.Lock()
		defer ob.mu.Unlock()
	} else {
		// for nonce > 0; we proceed only when we see the nonce-mark utxo
		preTxid, err := ob.getOutTxidByNonce(nonce-1, test)
		if err != nil {
			return nil, 0, err
		}
		ob.mu.Lock()
		defer ob.mu.Unlock()
		idx, err = ob.findNonceMarkUTXO(nonce-1, preTxid)
		if err != nil {
			return nil, 0, err
		}
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
	results := make([]btcjson.ListUnspentResult, right-left)
	copy(results, ob.utxos[left:right])

	// Note: always put nonce-mark as 1st input
	if idx >= 0 {
		if idx < left || idx >= right {
			total += ob.utxos[idx].Amount
			results = append([]btcjson.ListUnspentResult{ob.utxos[idx]}, results...)
		} else { // move nonce-mark to left
			for i := idx - left; i > 0; i-- {
				results[i], results[i-1] = results[i-1], results[i]
			}
		}
	}
	if total < amount {
		return nil, 0, fmt.Errorf("SelectUTXOs: not enough btc in reserve - available : %v , tx amount : %v", total, amount)
	}
	return results, total, nil
}

// Save successfully broadcasted transaction
func (ob *BitcoinChainClient) SaveBroadcastedTx(txHash string, nonce uint64) {
	outTxID := ob.GetTxID(nonce)
	ob.mu.Lock()
	ob.broadcastedTx[outTxID] = txHash
	ob.mu.Unlock()

	broadcastEntry := clienttypes.ToOutTxHashSQLType(txHash, outTxID)
	if err := ob.db.Save(&broadcastEntry).Error; err != nil {
		ob.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: error saving broadcasted tx")
	}
}

func (ob *BitcoinChainClient) GetPendingCctxParams(nonce uint64) (types.OutboundTxParams, error) {
	send, err := ob.zetaClient.GetCctxByNonce(ob.chain.ChainId, nonce)
	if err != nil {
		return types.OutboundTxParams{}, err
	}
	if send.GetCurrentOutTxParam() == nil { // never happen
		return types.OutboundTxParams{}, fmt.Errorf("GetPendingCctx: nil outbound tx params")
	}
	if send.CctxStatus.Status == types.CctxStatus_PendingOutbound || send.CctxStatus.Status == types.CctxStatus_PendingRevert {
		return *send.GetCurrentOutTxParam(), nil
	}
	return types.OutboundTxParams{}, fmt.Errorf("GetPendingCctx: not a pending cctx")
}

func (ob *BitcoinChainClient) observeOutTx() {
	ticker := NewDynamicTicker("Bitcoin_observeOutTx", ob.GetCoreParams().OutTxTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain, Ascending)
			if err != nil {
				ob.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: error GetAllOutTxTrackerByChain")
				continue
			}
			for _, tracker := range trackers {
				// get original cctx parameters
				outTxID := ob.GetTxID(tracker.Nonce)
				params, err := ob.GetPendingCctxParams(tracker.Nonce)
				if err != nil {
					ob.logger.ObserveOutTx.Info().Err(err).Msgf("observeOutTx: can't find pending cctx for nonce %d", tracker.Nonce)
					break
				}
				if tracker.Nonce != params.OutboundTxTssNonce { // Tanmay: it doesn't hurt to check
					ob.logger.ObserveOutTx.Error().Msgf("observeOutTx: tracker nonce %d not match cctx nonce %d", tracker.Nonce, params.OutboundTxTssNonce)
					break
				}
				if len(tracker.HashList) > 1 {
					ob.logger.ObserveOutTx.Warn().Msgf("observeOutTx: oops, outTxID %s got multiple (%d) outTx hashes", outTxID, len(tracker.HashList))
				}
				// verify outTx hashes
				for _, txHash := range tracker.HashList {
					_, err := ob.checkNSaveIncludedTx(txHash.TxHash, params)
					if err != nil {
						ob.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: checkNSaveIncludedTx failed")
					}
				}
			}
			ticker.UpdateInterval(ob.GetCoreParams().OutTxTicker, ob.logger.ObserveOutTx)
		case <-ob.stop:
			ob.logger.ObserveOutTx.Info().Msg("observeOutTx stopped")
			return
		}
	}
}

// checkNSaveIncludedTx either includes a new outTx or update an existing outTx result.
// Returns inMempool, error
func (ob *BitcoinChainClient) checkNSaveIncludedTx(txHash string, params types.OutboundTxParams) (bool, error) {
	outTxID := ob.GetTxID(params.OutboundTxTssNonce)
	hash, getTxResult, err := ob.GetTxResultByHash(txHash)
	if err != nil {
		return false, errors.Wrapf(err, "checkNSaveIncludedTx: error GetTxResultByHash: %s", txHash)
	}
	if getTxResult.Confirmations >= 0 { // check included tx only
		err = ob.checkTssOutTxResult(hash, getTxResult, params, params.OutboundTxTssNonce)
		if err != nil {
			return false, errors.Wrapf(err, "checkNSaveIncludedTx: error verify bitcoin outTx %s outTxID %s", txHash, outTxID)
		}

		ob.mu.Lock()
		defer ob.mu.Unlock()
		nonce, foundHash := ob.includedTxHashes[txHash]
		res, foundRes := ob.includedTxResults[outTxID]

		// include new outTx and enforce rigid 1-to-1 mapping: outTxID(nonce) <===> txHash
		if !foundHash && !foundRes {
			ob.includedTxHashes[txHash] = params.OutboundTxTssNonce
			ob.includedTxResults[outTxID] = *getTxResult
			if params.OutboundTxTssNonce >= ob.pendingNonce { // try increasing pending nonce on every newly included outTx
				ob.pendingNonce = params.OutboundTxTssNonce
			}
			ob.logger.ObserveOutTx.Info().Msgf("checkNSaveIncludedTx: included new bitcoin outTx %s outTxID %s", txHash, outTxID)
		}
		// update saved tx result as confirmations may increase
		if foundHash && foundRes {
			ob.includedTxResults[outTxID] = *getTxResult
			if getTxResult.Confirmations > res.Confirmations {
				ob.logger.ObserveOutTx.Info().Msgf("checkNSaveIncludedTx: bitcoin outTx %s got confirmations %d", txHash, getTxResult.Confirmations)
			}
		}
		if !foundHash && foundRes { // be alert for duplicate payment!!! As we got a new hash paying same cctx. It might happen (e.g. majority of signers get crupted)
			ob.logger.ObserveOutTx.Error().Msgf("checkNSaveIncludedTx: duplicate payment by bitcoin outTx %s outTxID %s, prior result %v, current result %v", txHash, outTxID, res, *getTxResult)
		}
		if foundHash && !foundRes {
			ob.logger.ObserveOutTx.Error().Msgf("checkNSaveIncludedTx: unreachable code path! outTx %s outTxID %s, prior nonce %d, current nonce %d", txHash, outTxID, nonce, params.OutboundTxTssNonce)
		}
		return false, nil
	}
	return true, nil // in mempool
}

// Basic TSS outTX checks:
//   - should be able to query the raw tx
//   - check if all inputs are segwit && TSS inputs
//
// Returns: true if outTx passes basic checks.
func (ob *BitcoinChainClient) checkTssOutTxResult(hash *chainhash.Hash, res *btcjson.GetTransactionResult, params types.OutboundTxParams, nonce uint64) error {
	rawResult, err := ob.getRawTxResult(hash, res)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: error GetRawTxResultByHash %s", hash.String())
	}
	err = ob.checkTSSVin(rawResult.Vin, nonce)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vin in outTx %s nonce %d", hash, nonce)
	}
	err = ob.checkTSSVout(rawResult.Vout, params, nonce)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vout in outTx %s nonce %d", hash, nonce)
	}
	return nil
}

func (ob *BitcoinChainClient) GetTxResultByHash(txID string) (*chainhash.Hash, *btcjson.GetTransactionResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTxResultByHash: error NewHashFromStr: %s", txID)
	}

	// The Bitcoin node has to be configured to watch TSS address
	txResult, err := ob.rpcClient.GetTransaction(hash)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetOutTxByTxHash: error GetTransaction %s", hash.String())
	}
	return hash, txResult, nil
}

func (ob *BitcoinChainClient) getRawTxResult(hash *chainhash.Hash, res *btcjson.GetTransactionResult) (btcjson.TxRawResult, error) {
	if res.Confirmations == 0 { // for pending tx, we query the raw tx directly
		rawResult, err := ob.rpcClient.GetRawTransactionVerbose(hash) // for pending tx, we query the raw tx
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error GetRawTransactionVerbose %s", res.TxID)
		}
		return *rawResult, nil
	} else if res.Confirmations > 0 { // for confirmed tx, we query the block
		blkHash, err := chainhash.NewHashFromStr(res.BlockHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error NewHashFromStr for block hash %s", res.BlockHash)
		}
		block, err := ob.rpcClient.GetBlockVerboseTx(blkHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error GetBlockVerboseTx %s", res.BlockHash)
		}
		if res.BlockIndex < 0 || res.BlockIndex >= int64(len(block.Tx)) {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: invalid outTx with invalid block index, TxID %s, BlockIndex %d", res.TxID, res.BlockIndex)
		}
		return block.Tx[res.BlockIndex], nil
	} else { // res.Confirmations < 0 (meaning not included)
		return btcjson.TxRawResult{}, fmt.Errorf("getRawTxResult: tx %s not included yet", hash)
	}
}

// Vin is valid if:
//   - The first input is the nonce-mark
//   - All inputs are from TSS address
func (ob *BitcoinChainClient) checkTSSVin(vins []btcjson.Vin, nonce uint64) error {
	// vins: [nonce-mark, UTXO1, UTXO2, ...]
	if len(vins) <= 1 {
		return fmt.Errorf("checkTSSVin: len(vins) <= 1")
	}
	pubKeyTss := hex.EncodeToString(ob.Tss.PubKeyCompressedBytes())
	for i, vin := range vins {
		// The length of the Witness should be always 2 for SegWit inputs.
		if len(vin.Witness) != 2 {
			return fmt.Errorf("checkTSSVin: expected 2 witness items, got %d", len(vin.Witness))
		}
		if vin.Witness[1] != pubKeyTss {
			return fmt.Errorf("checkTSSVin: witness pubkey %s not match TSS pubkey %s", vin.Witness[1], pubKeyTss)
		}
		// 1st vin: nonce-mark MUST come from prior TSS outTx
		if nonce > 0 && i == 0 {
			preTxid, err := ob.getOutTxidByNonce(nonce-1, false)
			if err != nil {
				return fmt.Errorf("checkTSSVin: error findTxIDByNonce %d", nonce-1)
			}
			// nonce-mark MUST the 1st output that comes from prior TSS outTx
			if vin.Txid != preTxid || vin.Vout != 0 {
				return fmt.Errorf("checkTSSVin: invalid nonce-mark txid %s vout %d, expected txid %s vout 0", vin.Txid, vin.Vout, preTxid)
			}
		}
	}
	return nil
}

// Vout is valid if:
//   - The first output is the nonce-mark
//   - The second output is the correct payment to recipient
//   - The third output is the change to TSS (optional)
func (ob *BitcoinChainClient) checkTSSVout(vouts []btcjson.Vout, params types.OutboundTxParams, nonce uint64) error {
	// vouts: [nonce-mark, payment to recipient, change to TSS (optional)]
	if !(len(vouts) == 2 || len(vouts) == 3) {
		return fmt.Errorf("checkTSSVout: invalid number of vouts: %d", len(vouts))
	}

	tssAddress := ob.Tss.BTCAddress()
	for _, vout := range vouts {
		amount, err := getSatoshis(vout.Value)
		if err != nil {
			return errors.Wrap(err, "checkTSSVout: error getting satoshis")
		}
		// decode P2WPKH scriptPubKey
		scriptPubKey := vout.ScriptPubKey.Hex
		decodedScriptPubKey, err := hex.DecodeString(scriptPubKey)
		if err != nil {
			return errors.Wrapf(err, "checkTSSVout: error decoding scriptPubKey %s", scriptPubKey)
		}
		if len(decodedScriptPubKey) != 22 { // P2WPKH script
			return fmt.Errorf("checkTSSVout: unsupported scriptPubKey: %s", scriptPubKey)
		}
		witnessVersion := decodedScriptPubKey[0]
		witnessProgram := decodedScriptPubKey[2:]
		if witnessVersion != 0 {
			return fmt.Errorf("checkTSSVout: unsupported witness in scriptPubKey %s", scriptPubKey)
		}
		recvAddress, err := ob.chain.BTCAddressFromWitnessProgram(witnessProgram)
		if err != nil {
			return errors.Wrapf(err, "checkTSSVout: error getting receiver from witness program %s", witnessProgram)
		}

		// 1st vout: nonce-mark
		if vout.N == 0 {
			if recvAddress != tssAddress {
				return fmt.Errorf("checkTSSVout: nonce-mark address %s not match TSS address %s", recvAddress, tssAddress)
			}
			if amount != NonceMarkAmount(nonce) {
				return fmt.Errorf("checkTSSVout: nonce-mark amount %d not match nonce-mark amount %d", amount, NonceMarkAmount(nonce))
			}
		}
		// 2nd vout: payment to recipient
		if vout.N == 1 {
			if recvAddress != params.Receiver {
				return fmt.Errorf("checkTSSVout: output address %s not match params receiver %s", recvAddress, params.Receiver)
			}
			// #nosec G701 always positive
			if uint64(amount) != params.Amount.Uint64() {
				return fmt.Errorf("checkTSSVout: output amount %d not match params amount %d", amount, params.Amount)
			}
		}
		// 3rd vout: change to TSS (optional)
		if vout.N == 2 {
			if recvAddress != tssAddress {
				return fmt.Errorf("checkTSSVout: change address %s not match TSS address %s", recvAddress, tssAddress)
			}
		}
	}
	return nil
}

func (ob *BitcoinChainClient) BuildBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutTxHashSQLType
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
		ob.SetLastBlockHeightScanned(bn)
	} else {
		ob.SetLastBlockHeightScanned(lastBlockNum.Num)

		//If persisted block number is too low, use the latest height
		if (bn - lastBlockNum.Num) > maxHeightDiff {
			ob.logger.ChainLogger.Info().Msgf("LastBlockNum too low: %d, scan from latest", lastBlockNum.Num)
			ob.SetLastBlockHeightScanned(bn)
		}
	}

	if ob.chain.ChainId == 18444 { // bitcoin regtest: start from block 100
		ob.SetLastBlockHeightScanned(100)
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
		&clienttypes.OutTxHashSQLType{},
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
	// #nosec G701 always in range
	return int64(nonce) + dustOffset // +2000 to avoid being a dust rejection
}
