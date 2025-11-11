package base

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	mathpkg "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/mode"
)

const (
	// blocksPerBatch is the number of blocks to group digests into a batch
	blocksPerBatch = 100

	// maxBatchSize is the maximum number of digests in a batch to sign
	// signing a big batch takes longer time and may cause delay and timeout
	maxBatchSize = 10

	// recentBlocks is the number of blocks to consider recent for TSS keysign
	// recent blocks
	recentBlocks = 15
)

type tssKeysignInfo struct {
	cctxHeight uint64
	zetaHeight uint64
	digest     ethcommon.Hash
	signature  [65]byte
}

type TSSKeysignBatch struct {
	digests    [][]byte
	nonceLow   uint64
	nonceHigh  uint64
	heightLow  uint64
	heightHigh uint64
}

func newKeysignBatch(nonceLow uint64, heightLow uint64) TSSKeysignBatch {
	return TSSKeysignBatch{
		digests:    make([][]byte, 0),
		nonceLow:   nonceLow,
		nonceHigh:  nonceLow,
		heightLow:  heightLow,
		heightHigh: heightLow,
	}
}

// addKeysignInfo adds a keysign info to the batch and updates the nonce and height
func (b *TSSKeysignBatch) addKeysignInfo(nonce uint64, info tssKeysignInfo) {
	b.digests = append(b.digests, info.digest.Bytes())

	if nonce > b.nonceHigh {
		b.nonceHigh = nonce
	}

	if info.cctxHeight > b.heightHigh {
		b.heightHigh = info.cctxHeight
	}
}

// Digests returns the digests in the batch
func (b *TSSKeysignBatch) Digests() [][]byte {
	return b.digests
}

// NonceLow returns the nonceLow of the batch
func (b *TSSKeysignBatch) NonceLow() uint64 {
	return b.nonceLow
}

// NonceHigh returns the nonceHigh of the batch
func (b *TSSKeysignBatch) NonceHigh() uint64 {
	return b.nonceHigh
}

// KeysignHeight calculates an artificial keysign height (based on current Zeta height) that uniquely identifies the batch
func (b *TSSKeysignBatch) KeysignHeight(zetaHeight uint64) uint64 {
	// #nosec G115 e2eTest - always in range
	zetaHeight32 := uint32(zetaHeight)

	// #nosec G115 e2eTest - always in range
	uniqueNonce32 := uint32(b.KeysignNonce())

	return mathpkg.CantorPair(zetaHeight32, uniqueNonce32)
}

// KeysignNonce returns the nonceLow as the identifier of the keysign batch
func (b *TSSKeysignBatch) KeysignNonce() uint64 {
	return b.nonceLow
}

// BatchNonce returns the batch number of the keysign batch
func (b *TSSKeysignBatch) BatchNonce() uint64 {
	return b.nonceLow/maxBatchSize + 1
}

// Signer is the base structure for grouping the common logic between chain signers.
// The common logic includes: chain, chainParams, contexts, tss, metrics, loggers etc.
type Signer struct {
	// chain contains static information about the external chain
	chain chains.Chain

	// tssSigner is the TSS signer
	tssSigner tssrepo.TSSClient

	// logger contains the loggers used by signer
	logger Logger

	// outboundBeingReported is a map of outbound being reported to tracker
	outboundBeingReported map[string]bool

	activeOutbounds map[string]time.Time

	// tssKeysignInfoMap maps nonce to TSS keysign information to be signed
	tssKeysignInfoMap map[uint64]*tssKeysignInfo

	// current batch nonce is the low nonce of the current batch of digests to sign
	currentBatchNonce uint64

	// activeBatchNumbers is a map of batch numbers being signed
	activeBatchNumbers map[uint64]bool

	// signedBatchNumbers is a map of batch numbers that have been signed
	signedBatchNumbers map[uint64]bool

	// mu protects fields from concurrent access
	// Note: base signer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu sync.RWMutex

	ClientMode mode.ClientMode
}

