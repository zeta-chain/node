package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// DynamicDepositorFeeHeight contains the starting height (Bitcoin mainnet) from which dynamic depositor fee will take effect
	DynamicDepositorFeeHeight = 834500

	// maxHeightDiff contains the max height diff in case the last block is too old when the observer starts
	maxHeightDiff = 10000

	// btcBlocksPerDay represents Bitcoin blocks per days for LRU block cache size
	btcBlocksPerDay = 144

	// bigValueSats contains the threshold to determine a big value in Bitcoin represents 2 BTC
	bigValueSats = 200000000

	// bigValueConfirmationCount represents the number of confirmation necessary for bigger values: 6 confirmations
	bigValueConfirmationCount = 6
)

var _ interfaces.ChainClient = &Client{}

// Logger contains list of loggers used by Bitcoin chain client
// TODO: Merge this logger with the one in evm
// https://github.com/zeta-chain/node/issues/2022
type Logger struct {
	// Chain is the parent logger for the chain
	Chain zerolog.Logger

	// InTx is the logger for incoming transactions
	InTx zerolog.Logger // The logger for incoming transactions

	// OutTx is the logger for outgoing transactions
	OutTx zerolog.Logger // The logger for outgoing transactions

	// UTXOS is the logger for UTXOs management
	UTXOS zerolog.Logger // The logger for UTXOs management

	// GasPrice is the logger for gas price
	GasPrice zerolog.Logger // The logger for gas price

	// Compliance is the logger for compliance checks
	Compliance zerolog.Logger // The logger for compliance checks
}

// BTCInTxEvent represents an incoming transaction event
type BTCInTxEvent struct {
	// FromAddress is the first input address
	FromAddress string

	// ToAddress is the TSS address
	ToAddress string

	// Value is the amount of BTC
	Value float64

	MemoBytes   []byte
	BlockNumber uint64
	TxHash      string
}

// BTCOutTxEvent contains bitcoin block and the header
type BTCBlockNHeader struct {
	Header *wire.BlockHeader
	Block  *btcjson.GetBlockVerboseTxResult
}

// Client represents a chain configuration for Bitcoin
// Filled with above constants depending on chain
type Client struct {
	BlockCache *lru.Cache

	// Mu is lock for all the maps, utxos and core params
	Mu *sync.Mutex

	Tss interfaces.TSSSigner

	chain            chains.Chain
	netParams        *chaincfg.Params
	rpcClient        interfaces.BTCRPCClient
	coreClient       interfaces.ZetaCoreClient
	lastBlock        int64
	lastBlockScanned int64
	pendingNonce     uint64
	utxos            []btcjson.ListUnspentResult
	params           observertypes.ChainParams
	coreContext      *context.ZetaCoreContext

	// includedTxHashes indexes included tx with tx hash
	includedTxHashes map[string]bool

	// includedTxResults indexes tx results with the outbound tx identifier
	includedTxResults map[string]*btcjson.GetTransactionResult

	// broadcastedTx indexes the outbound hash with the outbound tx identifier
	broadcastedTx map[string]string

	db     *gorm.DB
	stop   chan struct{}
	logger Logger
	ts     *metrics.TelemetryServer
}

func (ob *Client) WithZetaCoreClient(client *zetacore.Client) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.coreClient = client
}

func (ob *Client) WithLogger(logger zerolog.Logger) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.logger = Logger{
		Chain:    logger,
		InTx:     logger.With().Str("module", "WatchInTx").Logger(),
		OutTx:    logger.With().Str("module", "WatchOutTx").Logger(),
		UTXOS:    logger.With().Str("module", "WatchUTXOS").Logger(),
		GasPrice: logger.With().Str("module", "WatchGasPrice").Logger(),
	}
}

func (ob *Client) WithBtcClient(client *rpcclient.Client) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.rpcClient = client
}

func (ob *Client) WithChain(chain chains.Chain) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.chain = chain
}

func (ob *Client) SetChainParams(params observertypes.ChainParams) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.params = params
}

func (ob *Client) GetChainParams() observertypes.ChainParams {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.params
}

