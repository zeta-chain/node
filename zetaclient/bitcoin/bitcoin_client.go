package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ interfaces.ChainClient = &BitcoinChainClient{}

type BTCLog struct {
	ChainLogger   zerolog.Logger
	WatchInTx     zerolog.Logger
	ObserveOutTx  zerolog.Logger
	WatchUTXOS    zerolog.Logger
	WatchGasPrice zerolog.Logger
}

// BitcoinChainClient represents a chain configuration for Bitcoin
// Filled with above constants depending on chain
type BitcoinChainClient struct {
	*metricsPkg.ChainMetrics

	chain            common.Chain
	rpcClient        interfaces.BTCRPCClient
	zetaClient       interfaces.ZetaCoreBridger
	Tss              interfaces.TSSSigner
	lastBlock        int64
	lastBlockScanned int64
	BlockTime        uint64 // block time in seconds

	Mu                *sync.Mutex // lock for all the maps, utxos and core params
	pendingNonce      uint64
	includedTxHashes  map[string]bool                          // key: tx hash
	includedTxResults map[string]*btcjson.GetTransactionResult // key: chain-tss-nonce
	broadcastedTx     map[string]string                        // key: chain-tss-nonce, value: outTx hash
	utxos             []btcjson.ListUnspentResult
	params            observertypes.ChainParams

	db     *gorm.DB
	stop   chan struct{}
	logger BTCLog
	ts     *metricsPkg.TelemetryServer

	BlockCache *lru.Cache
}

const (
	minConfirmations = 0
	maxHeightDiff    = 10000
	btcBlocksPerDay  = 144
	DonationMessage  = "I am rich!"
)

func (ob *BitcoinChainClient) WithZetaClient(bridge *zetabridge.ZetaCoreBridge) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.zetaClient = bridge
}
func (ob *BitcoinChainClient) WithLogger(logger zerolog.Logger) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.logger = BTCLog{
		ChainLogger:   logger,
		WatchInTx:     logger.With().Str("module", "WatchInTx").Logger(),
		ObserveOutTx:  logger.With().Str("module", "observeOutTx").Logger(),
		WatchUTXOS:    logger.With().Str("module", "WatchUTXOS").Logger(),
		WatchGasPrice: logger.With().Str("module", "WatchGasPrice").Logger(),
	}
}

func (ob *BitcoinChainClient) WithBtcClient(client *rpcclient.Client) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.rpcClient = client
}

func (ob *BitcoinChainClient) WithChain(chain common.Chain) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.chain = chain
}

func (ob *BitcoinChainClient) SetChainParams(params observertypes.ChainParams) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.params = params
}

func (ob *BitcoinChainClient) GetChainParams() observertypes.ChainParams {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.params
}

// NewBitcoinClient returns a new configuration based on supplied target chain
func NewBitcoinClient(
	chain common.Chain,
	bridge interfaces.ZetaCoreBridger,
	tss interfaces.TSSSigner,
	dbpath string,
	metrics *metricsPkg.Metrics,
	logger zerolog.Logger,
	btcCfg config.BTCConfig,
	ts *metricsPkg.TelemetryServer,
) (*BitcoinChainClient, error) {
	ob := BitcoinChainClient{
		ChainMetrics: metricsPkg.NewChainMetrics(chain.ChainName.String(), metrics),
		ts:           ts,
	}
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.Mu = &sync.Mutex{}
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
	ob.includedTxHashes = make(map[string]bool)
	ob.includedTxResults = make(map[string]*btcjson.GetTransactionResult)
	ob.broadcastedTx = make(map[string]string)
	ob.params = btcCfg.ChainParams

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

	ob.BlockCache, err = lru.New(btcBlocksPerDay)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("failed to create bitcoin block cache")
		return nil, err
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
	go ob.ExternalChainWatcherForNewInboundTrackerSuggestions()
	go ob.RPCStatus()
}

