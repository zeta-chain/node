package zetaclient

import (
	"fmt"
	"math"
	"os"
	"sort"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/rs/zerolog"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	RescanBatchSize      uint64 = 1000
	MaxRetryOnMissedCctx uint64 = 300 // 30 minutes
)

// CctxScanner scans missed pending cctx and updates their status
type CctxScanner struct {
	tssPubkey              string
	db                     *gorm.DB
	logger                 *zerolog.Logger
	bridge                 *ZetaCoreBridge
	firstNonceToScan       map[int64]uint64                                   // chainID -> the nonce to scan from when zetaclient starts
	nextNonceToScan        map[int64]uint64                                   // chainID -> next nonce to scan from when zetaclient is running
	missedPendingCctx      map[int64]map[uint64]*crosschaintypes.CrossChainTx // chainID -> nonce -> missed pending cctx
	missedPendingCctxRetry map[int64]map[uint64]uint64                        // chainID -> nonce -> retry count
}

func NewCctxScanner(bridge *ZetaCoreBridge, dbpath string, memDB bool, tssPubkey string, logger *zerolog.Logger) (*CctxScanner, error) {
	sc := &CctxScanner{
		logger:                 logger,
		bridge:                 bridge,
		firstNonceToScan:       make(map[int64]uint64),
		nextNonceToScan:        make(map[int64]uint64),
		missedPendingCctx:      make(map[int64]map[uint64]*crosschaintypes.CrossChainTx),
		missedPendingCctxRetry: make(map[int64]map[uint64]uint64),
	}
	err := sc.LoadDB(dbpath, memDB)
	if err != nil {
		return nil, err
	}

	// on bootstrap or tss migration
	if tssPubkey != sc.tssPubkey {
		err = sc.Reset(tssPubkey)
		if err != nil {
			return nil, err
		}
	}
	return sc, nil
}

// Scan a new batch of missed pending cctx
func (sc *CctxScanner) ScanMissedPendingCctx(bn int64, chainID int64, pendingNonces *crosschaintypes.PendingNonces) []*crosschaintypes.CrossChainTx {
	// initialize missed cctx map
	if _, found := sc.missedPendingCctx[chainID]; !found {
		sc.missedPendingCctx[chainID] = make(map[uint64]*crosschaintypes.CrossChainTx)
	}

	// calculate nonce range to scan
	nonceFrom, found := sc.nextNonceToScan[chainID]
	if !found { // uses db nonce on restart
		nonceFrom = sc.firstNonceToScan[chainID]
		sc.logger.Info().Msgf("scanner: chain %d starts from nonce %d", chainID, nonceFrom)
	}
	nonceTo := nonceFrom + RescanBatchSize
	if nonceTo > uint64(pendingNonces.NonceLow) {
		nonceTo = uint64(pendingNonces.NonceLow)
	}

	// scans [fromNonce, toNonce) for missed pending cctx
	if nonceFrom < nonceTo {
		missedList, err := sc.bridge.GetAllPendingCctxInNonceRange(chainID, nonceFrom, nonceTo)
		if err != nil {
			sc.logger.Error().Err(err).Msgf("scanner: failed to get pending cctx for chain %d from nonce %d to %d", chainID, nonceFrom, nonceTo)
			return sc.AllMissedPendingCctxByChain(chainID)
		}
		sc.addMissedPendingCctx(chainID, nonceFrom, nonceTo, missedList)
	}
	return sc.AllMissedPendingCctxByChain(chainID)
}

// Note: deep clone is unnecessary as the cctx list is used in a single thread
func (sc *CctxScanner) AllMissedPendingCctxByChain(chainID int64) []*crosschaintypes.CrossChainTx {
	missed := make([]*crosschaintypes.CrossChainTx, 0)
	for _, send := range sc.missedPendingCctx[chainID] {
		missed = append(missed, send)
	}
	sort.Slice(missed, func(i, j int) bool {
		return missed[i].GetCurrentOutTxParam().OutboundTxTssNonce < missed[j].GetCurrentOutTxParam().OutboundTxTssNonce
	})
	return missed
}

func (sc *CctxScanner) IsMissedPendingCctx(chainID int64, nonce uint64) bool {
	_, found := sc.missedPendingCctx[chainID][nonce]
	return found
}

// Re-check and update missed cctx's status
func (sc *CctxScanner) UpdateMissedPendingCctx(chainID int64, nonce uint64, nonceLow uint64) {
	send, err := sc.bridge.GetCctxByNonce(chainID, nonce)
	if err != nil {
		sc.logger.Error().Err(err).Msgf("scanner: error GetCctxByNonce for chain %d nonce %d", chainID, nonce)
		return
	}
	if crosschaintypes.IsCctxStatusPending(send.CctxStatus.Status) {
		// increase retry count
		if _, found := sc.missedPendingCctxRetry[chainID]; !found {
			sc.missedPendingCctxRetry[chainID] = make(map[uint64]uint64)
		}
		sc.missedPendingCctxRetry[chainID][nonce]++

		// forget about this missed cctx as max retry reached, there could be 2 cases:
		//    1. No tracker.   For some reason, no observer had reported outtx hash.
		//    2. No true hash. Track exists but none of the hashes is true (can't be verified)
		if sc.missedPendingCctxRetry[chainID][nonce] == MaxRetryOnMissedCctx {
			sc.removeMissedPendingCctx(chainID, nonce, nonceLow)
			sc.logger.Warn().Msgf("scanner: forget about missed pending cctx for chain %d nonce %d", chainID, nonce)
		}
	} else { // no longer pending
		sc.removeMissedPendingCctx(chainID, nonce, nonceLow)
		sc.logger.Info().Msgf("scanner: removed missed pending cctx for chain %d nonce %d", chainID, nonce)
	}
}