// NewClient returns a new Bitcoin chain client
func NewClient(
	appcontext *context.AppContext,
	chain chains.Chain,
	coreClient interfaces.ZetaCoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
	loggers clientcommon.ClientLogger,
	btcCfg config.BTCConfig,
	ts *metrics.TelemetryServer,
) (*Client, error) {
	// initialize the BTCChainClient
	ob := Client{
		ts: ts,
	}
	ob.stop = make(chan struct{})
	ob.chain = chain

	// get the bitcoin network params
	netParams, err := chains.BitcoinNetParamsFromChainID(ob.chain.ChainId)
	if err != nil {
		return nil, fmt.Errorf("error getting net params for chain %d: %s", ob.chain.ChainId, err)
	}
	ob.netParams = netParams

	ob.Mu = &sync.Mutex{}

	chainLogger := loggers.Std.With().Str("chain", chain.ChainName.String()).Logger()
	ob.logger = Logger{
		Chain:      chainLogger,
		InTx:       chainLogger.With().Str("module", "WatchInTx").Logger(),
		OutTx:      chainLogger.With().Str("module", "WatchOutTx").Logger(),
		UTXOS:      chainLogger.With().Str("module", "WatchUTXOS").Logger(),
		GasPrice:   chainLogger.With().Str("module", "WatchGasPrice").Logger(),
		Compliance: loggers.Compliance,
	}

	ob.coreClient = coreClient
	ob.Tss = tss
	ob.coreContext = appcontext.ZetaCoreContext()
	ob.includedTxHashes = make(map[string]bool)
	ob.includedTxResults = make(map[string]*btcjson.GetTransactionResult)
	ob.broadcastedTx = make(map[string]string)

	// set the Bitcoin chain params
	_, chainParams, found := appcontext.ZetaCoreContext().GetBTCChainParams()
	if !found {
		return nil, fmt.Errorf("btc chains params not initialized")
	}
	ob.params = *chainParams

	// create the RPC client
	ob.logger.Chain.Info().Msgf("Chain %s endpoint %s", ob.chain.String(), btcCfg.RPCHost)
	connCfg := &rpcclient.ConnConfig{
		Host:         btcCfg.RPCHost,
		User:         btcCfg.RPCUsername,
		Pass:         btcCfg.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       btcCfg.RPCParams,
	}
	rpcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating rpc client: %s", err)
	}

	// try connection
	ob.rpcClient = rpcClient
	err = rpcClient.Ping()
	if err != nil {
		return nil, fmt.Errorf("error ping the bitcoin server: %s", err)
	}

	ob.BlockCache, err = lru.New(btcBlocksPerDay)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("failed to create bitcoin block cache")
		return nil, err
	}

	// load btc chain client DB
	err = ob.loadDB(dbpath)
	if err != nil {
		return nil, err
	}

	return &ob, nil
}

// Start starts the Go routine to observe the Bitcoin chain
func (ob *Client) Start() {
	ob.logger.Chain.Info().Msgf("Bitcoin client is starting")
	go ob.WatchInTx()        // watch bitcoin chain for incoming txs and post votes to zetacore
	go ob.WatchOutTx()       // watch bitcoin chain for outgoing txs status
	go ob.WatchUTXOS()       // watch bitcoin chain for UTXOs owned by the TSS address
	go ob.WatchGasPrice()    // watch bitcoin chain for gas rate and post to zetacore
	go ob.WatchIntxTracker() // watch zetacore for bitcoin intx trackers
	go ob.WatchRPCStatus()   // watch the RPC status of the bitcoin chain
}