func (ob *BitcoinChainClient) RPCStatus() {
	ob.logger.ChainLogger.Info().Msgf("RPCStatus is starting")
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			//ob.logger.ChainLogger.Info().Msgf("RPCStatus is running")
			bn, err := ob.rpcClient.GetBlockCount()
			if err != nil {
				ob.logger.ChainLogger.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			hash, err := ob.rpcClient.GetBlockHash(bn)
			if err != nil {
				ob.logger.ChainLogger.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			header, err := ob.rpcClient.GetBlockHeader(hash)
			if err != nil {
				ob.logger.ChainLogger.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			blockTime := header.Timestamp
			elapsedSeconds := time.Since(blockTime).Seconds()
			if elapsedSeconds > 1200 {
				ob.logger.ChainLogger.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			tssAddr := ob.Tss.BTCAddressWitnessPubkeyHash()
			res, err := ob.rpcClient.ListUnspentMinMaxAddresses(0, 1000000, []btcutil.Address{tssAddr})
			if err != nil {
				ob.logger.ChainLogger.Error().Err(err).Msg("RPC status check: can't list utxos of TSS address; wallet or loaded? TSS address is not imported? ")
				continue
			}
			if len(res) == 0 {
				ob.logger.ChainLogger.Error().Err(err).Msg("RPC status check: TSS address has no utxos; TSS address is not imported? ")
				continue
			}
			ob.logger.ChainLogger.Info().Msgf("[OK] RPC status check: latest block number %d, timestamp %s (%.fs ago), tss addr %s, #utxos: %d", bn, blockTime, elapsedSeconds, tssAddr, len(res))

		case <-ob.stop:
			return
		}
	}
}

func (ob *BitcoinChainClient) Stop() {
	ob.logger.ChainLogger.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop
	ob.logger.ChainLogger.Info().Msgf("%s observer stopped", ob.chain.String())
}

func (ob *BitcoinChainClient) SetLastBlockHeight(height int64) {
	if height < 0 {
		panic("lastBlock is negative")
	}
	atomic.StoreInt64(&ob.lastBlock, height)
}

func (ob *BitcoinChainClient) GetLastBlockHeight() int64 {
	height := atomic.LoadInt64(&ob.lastBlock)
	if height < 0 {
		panic("lastBlock is negative")
	}
	return height
}

func (ob *BitcoinChainClient) SetLastBlockHeightScanned(height int64) {
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	atomic.StoreInt64(&ob.lastBlockScanned, height)
	// #nosec G701 checked as positive
	ob.ts.SetLastScannedBlockNumber((ob.chain.ChainId), uint64(height))
}

func (ob *BitcoinChainClient) GetLastBlockHeightScanned() int64 {
	height := atomic.LoadInt64(&ob.lastBlockScanned)
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	return height
}

func (ob *BitcoinChainClient) GetPendingNonce() uint64 {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.pendingNonce
}

// GetBaseGasPrice ...
// TODO: implement
// https://github.com/zeta-chain/node/issues/868
func (ob *BitcoinChainClient) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

func (ob *BitcoinChainClient) WatchInTx() {
	ticker, err := NewDynamicTicker("Bitcoin_WatchInTx", ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.WatchInTx.Error().Err(err).Msg("WatchInTx error")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.observeInTx()
			if err != nil {
				ob.logger.WatchInTx.Error().Err(err).Msg("WatchInTx error observing in tx")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.WatchInTx)
		case <-ob.stop:
			ob.logger.WatchInTx.Info().Msg("WatchInTx stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) postBlockHeader(tip int64) error {
	ob.logger.WatchInTx.Info().Msgf("postBlockHeader: tip %d", tip)
	bn := tip
	res, err := ob.zetaClient.GetBlockHeaderStateByChain(ob.chain.ChainId)
	if err == nil && res.BlockHeaderState != nil && res.BlockHeaderState.EarliestHeight > 0 {
		bn = res.BlockHeaderState.LatestHeight + 1
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
		ob.logger.WatchInTx.Error().Err(err).Msgf("error serializing bitcoin block header: %d", bn)
		return err
	}
	blockHash := res2.Header.BlockHash()
	_, err = ob.zetaClient.PostAddBlockHeader(
		ob.chain.ChainId,
		blockHash[:],
		res2.Block.Height,
		common.NewBitcoinHeader(headerBuf.Bytes()),
	)
	ob.logger.WatchInTx.Info().Msgf("posted block header %d: %s", bn, blockHash)
	if err != nil { // error shouldn't block the process
		ob.logger.WatchInTx.Error().Err(err).Msgf("error posting bitcoin block header: %d", bn)
	}
	return err
}

func (ob *BitcoinChainClient) observeInTx() error {
	// make sure inbound TXS / Send is enabled by the protocol
	flags, err := ob.zetaClient.GetCrosschainFlags()
	if err != nil {
		return err
	}
	if !flags.IsInboundEnabled {
		return errors.New("inbound TXS / Send has been disabled by the protocol")
	}

	// get and update latest block height
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("observeInTxBTC: error getting block count: %s", err)
	}
	if cnt < 0 {
		return fmt.Errorf("observeInTxBTC: block count is negative: %d", cnt)
	}
	ob.SetLastBlockHeight(cnt)

	// skip if current height is too low
	// #nosec G701 always in range
	confirmedBlockNum := cnt - int64(ob.GetChainParams().ConfirmationCount)
	if confirmedBlockNum < 0 {
		return fmt.Errorf("observeInTxBTC: skipping observer, current block number %d is too low", cnt)
	}

	// skip if no new block is confirmed
	lastScanned := ob.GetLastBlockHeightScanned()
	if lastScanned >= confirmedBlockNum {
		return nil
	}

	// query incoming gas asset to TSS address
	{
		bn := lastScanned + 1
		res, err := ob.GetBlockByNumberCached(bn)
		if err != nil {
			ob.logger.WatchInTx.Error().Err(err).Msgf("observeInTxBTC: error getting bitcoin block %d", bn)
			return err
		}
		ob.logger.WatchInTx.Info().Msgf("observeInTxBTC: block %d has %d txs, current block %d, last block %d",
			bn, len(res.Block.Tx), cnt, lastScanned)

		// print some debug information
		if len(res.Block.Tx) > 1 {
			for idx, tx := range res.Block.Tx {
				ob.logger.WatchInTx.Debug().Msgf("BTC InTX |  %d: %s\n", idx, tx.Txid)
				for vidx, vout := range tx.Vout {
					ob.logger.WatchInTx.Debug().Msgf("vout %d \n value: %v\n scriptPubKey: %v\n", vidx, vout.Value, vout.ScriptPubKey.Hex)
				}
			}
		}

		// add block header to zetabridge
		if flags.BlockHeaderVerificationFlags != nil && flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
			err = ob.postBlockHeader(bn)
			if err != nil {
				ob.logger.WatchInTx.Warn().Err(err).Msgf("observeInTxBTC: error posting block header %d", bn)
			}
		}

		tssAddress := ob.Tss.BTCAddress()
		// #nosec G701 always positive
		inTxs := FilterAndParseIncomingTx(
			res.Block.Tx,
			uint64(res.Block.Height),
			tssAddress,
			&ob.logger.WatchInTx,
			ob.chain.ChainId,
		)

		// post inbound vote message to zetabridge
		for _, inTx := range inTxs {
			msg := ob.GetInboundVoteMessageFromBtcEvent(inTx)
			zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundExecutionGasLimit, msg)
			if err != nil {
				ob.logger.WatchInTx.Error().Err(err).Msgf("observeInTxBTC: error posting to zeta core for tx %s", inTx.TxHash)
				return err // we have to re-scan this block next time
			} else if zetaHash != "" {
				ob.logger.WatchInTx.Info().Msgf("observeInTxBTC: BTC deposit detected and reported: PostVoteInbound zeta tx: %s ballot %s", zetaHash, ballot)
			}
		}

		// Save LastBlockHeight
		ob.SetLastBlockHeightScanned(bn)
		// #nosec G701 always positive
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(uint64(bn))).Error; err != nil {
			ob.logger.WatchInTx.Error().Err(err).Msgf("observeInTxBTC: error writing last scanned block %d to db", bn)
		}
	}

	return nil
}

