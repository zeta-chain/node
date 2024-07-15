// Package observer implements the Bitcoin chain observer
package observer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sort"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/bg"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	// btcBlocksPerDay represents Bitcoin blocks per days for LRU block cache size
	btcBlocksPerDay = 144

	// RegnetStartBlock is the hardcoded start block for regnet
	RegnetStartBlock = 100

	// BigValueSats contains the threshold to determine a big value in Bitcoin represents 2 BTC
	BigValueSats = 200000000

	// BigValueConfirmationCount represents the number of confirmation necessary for bigger values: 6 confirmations
	BigValueConfirmationCount = 6
)

var _ interfaces.ChainObserver = &Observer{}

// Logger contains list of loggers used by Bitcoin chain observer
type Logger struct {
	// base.Logger contains a list of base observer loggers
	base.ObserverLogger

	// UTXOs is the logger for UTXOs management
	UTXOs zerolog.Logger
}

// BTCInboundEvent represents an incoming transaction event
// TODO(revamp): Move to inbound
type BTCInboundEvent struct {
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

// BTCBlockNHeader contains bitcoin block and the header
type BTCBlockNHeader struct {
	Header *wire.BlockHeader
	Block  *btcjson.GetBlockVerboseTxResult
}

// Observer is the Bitcoin chain observer
type Observer struct {
	// base.Observer implements the base chain observer
	base.Observer

	// netParams contains the Bitcoin network parameters
	netParams *chaincfg.Params

	// btcClient is the Bitcoin RPC client that interacts with the Bitcoin node
	btcClient interfaces.BTCRPCClient

	// pendingNonce is the outbound artificial pending nonce
	pendingNonce uint64

	// utxos contains the UTXOs owned by the TSS address
	utxos []btcjson.ListUnspentResult

	// includedTxHashes indexes included tx with tx hash
	includedTxHashes map[string]bool

	// includedTxResults indexes tx results with the outbound tx identifier
	includedTxResults map[string]*btcjson.GetTransactionResult

	// broadcastedTx indexes the outbound hash with the outbound tx identifier
	broadcastedTx map[string]string

	// logger contains the loggers used by the bitcoin observer
	logger Logger
}

// NewObserver returns a new Bitcoin chain observer
func NewObserver(
	chain chains.Chain,
	btcClient interfaces.BTCRPCClient,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*Observer, error) {
	// create base observer
	baseObserver, err := base.NewObserver(
		chain,
		chainParams,
		zetacoreClient,
		tss,
		btcBlocksPerDay,
		base.DefaultHeaderCacheSize,
		ts,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// get the bitcoin network params
	netParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
	if err != nil {
		return nil, fmt.Errorf("error getting net params for chain %d: %s", chain.ChainId, err)
	}

	// create bitcoin observer
	ob := &Observer{
		Observer:          *baseObserver,
		netParams:         netParams,
		btcClient:         btcClient,
		pendingNonce:      0,
		utxos:             []btcjson.ListUnspentResult{},
		includedTxHashes:  make(map[string]bool),
		includedTxResults: make(map[string]*btcjson.GetTransactionResult),
		broadcastedTx:     make(map[string]string),
		logger: Logger{
			ObserverLogger: *baseObserver.Logger(),
			UTXOs:          baseObserver.Logger().Chain.With().Str("module", "utxos").Logger(),
		},
	}

	// load btc chain observer DB
	err = ob.LoadDB(dbpath)
	if err != nil {
		return nil, err
	}

	return ob, nil
}

// BtcClient returns the btc client
func (ob *Observer) BtcClient() interfaces.BTCRPCClient {
	return ob.btcClient
}

// WithBtcClient attaches a new btc client to the observer
func (ob *Observer) WithBtcClient(client interfaces.BTCRPCClient) {
	ob.btcClient = client
}

// SetChainParams sets the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) SetChainParams(params observertypes.ChainParams) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.WithChainParams(params)
}

// GetChainParams returns the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) GetChainParams() observertypes.ChainParams {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.ChainParams()
}