// WatchRPCStatus watches the RPC status of the Bitcoin chain
func (ob *Client) WatchRPCStatus() {
	ob.logger.Chain.Info().Msgf("RPCStatus is starting")
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}

			bn, err := ob.rpcClient.GetBlockCount()
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}

			hash, err := ob.rpcClient.GetBlockHash(bn)
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}

			header, err := ob.rpcClient.GetBlockHeader(hash)
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}

			blockTime := header.Timestamp
			elapsedSeconds := time.Since(blockTime).Seconds()
			if elapsedSeconds > 1200 {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}

			tssAddr := ob.Tss.BTCAddressWitnessPubkeyHash()
			res, err := ob.rpcClient.ListUnspentMinMaxAddresses(0, 1000000, []btcutil.Address{tssAddr})
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: can't list utxos of TSS address; wallet or loaded? TSS address is not imported? ")
				continue
			}

			if len(res) == 0 {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: TSS address has no utxos; TSS address is not imported? ")
				continue
			}

			ob.logger.Chain.Info().Msgf("[OK] RPC status check: latest block number %d, timestamp %s (%.fs ago), tss addr %s, #utxos: %d", bn, blockTime, elapsedSeconds, tssAddr, len(res))

		case <-ob.stop:
			return
		}
	}
}

func (ob *Client) Stop() {
	ob.logger.Chain.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop
	ob.logger.Chain.Info().Msgf("%s observer stopped", ob.chain.String())
}

func (ob *Client) SetLastBlockHeight(height int64) {
	if height < 0 {
		panic("lastBlock is negative")
	}
	atomic.StoreInt64(&ob.lastBlock, height)
}

func (ob *Client) GetLastBlockHeight() int64 {
	height := atomic.LoadInt64(&ob.lastBlock)
	if height < 0 {
		panic("lastBlock is negative")
	}
	return height
}

func (ob *Client) SetLastBlockHeightScanned(height int64) {
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	atomic.StoreInt64(&ob.lastBlockScanned, height)
	metrics.LastScannedBlockNumber.WithLabelValues(ob.chain.ChainName.String()).Set(float64(height))
}

func (ob *Client) GetLastBlockHeightScanned() int64 {
	height := atomic.LoadInt64(&ob.lastBlockScanned)
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	return height
}

func (ob *Client) GetPendingNonce() uint64 {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.pendingNonce
}

// GetBaseGasPrice ...
// TODO: implement
// https://github.com/zeta-chain/node/issues/868
func (ob *Client) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

// ConfirmationsThreshold returns number of required Bitcoin confirmations depending on sent BTC amount.
func (ob *Client) ConfirmationsThreshold(amount *big.Int) int64 {
	if amount.Cmp(big.NewInt(bigValueSats)) >= 0 {
		return bigValueConfirmationCount
	}
	if bigValueConfirmationCount < ob.GetChainParams().ConfirmationCount {
		return bigValueConfirmationCount
	}

	// #nosec G701 always in range
	return int64(ob.GetChainParams().ConfirmationCount)
}

// WatchGasPrice watches Bitcoin chain for gas rate and post to zetacore
func (ob *Client) WatchGasPrice() {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice()
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.chain.ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchGasPrice", ob.GetChainParams().GasPriceTicker)
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msg("error creating ticker")
		return
	}
	ob.logger.GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.chain.ChainId, ob.GetChainParams().GasPriceTicker)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.chain.ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.logger.GasPrice)
		case <-ob.stop:
			ob.logger.GasPrice.Info().Msgf("WatchGasPrice stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *Client) PostGasPrice() error {
	if ob.chain.ChainId == 18444 { //bitcoin regtest; hardcode here since this RPC is not available on regtest
		blockNumber, err := ob.rpcClient.GetBlockCount()
		if err != nil {
			return err
		}

		// #nosec G701 always in range
		_, err = ob.coreClient.PostGasPrice(ob.chain, 1, "100", uint64(blockNumber))
		if err != nil {
			ob.logger.GasPrice.Err(err).Msg("PostGasPrice:")
			return err
		}
		return nil
	}

	// EstimateSmartFee returns the fees per kilobyte (BTC/kb) targeting given block confirmation
	feeResult, err := ob.rpcClient.EstimateSmartFee(1, &btcjson.EstimateModeEconomical)
	if err != nil {
		return err
	}
	if feeResult.Errors != nil || feeResult.FeeRate == nil {
		return fmt.Errorf("error getting gas price: %s", feeResult.Errors)
	}
	if *feeResult.FeeRate > math.MaxInt64 {
		return fmt.Errorf("gas price is too large: %f", *feeResult.FeeRate)
	}
	feeRatePerByte := FeeRateToSatPerByte(*feeResult.FeeRate)

	blockNumber, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}

	// #nosec G701 always positive
	_, err = ob.coreClient.PostGasPrice(ob.chain, feeRatePerByte.Uint64(), "100", uint64(blockNumber))
	if err != nil {
		ob.logger.GasPrice.Err(err).Msg("PostGasPrice:")
		return err
	}

	return nil
}