// SignBatch signs a batch of digests and adds the signatures to the cache.
func (s *Signer) SignBatch(ctx context.Context, batch TSSKeysignBatch, zetaHeight int64, critical bool) error {
	var (
		chainID      = s.Chain().ChainId
		digests      = batch.Digests()
		keysignNonce = batch.KeysignNonce()
		batchNumber  = batch.BatchNonce()
		//keysignHeight = batch.KeysignHeight(zetaHeight)
	)

	// for non-critical batches, we allow only one single batch keysign at the same time
	if !critical {
		// try to set batch as active
		wasActive := s.SetActiveBatch(keysignNonce)
		if wasActive {
			s.Logger().Std.Info().Uint64("batch_nonce", keysignNonce).Msg("batch is active, skipping")
			return nil
		}

		// clear active batch flag on exit
		defer func() {
			s.ClearActiveBatch(keysignNonce)
		}()
	}

	logger := s.Logger().
		Std.With().
		Uint64("batch_nonce", batch.BatchNonce()).
		Int64("zeta_height", zetaHeight).
		Logger()

	if critical {
		logger.Info().Msg("signing batch critical of digests")
	} else {
		logger.Info().Msg("signing batch of digests")
	}

	sigs, err := s.TSS().SignBatch(ctx, digests, batchNumber, keysignNonce, chainID)
	if err != nil {
		logger.Error().Err(err).Msgf("batch keysign failed batch_nonce=%d", batch.BatchNonce())
		return err
	}

	s.SetSignedBatch(batch.BatchNonce(), true)

	s.AddTSSSignatures(sigs, batch.NonceLow(), batch.NonceHigh())

	if critical {
		logger.Info().Msg("signed batch critical of digests")
	} else {
		logger.Info().Msg("signed batch of digests")
	}

	return nil
}

// AddKeysignInfo adds TSS keysign info to the cache.
func (s *Signer) AddKeysignInfo(nonce, zetaHeight, cctxHeight uint64, digest ethcommon.Hash) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if info, found := s.tssKeysignInfoMap[nonce]; found {
		if zetaHeight > info.zetaHeight {
			info.zetaHeight = zetaHeight
			s.logger.Std.Debug().
				Uint64(logs.FieldNonce, nonce).
				Uint64(logs.FieldBlock, zetaHeight).
				Msg("tx height has increased")

			if info.digest != digest {
				s.logger.Std.Info().Uint64(logs.FieldNonce, nonce).Msg("tx digest has changed")
			}
		}
		return
	}

	s.logger.Std.Debug().
		Uint64(logs.FieldNonce, nonce).
		Uint64(logs.FieldBlock, zetaHeight).
		Msg("added tx digest to sign")

	s.tssKeysignInfoMap[nonce] = &tssKeysignInfo{
		digest:     digest,
		signature:  [65]byte{},
		zetaHeight: zetaHeight,
		cctxHeight: cctxHeight,
	}
}

func (s *Signer) AddTSSSignatures(signatures [][65]byte, nonceLow, nonceHigh uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for nonce := nonceLow; nonce <= nonceHigh; nonce++ {
		txSig, found := s.tssKeysignInfoMap[nonce]
		if found {
			sigIndex := nonce - nonceLow
			txSig.signature = signatures[sigIndex]
		}
	}

	s.currentBatchNonce = nonceHigh + 1
	s.Logger().Std.Debug().Msgf("current batch nonce updated to %d", s.currentBatchNonce)
}

func (s *Signer) GetSignature(nonce uint64) ([65]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	txSig, found := s.tssKeysignInfoMap[nonce]
	if !found {
		return [65]byte{}, false
	}

	return txSig.signature, txSig.signature != [65]byte{}
}

func (s *Signer) IsBatchActive(batchNumber uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeBatchNumbers[batchNumber]
}

func (s *Signer) SetActiveBatch(batchNumber uint64) (wasActive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	wasActive = s.activeBatchNumbers[batchNumber]
	s.activeBatchNumbers[batchNumber] = true

	return wasActive
}

func (s *Signer) ClearActiveBatch(batchNumber uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.activeBatchNumbers, batchNumber)
}

func (s *Signer) SetSignedBatch(batchNumber uint64, signed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signedBatchNumbers[batchNumber] = signed
}

func (s *Signer) IsBatchSigned(batchNumber uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.signedBatchNumbers[batchNumber]
}