// ConfirmationsThreshold returns number of required Bitcoin confirmations depending on sent BTC amount.
func (ob *BitcoinChainClient) ConfirmationsThreshold(amount *big.Int) int64 {
	if amount.Cmp(big.NewInt(200000000)) >= 0 {
		return 6
	}
	return 2
}

// IsSendOutTxProcessed returns isIncluded(or inMempool), isConfirmed, Error
func (ob *BitcoinChainClient) IsSendOutTxProcessed(sendHash string, nonce uint64, _ common.CoinType, logger zerolog.Logger) (bool, bool, error) {
	outTxID := ob.GetTxID(nonce)
	logger.Info().Msgf("IsSendOutTxProcessed %s", outTxID)

	ob.Mu.Lock()
	txnHash, broadcasted := ob.broadcastedTx[outTxID]
	res, included := ob.includedTxResults[outTxID]
	ob.Mu.Unlock()

	// Get original cctx parameters
	params, err := ob.GetCctxParams(nonce)
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
		txResult, inMempool := ob.checkIncludedTx(txnHash, params)
		if txResult == nil {
			ob.logger.ObserveOutTx.Error().Err(err).Msg("IsSendOutTxProcessed: checkIncludedTx failed")
			return false, false, err
		} else if inMempool { // still in mempool (should avoid unnecessary Tss keysign)
			ob.logger.ObserveOutTx.Info().Msgf("IsSendOutTxProcessed: outTx %s is still in mempool", outTxID)
			return true, false, nil
		}
		// included
		ob.setIncludedTx(nonce, txResult)

		// Get tx result again in case it is just included
		res = ob.getIncludedTx(nonce)
		if res == nil {
			return false, false, nil
		}
		ob.logger.ObserveOutTx.Info().Msgf("IsSendOutTxProcessed: setIncludedTx succeeded for outTx %s", outTxID)
	}

	// It's safe to use cctx's amount to post confirmation because it has already been verified in observeOutTx()
	amountInSat := params.Amount.BigInt()
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		return true, false, nil
	}

	logger.Debug().Msgf("Bitcoin outTx confirmed: txid %s, amount %s\n", res.TxID, amountInSat.String())
	zetaHash, ballot, err := ob.zetaClient.PostVoteOutbound(
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
		logger.Error().Err(err).Msgf("IsSendOutTxProcessed: error confirming bitcoin outTx %s, nonce %d ballot %s", res.TxID, nonce, ballot)
	} else if zetaHash != "" {
		logger.Info().Msgf("IsSendOutTxProcessed: confirmed Bitcoin outTx %s, zeta tx hash %s nonce %d ballot %s", res.TxID, zetaHash, nonce, ballot)
	}
	return true, true, nil
}