// GetSenderAddressByVin get the sender address from the previous transaction
func GetSenderAddressByVin(rpcClient interfaces.BTCRPCClient, vin btcjson.Vin, net *chaincfg.Params) (string, error) {
	// query previous raw transaction by txid
	// GetTransaction requires reconfiguring the bitcoin node (txindex=1), so we use GetRawTransaction instead
	hash, err := chainhash.NewHashFromStr(vin.Txid)
	if err != nil {
		return "", err
	}

	tx, err := rpcClient.GetRawTransaction(hash)
	if err != nil {
		return "", errors.Wrapf(err, "error getting raw transaction %s", vin.Txid)
	}

	// #nosec G701 - always in range
	if len(tx.MsgTx().TxOut) <= int(vin.Vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vin.Vout, vin.Txid)
	}

	// decode sender address from previous pkScript
	pkScript := tx.MsgTx().TxOut[vin.Vout].PkScript
	scriptHex := hex.EncodeToString(pkScript)
	if IsPkScriptP2TR(pkScript) {
		return DecodeScriptP2TR(scriptHex, net)
	}
	if IsPkScriptP2WSH(pkScript) {
		return DecodeScriptP2WSH(scriptHex, net)
	}
	if IsPkScriptP2WPKH(pkScript) {
		return DecodeScriptP2WPKH(scriptHex, net)
	}
	if IsPkScriptP2SH(pkScript) {
		return DecodeScriptP2SH(scriptHex, net)
	}
	if IsPkScriptP2PKH(pkScript) {
		return DecodeScriptP2PKH(scriptHex, net)
	}

	// sender address not found, return nil and move on to the next tx
	return "", nil
}