// GetDigestBatches returns a list of digest batches partitioned by keysign height
func (s *Signer) GetDigestBatches() []TSSKeysignBatch {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.tssKeysignInfoMap) == 0 {
		return nil
	}

	// sort all allNonces in ascending order
	allNonces := make([]uint64, 0, len(s.tssKeysignInfoMap))
	for nonce := range s.tssKeysignInfoMap {
		allNonces = append(allNonces, nonce)
	}
	slices.Sort(allNonces)

	// wait for digests of all nonces to be present before partitioning (no gaps)
	// this is to make sure each batch of digests are contiguous and deterministic
	lowestNonce := allNonces[0]
	highestNonce := allNonces[len(allNonces)-1]
	// #nosec G115 always positive
	if uint64(len(allNonces)) != highestNonce-lowestNonce+1 {
		s.logger.Std.Info().Msg("waiting for digests of all nonces to partition")
		return nil
	}

	// sanity check
	// when sorted by nonce in ascending order, the cctx height should also be in ascending order
	// this error should NEVER happen but just in case, we log and return nil
	prevHeight := s.tssKeysignInfoMap[allNonces[0]].cctxHeight
	for i := 1; i < len(allNonces); i++ {
		thisHeight := s.tssKeysignInfoMap[allNonces[i]].cctxHeight
		if thisHeight < prevHeight {
			s.logger.Std.Error().Uint64(logs.FieldNonce, allNonces[i]).Msg("cctx carries wrong height")
			return nil
		}
		prevHeight = thisHeight
	}

	var (
		// the batches to return
		batches = make([]TSSKeysignBatch, 0)

		// the current batch being collected
		lowestHeight = s.tssKeysignInfoMap[lowestNonce].cctxHeight
		currentBatch = newKeysignBatch(lowestNonce, lowestHeight)

		// round nonce up to nearest multiple of maxBatchSize, indicating the upper bound of each batch
		// example upper bounds are: 19, 39, 59, 79, 99, ...
		// corresponding ranges are: [0, 19], [20, 39], [40, 59], [60, 79], [80, 99], ...
		currentBatchUpperBound = (lowestNonce/maxBatchSize+1)*maxBatchSize - 1
	)

	// group digests into batches by nonce range
	for _, nonce := range allNonces {
		info := s.tssKeysignInfoMap[nonce]

		if nonce <= currentBatchUpperBound {
			currentBatch.addKeysignInfo(nonce, *info)
		} else {
			// end current batch when nonce exceeds current batch upper bound
			batches = append(batches, currentBatch)

			// start new batch
			currentBatchUpperBound += maxBatchSize
			//s.logger.Std.Info().Msgf("started new batch: %d, %d, %d, %d", len(batches), len(currentBatch.Digests()), currentBatchUpperBound, nonce)
			currentBatch = newKeysignBatch(nonce, info.cctxHeight)
			currentBatch.addKeysignInfo(nonce, *info)
		}
	}

	// add last batch
	if len(currentBatch.Digests()) > 0 && len(currentBatch.Digests()) == maxBatchSize {
		batches = append(batches, currentBatch)
		s.logger.Std.Info().Msgf("added last batch: %d, %d", len(batches), len(currentBatch.Digests()))
	}

	// // if last batch is not full, wait
	// if len(batches[len(batches)-1].Digests()) < maxBatchSize {
	// 	s.logger.Std.Info().Msgf("waiting for last batch to be full: %d, %d", len(batches), len(batches[len(batches)-1].Digests()))
	// 	return nil
	// }

	// if last batch is not full, sign each digest individually

	return batches
}

func (s *Signer) GetNextDigestBatch() ([][]byte, uint64, uint64, uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// sort all nonces in ascending order
	nonces := make([]uint64, 0, len(s.tssKeysignInfoMap))
	for nonce := range s.tssKeysignInfoMap {
		nonces = append(nonces, nonce)
	}
	slices.Sort(nonces)

	// early return if less than batch size nonces to sign
	if len(nonces) > 0 && nonces[len(nonces)-1] < s.currentBatchNonce+maxBatchSize-1 {
		s.Logger().
			Std.Info().
			Msgf("last nonce %d < %d, total %d", nonces[len(nonces)-1], s.currentBatchNonce+maxBatchSize, len(nonces))
		return nil, 0, 0, 0
	}

	// take next batch of digests to sign
	var (
		maxHeight = uint64(0)
		nonceLow  = s.currentBatchNonce
		nonceHigh = nonceLow + maxBatchSize - 1
		digests   = make([][]byte, 0, maxBatchSize)
	)

	for nonce := nonceLow; nonce <= nonceHigh; nonce++ {
		// early return if any one of the digests in range is not found
		info, found := s.tssKeysignInfoMap[nonce]
		if !found {
			return nil, 0, 0, 0
		}

		// if no signature available, add digest to batch
		if info.signature == [65]byte{} {
			nonceLow = min(nonce, nonceLow)
			nonceHigh = max(nonce, nonceHigh)
			maxHeight = max(info.zetaHeight, maxHeight)
			digests = append(digests, info.digest.Bytes())
		}
	}

	// round maxHeight up to nearest multiple of 10
	batchHeight := (maxHeight + 4) / 5 * 5

	return digests, batchHeight, nonceLow, nonceHigh
}