func (ob *BitcoinChainClient) WatchGasPrice() {
	ticker, err := NewDynamicTicker("Bitcoin_WatchGasPrice", ob.GetChainParams().GasPriceTicker)
	if err != nil {
		ob.logger.WatchGasPrice.Error().Err(err).Msg("WatchGasPrice error")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.WatchGasPrice.Error().Err(err).Msg("PostGasPrice error on " + ob.chain.String())
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.logger.WatchGasPrice)
		case <-ob.stop:
			ob.logger.WatchGasPrice.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) PostGasPrice() error {
	if ob.chain.ChainId == 18444 { //bitcoin regtest; hardcode here since this RPC is not available on regtest
		bn, err := ob.rpcClient.GetBlockCount()
		if err != nil {
			return err
		}
		// #nosec G701 always in range
		zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, 1, "100", uint64(bn))
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
	if *feeResult.FeeRate > math.MaxInt64 {
		return fmt.Errorf("gas price is too large: %f", *feeResult.FeeRate)
	}
	feeRatePerByte := FeeRateToSatPerByte(*feeResult.FeeRate)
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}
	// #nosec G701 always positive
	zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, feeRatePerByte.Uint64(), "100", uint64(bn))
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

// FilterAndParseIncomingTx given txs list returned by the "getblock 2" RPC command, return the txs that are relevant to us
// relevant tx must have the following vouts as the first two vouts:
// vout0: p2wpkh to the TSS address (targetAddress)
// vout1: OP_RETURN memo, base64 encoded
func FilterAndParseIncomingTx(
	txs []btcjson.TxRawResult,
	blockNumber uint64,
	targetAddress string,
	logger *zerolog.Logger,
	chainID int64,
) []*BTCInTxEvnet {
	inTxs := make([]*BTCInTxEvnet, 0)
	for idx, tx := range txs {
		if idx == 0 {
			continue // the first tx is coinbase; we do not process coinbase tx
		}
		inTx, err := GetBtcEvent(tx, targetAddress, blockNumber, logger, chainID)
		if err != nil {
			logger.Error().Err(err).Msg("error getting btc event")
			continue
		}
		if inTx != nil {
			inTxs = append(inTxs, inTx)
		}
	}
	return inTxs
}