// Start starts the Go routine processes to observe the Bitcoin chain
func (ob *Observer) Start(ctx context.Context) {
	ob.Logger().Chain.Info().Msgf("observer is starting for chain %d", ob.Chain().ChainId)

	// watch bitcoin chain for incoming txs and post votes to zetacore
	bg.Work(ctx, ob.WatchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))

	// watch bitcoin chain for outgoing txs status
	bg.Work(ctx, ob.WatchOutbound, bg.WithName("WatchOutbound"), bg.WithLogger(ob.Logger().Outbound))

	// watch bitcoin chain for UTXOs owned by the TSS address
	bg.Work(ctx, ob.WatchUTXOs, bg.WithName("WatchUTXOs"), bg.WithLogger(ob.Logger().Outbound))

	// watch bitcoin chain for gas rate and post to zetacore
	bg.Work(ctx, ob.WatchGasPrice, bg.WithName("WatchGasPrice"), bg.WithLogger(ob.Logger().GasPrice))

	// watch zetacore for bitcoin inbound trackers
	bg.Work(ctx, ob.WatchInboundTracker, bg.WithName("WatchInboundTracker"), bg.WithLogger(ob.Logger().Inbound))

	// watch the RPC status of the bitcoin chain
	bg.Work(ctx, ob.WatchRPCStatus, bg.WithName("WatchRPCStatus"), bg.WithLogger(ob.Logger().Chain))
}

// WatchRPCStatus watches the RPC status of the Bitcoin chain
// TODO(revamp): move ticker related functions to a specific file
// TODO(revamp): move inner logic in a separate function
func (ob *Observer) WatchRPCStatus(_ context.Context) error {
	ob.logger.Chain.Info().Msgf("RPCStatus is starting")
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}

			bn, err := ob.btcClient.GetBlockCount()
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}

			hash, err := ob.btcClient.GetBlockHash(bn)
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}

			header, err := ob.btcClient.GetBlockHeader(hash)
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

			tssAddr := ob.TSS().BTCAddressWitnessPubkeyHash()
			res, err := ob.btcClient.ListUnspentMinMaxAddresses(0, 1000000, []btcutil.Address{tssAddr})
			if err != nil {
				ob.logger.Chain.Error().
					Err(err).
					Msg("RPC status check: can't list utxos of TSS address; wallet or loaded? TSS address is not imported? ")
				continue
			}

			if len(res) == 0 {
				ob.logger.Chain.Error().
					Err(err).
					Msg("RPC status check: TSS address has no utxos; TSS address is not imported? ")
				continue
			}

			ob.logger.Chain.Info().
				Msgf("[OK] RPC status check: latest block number %d, timestamp %s (%.fs ago), tss addr %s, #utxos: %d", bn, blockTime, elapsedSeconds, tssAddr, len(res))

		case <-ob.StopChannel():
			return nil
		}
	}
}

// GetPendingNonce returns the artificial pending nonce
// Note: pending nonce is accessed concurrently
func (ob *Observer) GetPendingNonce() uint64 {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.pendingNonce
}

// ConfirmationsThreshold returns number of required Bitcoin confirmations depending on sent BTC amount.
func (ob *Observer) ConfirmationsThreshold(amount *big.Int) int64 {
	if amount.Cmp(big.NewInt(BigValueSats)) >= 0 {
		return BigValueConfirmationCount
	}
	if BigValueConfirmationCount < ob.GetChainParams().ConfirmationCount {
		return BigValueConfirmationCount
	}

	// #nosec G115 always in range
	return int64(ob.GetChainParams().ConfirmationCount)
}