// WatchUTXOS watches bitcoin chain for UTXOs owned by the TSS address
func (ob *Client) WatchUTXOS() {
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchUTXOS", ob.GetChainParams().WatchUtxoTicker)
	if err != nil {
		ob.logger.UTXOS.Error().Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.FetchUTXOS()
			if err != nil {
				ob.logger.UTXOS.Error().Err(err).Msg("error fetching btc utxos")
			}
			ticker.UpdateInterval(ob.GetChainParams().WatchUtxoTicker, ob.logger.UTXOS)
		case <-ob.stop:
			ob.logger.UTXOS.Info().Msgf("WatchUTXOS stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *Client) FetchUTXOS() error {
	defer func() {
		if err := recover(); err != nil {
			ob.logger.UTXOS.Error().Msgf("BTC fetchUTXOS: caught panic error: %v", err)
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

	// List all unspent UTXOs (160ms)
	tssAddr := ob.Tss.BTCAddress()
	address, err := chains.DecodeBtcAddress(tssAddr, ob.chain.ChainId)
	if err != nil {
		return fmt.Errorf("btc: error decoding wallet address (%s) : %s", tssAddr, err.Error())
	}
	utxos, err := ob.rpcClient.ListUnspentMinMaxAddresses(0, maxConfirmations, []btcutil.Address{address})
	if err != nil {
		return err
	}

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

	// filter UTXOs good to spend for next TSS transaction
	utxosFiltered := make([]btcjson.ListUnspentResult, 0)
	for _, utxo := range utxos {
		// UTXOs big enough to cover the cost of spending themselves
		if utxo.Amount < DefaultDepositorFee {
			continue
		}
		// we don't want to spend other people's unconfirmed UTXOs as they may not be safe to spend
		if utxo.Confirmations == 0 {
			if !ob.isTssTransaction(utxo.TxID) {
				continue
			}
		}
		utxosFiltered = append(utxosFiltered, utxo)
	}

	ob.Mu.Lock()
	metrics.NumberOfUTXO.Set(float64(len(utxosFiltered)))
	ob.utxos = utxosFiltered
	ob.Mu.Unlock()
	return nil
}

// SelectUTXOs selects a sublist of utxos to be used as inputs.
//
// Parameters:
//   - amount: The desired minimum total value of the selected UTXOs.
//   - utxos2Spend: The maximum number of UTXOs to spend.
//   - nonce: The nonce of the outbound transaction.
//   - consolidateRank: The rank below which UTXOs will be consolidated.
//   - test: true for unit test only.
//
// Returns:
//   - a sublist (includes previous nonce-mark) of UTXOs or an error if the qualifying sublist cannot be found.
//   - the total value of the selected UTXOs.
//   - the number of consolidated UTXOs.
//   - the total value of the consolidated UTXOs.
func (ob *Client) SelectUTXOs(
	amount float64,
	utxosToSpend uint16,
	nonce uint64,
	consolidateRank uint16,
	test bool,
) ([]btcjson.ListUnspentResult, float64, uint16, float64, error) {
	idx := -1
	if nonce == 0 {
		// for nonce = 0; make exception; no need to include nonce-mark utxo
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
	} else {
		// for nonce > 0; we proceed only when we see the nonce-mark utxo
		preTxid, err := ob.getOutTxidByNonce(nonce-1, test)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		idx, err = ob.findNonceMarkUTXO(nonce-1, preTxid)
		if err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// select smallest possible UTXOs to make payment
	total := 0.0
	left, right := 0, 0
	for total < amount && right < len(ob.utxos) {
		if utxosToSpend > 0 { // expand sublist
			total += ob.utxos[right].Amount
			right++
			utxosToSpend--
		} else { // pop the smallest utxo and append the current one
			total -= ob.utxos[left].Amount
			total += ob.utxos[right].Amount
			left++
			right++
		}
	}
	results := make([]btcjson.ListUnspentResult, right-left)
	copy(results, ob.utxos[left:right])

	// include nonce-mark as the 1st input
	if idx >= 0 { // for nonce > 0
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
		return nil, 0, 0, 0, fmt.Errorf("SelectUTXOs: not enough btc in reserve - available : %v , tx amount : %v", total, amount)
	}

	// consolidate biggest possible UTXOs to maximize consolidated value
	// consolidation happens only when there are more than (or equal to) consolidateRank (10) UTXOs
	utxoRank, consolidatedUtxo, consolidatedValue := uint16(0), uint16(0), 0.0
	for i := len(ob.utxos) - 1; i >= 0 && utxosToSpend > 0; i-- { // iterate over UTXOs big-to-small
		if i != idx && (i < left || i >= right) { // exclude nonce-mark and already selected UTXOs
			utxoRank++
			if utxoRank >= consolidateRank { // consolication starts from the 10-ranked UTXO based on value
				utxosToSpend--
				consolidatedUtxo++
				total += ob.utxos[i].Amount
				consolidatedValue += ob.utxos[i].Amount
				results = append(results, ob.utxos[i])
			}
		}
	}

	return results, total, consolidatedUtxo, consolidatedValue, nil
}

// SaveBroadcastedTx saves successfully broadcasted transaction
func (ob *Client) SaveBroadcastedTx(txHash string, nonce uint64) {
	outTxID := ob.GetTxID(nonce)
	ob.Mu.Lock()
	ob.broadcastedTx[outTxID] = txHash
	ob.Mu.Unlock()

	broadcastEntry := clienttypes.ToOutTxHashSQLType(txHash, outTxID)
	if err := ob.db.Save(&broadcastEntry).Error; err != nil {
		ob.logger.OutTx.Error().Err(err).Msgf("SaveBroadcastedTx: error saving broadcasted txHash %s for outTx %s", txHash, outTxID)
	}
	ob.logger.OutTx.Info().Msgf("SaveBroadcastedTx: saved broadcasted txHash %s for outTx %s", txHash, outTxID)
}

// GetTxResultByHash gets the transaction result by hash
func GetTxResultByHash(rpcClient interfaces.BTCRPCClient, txID string) (*chainhash.Hash, *btcjson.GetTransactionResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTxResultByHash: error NewHashFromStr: %s", txID)
	}

	// The Bitcoin node has to be configured to watch TSS address
	txResult, err := rpcClient.GetTransaction(hash)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetOutTxByTxHash: error GetTransaction %s", hash.String())
	}
	return hash, txResult, nil
}

// GetRawTxResult gets the raw tx result
func GetRawTxResult(rpcClient interfaces.BTCRPCClient, hash *chainhash.Hash, res *btcjson.GetTransactionResult) (btcjson.TxRawResult, error) {
	if res.Confirmations == 0 { // for pending tx, we query the raw tx directly
		rawResult, err := rpcClient.GetRawTransactionVerbose(hash) // for pending tx, we query the raw tx
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error GetRawTransactionVerbose %s", res.TxID)
		}
		return *rawResult, nil
	} else if res.Confirmations > 0 { // for confirmed tx, we query the block
		blkHash, err := chainhash.NewHashFromStr(res.BlockHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error NewHashFromStr for block hash %s", res.BlockHash)
		}
		block, err := rpcClient.GetBlockVerboseTx(blkHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error GetBlockVerboseTx %s", res.BlockHash)
		}
		if res.BlockIndex < 0 || res.BlockIndex >= int64(len(block.Tx)) {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: invalid outTx with invalid block index, TxID %s, BlockIndex %d", res.TxID, res.BlockIndex)
		}
		return block.Tx[res.BlockIndex], nil
	}

	// res.Confirmations < 0 (meaning not included)
	return btcjson.TxRawResult{}, fmt.Errorf("getRawTxResult: tx %s not included yet", hash)
}

func (ob *Client) BuildBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutTxHashSQLType
	if err := ob.db.Find(&broadcastedTransactions).Error; err != nil {
		ob.logger.Chain.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, entry := range broadcastedTransactions {
		ob.broadcastedTx[entry.Key] = entry.Hash
	}
	return nil
}

func (ob *Client) LoadLastBlock() error {
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}

	//Load persisted block number
	var lastBlockNum clienttypes.LastBlockSQLType
	if err := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumID).Error; err != nil {
		ob.logger.Chain.Info().Msg("LastBlockNum not found in DB, scan from latest")
		ob.SetLastBlockHeightScanned(bn)
	} else {
		// #nosec G701 always in range
		lastBN := int64(lastBlockNum.Num)
		ob.SetLastBlockHeightScanned(lastBN)

		//If persisted block number is too low, use the latest height
		if (bn - lastBN) > maxHeightDiff {
			ob.logger.Chain.Info().Msgf("LastBlockNum too low: %d, scan from latest", lastBlockNum.Num)
			ob.SetLastBlockHeightScanned(bn)
		}
	}

	if ob.chain.ChainId == 18444 { // bitcoin regtest: start from block 100
		ob.SetLastBlockHeightScanned(100)
	}
	ob.logger.Chain.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.GetLastBlockHeightScanned())

	return nil
}