func (ob *BitcoinChainClient) GetInboundVoteMessageFromBtcEvent(inTx *BTCInTxEvnet) *types.MsgVoteOnObservedInboundTx {
	ob.logger.WatchInTx.Debug().Msgf("Processing inTx: %s", inTx.TxHash)
	amount := big.NewFloat(inTx.Value)
	amount = amount.Mul(amount, big.NewFloat(1e8))
	amountInt, _ := amount.Int(nil)
	message := hex.EncodeToString(inTx.MemoBytes)
	return zetabridge.GetInBoundVoteMessage(
		inTx.FromAddress,
		ob.chain.ChainId,
		inTx.FromAddress,
		inTx.FromAddress,
		ob.zetaClient.ZetaChain().ChainId,
		cosmosmath.NewUintFromBigInt(amountInt),
		message,
		inTx.TxHash,
		inTx.BlockNumber,
		0,
		common.CoinType_Gas,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		0,
	)
}

func GetBtcEvent(
	tx btcjson.TxRawResult,
	targetAddress string,
	blockNumber uint64,
	logger *zerolog.Logger,
	chainID int64,
) (*BTCInTxEvnet, error) {
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
				return nil, err
			}

			bitcoinNetParams, err := common.BitcoinNetParamsFromChainID(chainID)
			if err != nil {
				return nil, fmt.Errorf("btc: error getting bitcoin net params : %v", err)
			}
			wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(hash, bitcoinNetParams)
			if err != nil {
				return nil, err
			}
			if wpkhAddress.EncodeAddress() != targetAddress {
				return nil, err
			}
			// deposit amount has to be no less than the minimum depositor fee
			if out.Value < BtcDepositorFeeMin {
				return nil, fmt.Errorf("btc deposit amount %v in txid %s is less than minimum depositor fee %v", value, tx.Txid, BtcDepositorFeeMin)
			}
			value = out.Value - BtcDepositorFeeMin

			out = tx.Vout[1]
			script = out.ScriptPubKey.Hex
			if len(script) >= 4 && script[:2] == "6a" { // OP_RETURN
				memoSize, err := strconv.ParseInt(script[2:4], 16, 32)
				if err != nil {
					return nil, errors.Wrapf(err, "error decoding pubkey hash")
				}
				if int(memoSize) != (len(script)-4)/2 {
					return nil, fmt.Errorf("memo size mismatch: %d != %d", memoSize, (len(script)-4)/2)
				}
				memoBytes, err := hex.DecodeString(script[4:])
				if err != nil {
					logger.Warn().Err(err).Msgf("error hex decoding memo")
					return nil, fmt.Errorf("error hex decoding memo: %s", err)
				}
				if bytes.Equal(memoBytes, []byte(DonationMessage)) {
					logger.Info().Msgf("donation tx: %s; value %f", tx.Txid, value)
					return nil, fmt.Errorf("donation tx: %s; value %f", tx.Txid, value)
				}
				memo = memoBytes
				found = true
			}
		}
	}
	if found {
		logger.Info().Msgf("found bitcoin intx: %s", tx.Txid)
		var fromAddress string
		if len(tx.Vin) > 0 {
			vin := tx.Vin[0]
			//log.Info().Msgf("vin: %v", vin.Witness)
			if len(vin.Witness) == 2 {
				pk := vin.Witness[1]
				pkBytes, err := hex.DecodeString(pk)
				if err != nil {
					return nil, errors.Wrapf(err, "error decoding pubkey")
				}
				hash := btcutil.Hash160(pkBytes)

				bitcoinNetParams, err := common.BitcoinNetParamsFromChainID(chainID)
				if err != nil {
					return nil, fmt.Errorf("btc: error getting bitcoin net params : %v", err)
				}

				addr, err := btcutil.NewAddressWitnessPubKeyHash(hash, bitcoinNetParams)
				if err != nil {
					return nil, errors.Wrapf(err, "error decoding pubkey hash")
				}
				fromAddress = addr.EncodeAddress()
			}
		}
		return &BTCInTxEvnet{
			FromAddress: fromAddress,
			ToAddress:   targetAddress,
			Value:       value,
			MemoBytes:   memo,
			BlockNumber: blockNumber,
			TxHash:      tx.Txid,
		}, nil
	}
	return nil, nil
}