// WatchGasPrice watches Bitcoin chain for gas rate and post to zetacore
// TODO(revamp): move ticker related functions to a specific file
// TODO(revamp): move inner logic in a separate function
func (ob *Observer) WatchGasPrice(ctx context.Context) error {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice(ctx)
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchGasPrice", ob.GetChainParams().GasPriceTicker)
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msg("error creating ticker")
		return err
	}
	ob.logger.GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.Chain().ChainId, ob.GetChainParams().GasPriceTicker)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.PostGasPrice(ctx)
			if err != nil {
				ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.logger.GasPrice)
		case <-ob.StopChannel():
			ob.logger.GasPrice.Info().Msgf("WatchGasPrice stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// PostGasPrice posts gas price to zetacore
// TODO(revamp): move to gas price file
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	var (
		err              error
		feeRateEstimated uint64
	)

	// special handle regnet and testnet gas rate
	// regnet:  RPC 'EstimateSmartFee' is not available
	// testnet: RPC 'EstimateSmartFee' returns unreasonable high gas rate
	if ob.Chain().NetworkType != chains.NetworkType_mainnet {
		feeRateEstimated, err = ob.specialHandleFeeRate()
		if err != nil {
			return errors.Wrap(err, "unable to execute specialHandleFeeRate")
		}
	} else {
		// EstimateSmartFee returns the fees per kilobyte (BTC/kb) targeting given block confirmation
		feeResult, err := ob.btcClient.EstimateSmartFee(1, &btcjson.EstimateModeEconomical)
		if err != nil {
			return errors.Wrap(err, "unable to estimate smart fee")
		}
		if feeResult.Errors != nil || feeResult.FeeRate == nil {
			return fmt.Errorf("error getting gas price: %s", feeResult.Errors)
		}
		if *feeResult.FeeRate > math.MaxInt64 {
			return fmt.Errorf("gas price is too large: %f", *feeResult.FeeRate)
		}
		feeRateEstimated = bitcoin.FeeRateToSatPerByte(*feeResult.FeeRate).Uint64()
	}

	// query the current block number
	blockNumber, err := ob.btcClient.GetBlockCount()
	if err != nil {
		return err
	}

	// #nosec G115 always positive
	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), feeRateEstimated, "100", uint64(blockNumber))
	if err != nil {
		ob.logger.GasPrice.Err(err).Msg("err PostGasPrice")
		return err
	}

	return nil
}

// GetSenderAddressByVin get the sender address from the previous transaction
// TODO(revamp): move in upper package to separate file (e.g., rpc.go)
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

	// #nosec G115 - always in range
	if len(tx.MsgTx().TxOut) <= int(vin.Vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vin.Vout, vin.Txid)
	}

	// decode sender address from previous pkScript
	pkScript := tx.MsgTx().TxOut[vin.Vout].PkScript
	scriptHex := hex.EncodeToString(pkScript)
	if bitcoin.IsPkScriptP2TR(pkScript) {
		return bitcoin.DecodeScriptP2TR(scriptHex, net)
	}
	if bitcoin.IsPkScriptP2WSH(pkScript) {
		return bitcoin.DecodeScriptP2WSH(scriptHex, net)
	}
	if bitcoin.IsPkScriptP2WPKH(pkScript) {
		return bitcoin.DecodeScriptP2WPKH(scriptHex, net)
	}
	if bitcoin.IsPkScriptP2SH(pkScript) {
		return bitcoin.DecodeScriptP2SH(scriptHex, net)
	}
	if bitcoin.IsPkScriptP2PKH(pkScript) {
		return bitcoin.DecodeScriptP2PKH(scriptHex, net)
	}

	// sender address not found, return nil and move on to the next tx
	return "", nil
}