func (ob *Client) GetBlockByNumberCached(blockNumber int64) (*BTCBlockNHeader, error) {
	if result, ok := ob.BlockCache.Get(blockNumber); ok {
		return result.(*BTCBlockNHeader), nil
	}
	// Get the block hash
	hash, err := ob.rpcClient.GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	// Get the block header
	header, err := ob.rpcClient.GetBlockHeader(hash)
	if err != nil {
		return nil, err
	}
	// Get the block with verbose transactions
	block, err := ob.rpcClient.GetBlockVerboseTx(hash)
	if err != nil {
		return nil, err
	}
	blockNheader := &BTCBlockNHeader{
		Header: header,
		Block:  block,
	}
	ob.BlockCache.Add(blockNumber, blockNheader)
	ob.BlockCache.Add(hash, blockNheader)
	return blockNheader, nil
}

// isTssTransaction checks if a given transaction was sent by TSS itself.
// An unconfirmed transaction is safe to spend only if it was sent by TSS and verified by ourselves.
func (ob *Client) isTssTransaction(txid string) bool {
	_, found := ob.includedTxHashes[txid]
	return found
}

// refreshPendingNonce tries increasing the artificial pending nonce of outTx (if lagged behind).
// There could be many (unpredictable) reasons for a pending nonce lagging behind, for example:
// 1. The zetaclient gets restarted.
// 2. The tracker is missing in zetacore.
func (ob *Client) refreshPendingNonce() {
	// get pending nonces from zetacore
	p, err := ob.coreClient.GetPendingNoncesByChain(ob.chain.ChainId)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting pending nonces")
	}

	// increase pending nonce if lagged behind
	ob.Mu.Lock()
	pendingNonce := ob.pendingNonce
	ob.Mu.Unlock()

	// #nosec G701 always non-negative
	nonceLow := uint64(p.NonceLow)
	if nonceLow > pendingNonce {
		// get the last included outTx hash
		txid, err := ob.getOutTxidByNonce(nonceLow-1, false)
		if err != nil {
			ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting last outTx txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		ob.pendingNonce = nonceLow
		ob.logger.Chain.Info().Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", ob.pendingNonce, txid)
	}
}