func (ob *BitcoinChainClient) WatchUTXOS() {
	ticker, err := NewDynamicTicker("Bitcoin_WatchUTXOS", ob.GetChainParams().WatchUtxoTicker)
	if err != nil {
		ob.logger.WatchUTXOS.Error().Err(err).Msg("WatchUTXOS error")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.FetchUTXOS()
			if err != nil {
				ob.logger.WatchUTXOS.Error().Err(err).Msg("error fetching btc utxos")
			}
			ticker.UpdateInterval(ob.GetChainParams().WatchUtxoTicker, ob.logger.WatchUTXOS)
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

	// List all unspent UTXOs (160ms)
	tssAddr := ob.Tss.BTCAddress()
	address, err := common.DecodeBtcAddress(tssAddr, ob.chain.ChainId)
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
		if utxo.Amount < BtcDepositorFeeMin {
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
	ob.ts.SetNumberOfUTXOs(len(utxosFiltered))
	ob.utxos = utxosFiltered
	ob.Mu.Unlock()
	return nil
}

// isTssTransaction checks if a given transaction was sent by TSS itself.
// An unconfirmed transaction is safe to spend only if it was sent by TSS and verified by ourselves.
func (ob *BitcoinChainClient) isTssTransaction(txid string) bool {
	_, found := ob.includedTxHashes[txid]
	return found
}

// refreshPendingNonce tries increasing the artificial pending nonce of outTx (if lagged behind).
// There could be many (unpredictable) reasons for a pending nonce lagging behind, for example:
// 1. The zetaclient gets restarted.
// 2. The tracker is missing in zetabridge.
func (ob *BitcoinChainClient) refreshPendingNonce() {
	// get pending nonces from zetabridge
	p, err := ob.zetaClient.GetPendingNoncesByChain(ob.chain.ChainId)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("refreshPendingNonce: error getting pending nonces")
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
			ob.logger.ChainLogger.Error().Err(err).Msg("refreshPendingNonce: error getting last outTx txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		ob.pendingNonce = nonceLow
		ob.logger.ChainLogger.Info().Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", ob.pendingNonce, txid)
	}
}

func (ob *BitcoinChainClient) getOutTxidByNonce(nonce uint64, test bool) (string, error) {

	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetabridge based on majority vote.
	if res := ob.getIncludedTx(nonce); res != nil {
		return res.TxID, nil
	}
	if !test { // if not unit test, get cctx from zetabridge
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
	amount := common.NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := GetSatoshis(utxo.Amount)
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
func (ob *BitcoinChainClient) SelectUTXOs(amount float64, utxosToSpend uint16, nonce uint64, consolidateRank uint16, test bool) ([]btcjson.ListUnspentResult, float64, uint16, float64, error) {
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
func (ob *BitcoinChainClient) SaveBroadcastedTx(txHash string, nonce uint64) {
	outTxID := ob.GetTxID(nonce)
	ob.Mu.Lock()
	ob.broadcastedTx[outTxID] = txHash
	ob.Mu.Unlock()

	broadcastEntry := clienttypes.ToOutTxHashSQLType(txHash, outTxID)
	if err := ob.db.Save(&broadcastEntry).Error; err != nil {
		ob.logger.ObserveOutTx.Error().Err(err).Msgf("SaveBroadcastedTx: error saving broadcasted txHash %s for outTx %s", txHash, outTxID)
	}
	ob.logger.ObserveOutTx.Info().Msgf("SaveBroadcastedTx: saved broadcasted txHash %s for outTx %s", txHash, outTxID)
}

func (ob *BitcoinChainClient) GetCctxParams(nonce uint64) (types.OutboundTxParams, error) {
	send, err := ob.zetaClient.GetCctxByNonce(ob.chain.ChainId, nonce)
	if err != nil {
		return types.OutboundTxParams{}, err
	}
	if send.GetCurrentOutTxParam() == nil { // never happen
		return types.OutboundTxParams{}, fmt.Errorf("GetPendingCctx: nil outbound tx params")
	}
	return *send.GetCurrentOutTxParam(), nil
}

func (ob *BitcoinChainClient) observeOutTx() {
	ticker, err := NewDynamicTicker("Bitcoin_observeOutTx", ob.GetChainParams().OutTxTicker)
	if err != nil {
		ob.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: error creating ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain.ChainId, interfaces.Ascending)
			if err != nil {
				ob.logger.ObserveOutTx.Error().Err(err).Msg("observeOutTx: error GetAllOutTxTrackerByChain")
				continue
			}
			for _, tracker := range trackers {
				// get original cctx parameters
				outTxID := ob.GetTxID(tracker.Nonce)
				params, err := ob.GetCctxParams(tracker.Nonce)
				if err != nil {
					ob.logger.ObserveOutTx.Info().Err(err).Msgf("observeOutTx: can't find cctx for nonce %d", tracker.Nonce)
					break
				}
				if tracker.Nonce != params.OutboundTxTssNonce { // Tanmay: it doesn't hurt to check
					ob.logger.ObserveOutTx.Error().Msgf("observeOutTx: tracker nonce %d not match cctx nonce %d", tracker.Nonce, params.OutboundTxTssNonce)
					break
				}
				if len(tracker.HashList) > 1 {
					ob.logger.ObserveOutTx.Warn().Msgf("observeOutTx: oops, outTxID %s got multiple (%d) outTx hashes", outTxID, len(tracker.HashList))
				}
				// iterate over all txHashes to find the truly included one.
				// we do it this (inefficient) way because we don't rely on the first one as it may be a false positive (for unknown reason).
				txCount := 0
				var txResult *btcjson.GetTransactionResult
				for _, txHash := range tracker.HashList {
					result, inMempool := ob.checkIncludedTx(txHash.TxHash, params)
					if result != nil && !inMempool { // included
						txCount++
						txResult = result
						ob.logger.ObserveOutTx.Info().Msgf("observeOutTx: included outTx %s for chain %d nonce %d", txHash.TxHash, ob.chain.ChainId, tracker.Nonce)
						if txCount > 1 {
							ob.logger.ObserveOutTx.Error().Msgf(
								"observeOutTx: checkIncludedTx passed, txCount %d chain %d nonce %d result %v", txCount, ob.chain.ChainId, tracker.Nonce, result)
						}
					}
				}
				if txCount == 1 { // should be only one txHash included for each nonce
					ob.setIncludedTx(tracker.Nonce, txResult)
				} else if txCount > 1 {
					ob.removeIncludedTx(tracker.Nonce) // we can't tell which txHash is true, so we remove all (if any) to be safe
					ob.logger.ObserveOutTx.Error().Msgf("observeOutTx: included multiple (%d) outTx for chain %d nonce %d", txCount, ob.chain.ChainId, tracker.Nonce)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutTxTicker, ob.logger.ObserveOutTx)
		case <-ob.stop:
			ob.logger.ObserveOutTx.Info().Msg("observeOutTx stopped")
			return
		}
	}
}

// checkIncludedTx checks if a txHash is included and returns (txResult, inMempool)
// Note: if txResult is nil, then inMempool flag should be ignored.
func (ob *BitcoinChainClient) checkIncludedTx(txHash string, params types.OutboundTxParams) (*btcjson.GetTransactionResult, bool) {
	outTxID := ob.GetTxID(params.OutboundTxTssNonce)
	hash, getTxResult, err := ob.GetTxResultByHash(txHash)
	if err != nil {
		ob.logger.ObserveOutTx.Error().Err(err).Msgf("checkIncludedTx: error GetTxResultByHash: %s", txHash)
		return nil, false
	}
	if txHash != getTxResult.TxID { // just in case, we'll use getTxResult.TxID later
		ob.logger.ObserveOutTx.Error().Msgf("checkIncludedTx: inconsistent txHash %s and getTxResult.TxID %s", txHash, getTxResult.TxID)
		return nil, false
	}
	if getTxResult.Confirmations >= 0 { // check included tx only
		err = ob.checkTssOutTxResult(hash, getTxResult, params, params.OutboundTxTssNonce)
		if err != nil {
			ob.logger.ObserveOutTx.Error().Err(err).Msgf("checkIncludedTx: error verify bitcoin outTx %s outTxID %s", txHash, outTxID)
			return nil, false
		}
		return getTxResult, false // included
	}
	return getTxResult, true // in mempool
}

// setIncludedTx saves included tx result in memory
func (ob *BitcoinChainClient) setIncludedTx(nonce uint64, getTxResult *btcjson.GetTransactionResult) {
	txHash := getTxResult.TxID
	outTxID := ob.GetTxID(nonce)

	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	res, found := ob.includedTxResults[outTxID]

	if !found { // not found.
		ob.includedTxHashes[txHash] = true
		ob.includedTxResults[outTxID] = getTxResult // include new outTx and enforce rigid 1-to-1 mapping: nonce <===> txHash
		if nonce >= ob.pendingNonce {               // try increasing pending nonce on every newly included outTx
			ob.pendingNonce = nonce + 1
		}
		ob.logger.ObserveOutTx.Info().Msgf("setIncludedTx: included new bitcoin outTx %s outTxID %s pending nonce %d", txHash, outTxID, ob.pendingNonce)
	} else if txHash == res.TxID { // found same hash.
		ob.includedTxResults[outTxID] = getTxResult // update tx result as confirmations may increase
		if getTxResult.Confirmations > res.Confirmations {
			ob.logger.ObserveOutTx.Info().Msgf("setIncludedTx: bitcoin outTx %s got confirmations %d", txHash, getTxResult.Confirmations)
		}
	} else { // found other hash.
		// be alert for duplicate payment!!! As we got a new hash paying same cctx (for whatever reason).
		delete(ob.includedTxResults, outTxID) // we can't tell which txHash is true, so we remove all to be safe
		ob.logger.ObserveOutTx.Error().Msgf("setIncludedTx: duplicate payment by bitcoin outTx %s outTxID %s, prior outTx %s", txHash, outTxID, res.TxID)
	}
}

// getIncludedTx gets the receipt and transaction from memory
func (ob *BitcoinChainClient) getIncludedTx(nonce uint64) *btcjson.GetTransactionResult {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.includedTxResults[ob.GetTxID(nonce)]
}

// removeIncludedTx removes included tx from memory
func (ob *BitcoinChainClient) removeIncludedTx(nonce uint64) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	txResult, found := ob.includedTxResults[ob.GetTxID(nonce)]
	if found {
		delete(ob.includedTxHashes, txResult.TxID)
		delete(ob.includedTxResults, ob.GetTxID(nonce))
	}
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
	}

	// res.Confirmations < 0 (meaning not included)
	return btcjson.TxRawResult{}, fmt.Errorf("getRawTxResult: tx %s not included yet", hash)
}

// checkTSSVin checks vin is valid if:
//   - The first input is the nonce-mark
//   - All inputs are from TSS address
func (ob *BitcoinChainClient) checkTSSVin(vins []btcjson.Vin, nonce uint64) error {
	// vins: [nonce-mark, UTXO1, UTXO2, ...]
	if nonce > 0 && len(vins) <= 1 {
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

// checkTSSVout vout is valid if:
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
		amount, err := GetSatoshis(vout.Value)
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
			if amount != common.NonceMarkAmount(nonce) {
				return fmt.Errorf("checkTSSVout: nonce-mark amount %d not match nonce-mark amount %d", amount, common.NonceMarkAmount(nonce))
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
		// #nosec G701 always in range
		lastBN := int64(lastBlockNum.Num)
		ob.SetLastBlockHeightScanned(lastBN)

		//If persisted block number is too low, use the latest height
		if (bn - lastBN) > maxHeightDiff {
			ob.logger.ChainLogger.Info().Msgf("LastBlockNum too low: %d, scan from latest", lastBlockNum.Num)
			ob.SetLastBlockHeightScanned(bn)
		}
	}

	if ob.chain.ChainId == 18444 { // bitcoin regtest: start from block 100
		ob.SetLastBlockHeightScanned(100)
	}
	ob.logger.ChainLogger.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.GetLastBlockHeightScanned())

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

type BTCBlockNHeader struct {
	Header *wire.BlockHeader
	Block  *btcjson.GetBlockVerboseTxResult
}

func (ob *BitcoinChainClient) GetBlockByNumberCached(blockNumber int64) (*BTCBlockNHeader, error) {
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