// WatchUTXOs watches bitcoin chain for UTXOs owned by the TSS address
// TODO(revamp): move ticker related functions to a specific file
func (ob *Observer) WatchUTXOs(ctx context.Context) error {
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchUTXOs", ob.GetChainParams().WatchUtxoTicker)
	if err != nil {
		ob.logger.UTXOs.Error().Err(err).Msg("error creating ticker")
		return err
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.FetchUTXOs(ctx)
			if err != nil {
				ob.logger.UTXOs.Error().Err(err).Msg("error fetching btc utxos")
			}
			ticker.UpdateInterval(ob.GetChainParams().WatchUtxoTicker, ob.logger.UTXOs)
		case <-ob.StopChannel():
			ob.logger.UTXOs.Info().Msgf("WatchUTXOs stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// FetchUTXOs fetches TSS-owned UTXOs from the Bitcoin node
// TODO(revamp): move to UTXO file
func (ob *Observer) FetchUTXOs(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			ob.logger.UTXOs.Error().Msgf("BTC FetchUTXOs: caught panic error: %v", err)
		}
	}()

	// This is useful when a zetaclient's pending nonce lagged behind for whatever reason.
	ob.refreshPendingNonce(ctx)

	// get the current block height.
	bh, err := ob.btcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("btc: error getting block height : %v", err)
	}
	maxConfirmations := int(bh)

	// List all unspent UTXOs (160ms)
	tssAddr := ob.TSS().BTCAddress()
	address, err := chains.DecodeBtcAddress(tssAddr, ob.Chain().ChainId)
	if err != nil {
		return fmt.Errorf("btc: error decoding wallet address (%s) : %s", tssAddr, err.Error())
	}
	utxos, err := ob.btcClient.ListUnspentMinMaxAddresses(0, maxConfirmations, []btcutil.Address{address})
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
		if utxo.Amount < bitcoin.DefaultDepositorFee {
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

	ob.Mu().Lock()
	ob.TelemetryServer().SetNumberOfUTXOs(len(utxosFiltered))
	ob.utxos = utxosFiltered
	ob.Mu().Unlock()
	return nil
}

// SaveBroadcastedTx saves successfully broadcasted transaction
// TODO(revamp): move to db file
func (ob *Observer) SaveBroadcastedTx(txHash string, nonce uint64) {
	outboundID := ob.GetTxID(nonce)
	ob.Mu().Lock()
	ob.broadcastedTx[outboundID] = txHash
	ob.Mu().Unlock()

	broadcastEntry := clienttypes.ToOutboundHashSQLType(txHash, outboundID)
	if err := ob.DB().Save(&broadcastEntry).Error; err != nil {
		ob.logger.Outbound.Error().
			Err(err).
			Msgf("SaveBroadcastedTx: error saving broadcasted txHash %s for outbound %s", txHash, outboundID)
	}
	ob.logger.Outbound.Info().Msgf("SaveBroadcastedTx: saved broadcasted txHash %s for outbound %s", txHash, outboundID)
}

// GetBlockByNumberCached gets cached block (and header) by block number
func (ob *Observer) GetBlockByNumberCached(blockNumber int64) (*BTCBlockNHeader, error) {
	if result, ok := ob.BlockCache().Get(blockNumber); ok {
		if block, ok := result.(*BTCBlockNHeader); ok {
			return block, nil
		}
		return nil, errors.New("cached value is not of type *BTCBlockNHeader")
	}

	// Get the block hash
	hash, err := ob.btcClient.GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	// Get the block header
	header, err := ob.btcClient.GetBlockHeader(hash)
	if err != nil {
		return nil, err
	}
	// Get the block with verbose transactions
	block, err := ob.btcClient.GetBlockVerboseTx(hash)
	if err != nil {
		return nil, err
	}
	blockNheader := &BTCBlockNHeader{
		Header: header,
		Block:  block,
	}
	ob.BlockCache().Add(blockNumber, blockNheader)
	ob.BlockCache().Add(hash, blockNheader)
	return blockNheader, nil
}

// LoadDB open sql database and load data into Bitcoin observer
func (ob *Observer) LoadDB(dbPath string) error {
	if dbPath == "" {
		return errors.New("empty db path")
	}

	// open database, the custom dbName is used here for backward compatibility
	err := ob.OpenDB(dbPath, "btc_chain_client")
	if err != nil {
		return errors.Wrapf(err, "error OpenDB for chain %d", ob.Chain().ChainId)
	}

	// run auto migration
	// transaction result table is used nowhere but we still run migration in case they are needed in future
	err = ob.DB().AutoMigrate(
		&clienttypes.TransactionResultSQLType{},
		&clienttypes.OutboundHashSQLType{},
	)
	if err != nil {
		return errors.Wrapf(err, "error AutoMigrate for chain %d", ob.Chain().ChainId)
	}

	// load last scanned block
	err = ob.LoadLastBlockScanned()
	if err != nil {
		return err
	}

	// load broadcasted transactions
	err = ob.LoadBroadcastedTxMap()
	return err
}

// LoadLastBlockScanned loads the last scanned block from the database
func (ob *Observer) LoadLastBlockScanned() error {
	err := ob.Observer.LoadLastBlockScanned(ob.Logger().Chain)
	if err != nil {
		return errors.Wrapf(err, "error LoadLastBlockScanned for chain %d", ob.Chain().ChainId)
	}

	// observer will scan from the last block when 'lastBlockScanned == 0', this happens when:
	// 1. environment variable is set explicitly to "latest"
	// 2. environment variable is empty and last scanned block is not found in DB
	if ob.LastBlockScanned() == 0 {
		blockNumber, err := ob.btcClient.GetBlockCount()
		if err != nil {
			return errors.Wrapf(err, "error GetBlockCount for chain %d", ob.Chain().ChainId)
		}
		// #nosec G115 always positive
		ob.WithLastBlockScanned(uint64(blockNumber))
	}

	// bitcoin regtest starts from hardcoded block 100
	if chains.IsBitcoinRegnet(ob.Chain().ChainId) {
		ob.WithLastBlockScanned(RegnetStartBlock)
	}
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from block %d", ob.Chain().ChainId, ob.LastBlockScanned())

	return nil
}

// LoadBroadcastedTxMap loads broadcasted transactions from the database
func (ob *Observer) LoadBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutboundHashSQLType
	if err := ob.DB().Find(&broadcastedTransactions).Error; err != nil {
		ob.logger.Chain.Error().Err(err).Msgf("error iterating over db for chain %d", ob.Chain().ChainId)
		return err
	}
	for _, entry := range broadcastedTransactions {
		ob.broadcastedTx[entry.Key] = entry.Hash
	}
	return nil
}