// NewSigner creates a new base signer.
func NewSigner(
	chain chains.Chain,
	tssSigner tssrepo.TSSClient,
	logger Logger,
	clientMode mode.ClientMode,
) *Signer {
	withLogFields := func(log zerolog.Logger) zerolog.Logger {
		return log.With().
			Str(logs.FieldModule, logs.ModNameSigner).
			Int64(logs.FieldChain, chain.ChainId).
			Logger()
	}

	return &Signer{
		chain:                 chain,
		tssSigner:             tssSigner,
		outboundBeingReported: make(map[string]bool),
		activeOutbounds:       make(map[string]time.Time),
		tssKeysignInfoMap:     make(map[uint64]*tssKeysignInfo),
		currentBatchNonce:     0,
		activeBatchNumbers:    make(map[uint64]bool),
		signedBatchNumbers:    make(map[uint64]bool),
		logger: Logger{
			Std:        withLogFields(logger.Std),
			Compliance: withLogFields(logger.Compliance),
		},
		ClientMode: clientMode,
	}
}

// Chain returns the chain for the signer.
func (s *Signer) Chain() chains.Chain {
	return s.chain
}

// TSS returns the tss signer for the signer.
func (s *Signer) TSS() tssrepo.TSSClient {
	return s.tssSigner
}

// Logger returns the logger for the signer.
func (s *Signer) Logger() *Logger {
	return &s.logger
}

// SetBeingReportedFlag sets the outbound as being reported if not already set.
// Returns true if the outbound is already being reported.
// This method is used by outbound tracker reporter to avoid repeated reporting of same hash.
func (s *Signer) SetBeingReportedFlag(hash string) (alreadySet bool) {
	s.Lock()
	defer s.Unlock()

	alreadySet = s.outboundBeingReported[hash]
	if !alreadySet {
		// mark as being reported
		s.outboundBeingReported[hash] = true
	}
	return
}

// ClearBeingReportedFlag clears the being reported flag for the outbound.
func (s *Signer) ClearBeingReportedFlag(hash string) {
	s.Lock()
	defer s.Unlock()
	delete(s.outboundBeingReported, hash)
}

// Exported for unit tests

// GetReportedTxList returns a list of outboundHash being reported.
// TODO: investigate pointer usage
// https://github.com/zeta-chain/node/issues/2084
func (s *Signer) GetReportedTxList() *map[string]bool {
	return &s.outboundBeingReported
}

// Lock locks the signer.
func (s *Signer) Lock() {
	s.mu.Lock()
}

// Unlock unlocks the signer.
func (s *Signer) Unlock() {
	s.mu.Unlock()
}

// MarkOutbound marks the outbound as active.
func (s *Signer) MarkOutbound(outboundID string, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	startedAt, found := s.activeOutbounds[outboundID]

	switch {
	case active == found:
		// no-op
	case active:
		now := time.Now().UTC()
		s.activeOutbounds[outboundID] = now

		s.logger.Std.Debug().
			Bool("outbound_active", active).
			Str(logs.FieldOutboundID, outboundID).
			Time("outbound_timestamp", now).
			Int("outbound_total", len(s.activeOutbounds)).
			Msg("MarkOutbound")
	default:
		timeTaken := time.Since(startedAt)

		s.logger.Std.Debug().
			Bool("outbound_active", active).
			Str(logs.FieldOutboundID, outboundID).
			Float64("outbound_time_taken", timeTaken.Seconds()).
			Int("outbound_total", len(s.activeOutbounds)).
			Msg("MarkOutbound")

		delete(s.activeOutbounds, outboundID)
	}
}

func (s *Signer) IsOutboundActive(outboundID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.activeOutbounds[outboundID]
	return found
}

// PassesCompliance checks if the cctx passes the compliance check and prints compliance log.
func (s *Signer) PassesCompliance(cctx *types.CrossChainTx) bool {
	if !compliance.IsCCTXRestricted(cctx) {
		return true
	}

	params := cctx.GetCurrentOutboundParam()

	compliance.PrintComplianceLog(
		s.Logger().Std,
		s.Logger().Compliance,
		true,
		s.Chain().ChainId,
		cctx.Index,
		cctx.InboundParams.Sender,
		params.Receiver,
		&params.CoinType,
	)

	return false
}

// OutboundID returns the outbound ID.
func OutboundID(index string, receiverChainID int64, nonce uint64) string {
	return fmt.Sprintf("%s-%d-%d", index, receiverChainID, nonce)
}

// OutboundIDFromCCTX returns the outbound ID from the cctx.
func OutboundIDFromCCTX(cctx *types.CrossChainTx) string {
	index, params := cctx.GetIndex(), cctx.GetCurrentOutboundParam()
	return OutboundID(index, params.ReceiverChainId, params.TssNonce)
}