func (ob *Client) getOutTxidByNonce(nonce uint64, test bool) (string, error) {

	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetacore based on majority vote.
	if res := ob.getIncludedTx(nonce); res != nil {
		return res.TxID, nil
	}
	if !test { // if not unit test, get cctx from zetacore
		send, err := ob.coreClient.GetCctxByNonce(ob.chain.ChainId, nonce)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting cctx for nonce %d", nonce)
		}
		txid := send.GetCurrentOutTxParam().OutboundTxHash
		if txid == "" {
			return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
		}
		// make sure it's a real Bitcoin txid
		_, getTxResult, err := GetTxResultByHash(ob.rpcClient, txid)
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

func (ob *Client) findNonceMarkUTXO(nonce uint64, txid string) (int, error) {
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	amount := chains.NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := GetSatoshis(utxo.Amount)
		if err != nil {
			ob.logger.OutTx.Error().Err(err).Msgf("findNonceMarkUTXO: error getting satoshis for utxo %v", utxo)
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == txid && utxo.Vout == 0 {
			ob.logger.OutTx.Info().Msgf("findNonceMarkUTXO: found nonce-mark utxo with txid %s, amount %d satoshi", utxo.TxID, sats)
			return i, nil
		}
	}
	return -1, fmt.Errorf("findNonceMarkUTXO: cannot find nonce-mark utxo with nonce %d", nonce)
}

// postBlockHeader posts block header to zetacore
func (ob *Client) postBlockHeader(tip int64) error {
	ob.logger.InTx.Info().Msgf("postBlockHeader: tip %d", tip)
	bn := tip
	res, err := ob.coreClient.GetBlockHeaderChainState(ob.chain.ChainId)
	if err == nil && res.ChainState != nil && res.ChainState.EarliestHeight > 0 {
		bn = res.ChainState.LatestHeight + 1
	}
	if bn > tip {
		return fmt.Errorf("postBlockHeader: must post block confirmed block header: %d > %d", bn, tip)
	}
	res2, err := ob.GetBlockByNumberCached(bn)
	if err != nil {
		return fmt.Errorf("error getting bitcoin block %d: %s", bn, err)
	}

	var headerBuf bytes.Buffer
	err = res2.Header.Serialize(&headerBuf)
	if err != nil { // should never happen
		ob.logger.InTx.Error().Err(err).Msgf("error serializing bitcoin block header: %d", bn)
		return err
	}
	blockHash := res2.Header.BlockHash()
	_, err = ob.coreClient.PostVoteBlockHeader(
		ob.chain.ChainId,
		blockHash[:],
		res2.Block.Height,
		proofs.NewBitcoinHeader(headerBuf.Bytes()),
	)
	ob.logger.InTx.Info().Msgf("posted block header %d: %s", bn, blockHash)
	if err != nil { // error shouldn't block the process
		ob.logger.InTx.Error().Err(err).Msgf("error posting bitcoin block header: %d", bn)
	}
	return err
}

func (ob *Client) loadDB(dbpath string) error {
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