// specialHandleFeeRate handles the fee rate for regnet and testnet
func (ob *Observer) specialHandleFeeRate() (uint64, error) {
	switch ob.Chain().NetworkType {
	case chains.NetworkType_privnet:
		// hardcode gas price for regnet
		return 1, nil
	case chains.NetworkType_testnet:
		feeRateEstimated, err := rpc.GetRecentFeeRate(ob.btcClient, ob.netParams)
		if err != nil {
			return 0, errors.Wrapf(err, "error GetRecentFeeRate")
		}
		return feeRateEstimated, nil
	default:
		return 0, fmt.Errorf(" unsupported bitcoin network type %d", ob.Chain().NetworkType)
	}
}

// isTssTransaction checks if a given transaction was sent by TSS itself.
// An unconfirmed transaction is safe to spend only if it was sent by TSS and verified by ourselves.
func (ob *Observer) isTssTransaction(txid string) bool {
	_, found := ob.includedTxHashes[txid]
	return found
}

// postBlockHeader posts block header to zetacore
// TODO(revamp): move to block header file
func (ob *Observer) postBlockHeader(ctx context.Context, tip int64) error {
	ob.logger.Inbound.Info().Msgf("postBlockHeader: tip %d", tip)
	bn := tip
	chainState, err := ob.ZetacoreClient().GetBlockHeaderChainState(ctx, ob.Chain().ChainId)
	if err == nil && chainState != nil && chainState.EarliestHeight > 0 {
		bn = chainState.LatestHeight + 1
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
		ob.logger.Inbound.Error().Err(err).Msgf("error serializing bitcoin block header: %d", bn)
		return err
	}
	blockHash := res2.Header.BlockHash()
	_, err = ob.ZetacoreClient().PostVoteBlockHeader(
		ctx,
		ob.Chain().ChainId,
		blockHash[:],
		res2.Block.Height,
		proofs.NewBitcoinHeader(headerBuf.Bytes()),
	)
	ob.logger.Inbound.Info().Msgf("posted block header %d: %s", bn, blockHash)
	if err != nil { // error shouldn't block the process
		ob.logger.Inbound.Error().Err(err).Msgf("error posting bitcoin block header: %d", bn)
	}
	return err
}