func (sc *CctxScanner) addMissedPendingCctx(chainID int64, nonceFrom uint64, nonceTo uint64, missedList []*crosschaintypes.CrossChainTx) {
	nonces := make([]uint64, 0)
	for _, send := range missedList {
		nonce := send.GetCurrentOutTxParam().OutboundTxTssNonce
		nonces = append(nonces, nonce)
		sc.missedPendingCctx[chainID][nonce] = send
	}
	sc.nextNonceToScan[chainID] = nonceTo
	if len(nonces) > 0 {
		sc.logger.Info().Msgf("scanner: found missed pending cctx for chain %d with nonces %v", chainID, nonces)
	}
	sc.saveFirstNonceToScan(chainID)
}

func (sc *CctxScanner) removeMissedPendingCctx(chainID int64, nonce uint64, nonceLow uint64) {
	delete(sc.missedPendingCctx[chainID], nonce)
	sc.saveFirstNonceToScan(chainID)
}

func (sc *CctxScanner) saveFirstNonceToScan(chainID int64) {
	firstNonceToScan := uint64(math.MaxUint64)
	if len(sc.missedPendingCctx[chainID]) > 0 { // save the lowest nonce for future restart if there ARE missed pending cctx
		for nonceMissed := range sc.missedPendingCctx[chainID] {
			if nonceMissed < firstNonceToScan {
				firstNonceToScan = nonceMissed
			}
		}
	} else { // save the 'NextNonceToScan' for future restart if no missed pending cctx found so far
		firstNonceToScan = sc.nextNonceToScan[chainID]
	}
	if firstNonceToScan < uint64(math.MaxUint64) {
		if err := sc.db.Save(clienttypes.ToFirstNonceToScanSQLType(chainID, firstNonceToScan)).Error; err != nil {
			sc.logger.Error().Err(err).Msgf("scanner: error writing firstNonceToScan for chain %d nonce %d", chainID, firstNonceToScan)
		}
	}
}

// LoadDB open sql database and load data into scanner
func (sc *CctxScanner) LoadDB(dbpath string, memDB bool) error {
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		err := os.MkdirAll(dbpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	path := dbpath
	if !memDB { // memDB is used for uint test only
		path = fmt.Sprintf("%s/scanner", dbpath)
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	sc.db = db

	err = db.AutoMigrate(&clienttypes.CurrentTssSQLType{},
		&clienttypes.FirstNonceToScanSQLType{})
	if err != nil {
		return err
	}

	// Load current tss pubkey
	sc.loadCurrentTssPubkey()

	// Load first nonce for each chain to start scanning from
	err = sc.buildFirstNonceToScanMap()

	return err
}

func (sc *CctxScanner) Reset(tssPubkey string) error {
	sc.tssPubkey = tssPubkey
	sc.firstNonceToScan = make(map[int64]uint64)
	sc.nextNonceToScan = make(map[int64]uint64)
	sc.missedPendingCctx = make(map[int64]map[uint64]*crosschaintypes.CrossChainTx)
	sc.missedPendingCctxRetry = make(map[int64]map[uint64]uint64)

	// save current tss pubkey
	if err := sc.db.Save(clienttypes.ToCurrentTssSQLType(tssPubkey)).Error; err != nil {
		sc.logger.Error().Err(err).Msgf("scanner: error writing current tss pubkey %s", tssPubkey)
		return err
	}

	// clean db, GORM uses pluralizes struct name to snake_cases as table name
	if err := sc.db.Exec("DELETE FROM first_nonce_to_scan_sql_types").Error; err != nil {
		sc.logger.Error().Err(err).Msg("scanner: error cleaning FirstNonceToScan db")
		return err
	}
	sc.logger.Info().Msgf("scanner: reset db successfully for tss pubkey %s", tssPubkey)

	return nil
}

func (sc *CctxScanner) loadCurrentTssPubkey() {
	var tss clienttypes.CurrentTssSQLType
	if err := sc.db.First(&tss, clienttypes.CurrentTssID).Error; err != nil {
		sc.tssPubkey = "" // record not found
	}
	sc.tssPubkey = tss.TssPubkey
}

func (sc *CctxScanner) buildFirstNonceToScanMap() error {
	var firstNonces []clienttypes.FirstNonceToScanSQLType
	if err := sc.db.Find(&firstNonces).Error; err != nil {
		sc.logger.Error().Err(err).Msg("scanner: error iterating over FirstNonceToScan db")
		return err
	}
	for _, nonce := range firstNonces {
		sc.firstNonceToScan[nonce.ID] = nonce.FirstNonce
		sc.nextNonceToScan[nonce.ID] = nonce.FirstNonce
		sc.logger.Info().Msgf("scanner: first nonce to scan for chain %d is %d", nonce.ID, nonce.FirstNonce)
	}
	return nil
}
