package base

import (
	"context"
	"slices"

	ethcommon "github.com/ethereum/go-ethereum/common"

	mathpkg "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// batchSize is the number of digests in a keysign batch
	// signing a 10-digest batch takes about 3~4 seconds on average
	batchSize = 10
)

// tssKeysignInfo contains one TSS keysign information
type tssKeysignInfo struct {
	cctxHeight uint64
	digest     ethcommon.Hash
	signature  [65]byte
}

// TSSKeysignBatch contains a batch of TSS keysign information
type TSSKeysignBatch struct {
	digests    [][]byte
	nonceLow   uint64
	nonceHigh  uint64
	heightLow  uint64
	heightHigh uint64
}

func newtssKeysignBatch() TSSKeysignBatch {
	return TSSKeysignBatch{
		digests: make([][]byte, 0),
	}
}

// addKeysignInfo adds one TSS keysign info to the batch and updates the nonce and height
func (b *TSSKeysignBatch) addKeysignInfo(nonce uint64, info tssKeysignInfo) {
	b.digests = append(b.digests, info.digest.Bytes())

	// initialize on first call
	if len(b.digests) == 1 {
		b.nonceLow = nonce
		b.nonceHigh = nonce
		b.heightLow = info.cctxHeight
		b.heightHigh = info.cctxHeight
		return
	}

	// update nonceLow and heightLow
	if nonce < b.nonceLow {
		b.nonceLow = nonce
	} else if nonce > b.nonceHigh {
		b.nonceHigh = nonce
	}

	// update heightLow and heightHigh
	if info.cctxHeight > b.heightLow {
		b.heightLow = info.cctxHeight
	} else if info.cctxHeight > b.heightHigh {
		b.heightHigh = info.cctxHeight
	}
}

// KeysignHeight calculates an artificial keysign height tweaked with chainID
func (b *TSSKeysignBatch) KeysignHeight(chainID int64, height uint64) uint64 {
	// #nosec G115 e2eTest - always in range
	zetaHeight32 := uint32(height)

	// #nosec G115 e2eTest - always in range
	chainID32 := uint32(chainID)

	return mathpkg.CantorPair(zetaHeight32, chainID32)
}

// BatchNumber returns the batch number of the keysign batch
func (b *TSSKeysignBatch) BatchNumber() uint64 {
	return NonceToBatchNumber(b.nonceLow)
}

// NonceToBatchNumber returns the batch number for a given nonce.
// For example: nonce 1 falls into batch 0, nonce 10 falls into batch 1, ...
func NonceToBatchNumber(nonce uint64) uint64 {
	return nonce / batchSize
}

func (s *Signer) GetSignature(nonce uint64) ([65]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, found := s.tssKeysignInfoMap[nonce]
	if !found {
		return [65]byte{}, false
	}

	return info.signature, info.signature != [65]byte{}
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

// SignBatch signs a batch of digests and adds the signatures to the cache.
func (s *Signer) SignBatch(ctx context.Context, batch TSSKeysignBatch, zetaHeight int64) error {
	var (
		chainID      = s.Chain().ChainId
		digests      = batch.digests
		keysignNonce = batch.nonceHigh
		batchNumber  = batch.BatchNumber()

		// it's an artificial height to uniquely identify the batch; added 1 to avoid 0 height
		keysignHeight = batchNumber + 1
	)

	logger := s.Logger().
		Std.With().
		Int64("height", zetaHeight).
		Uint64("batch", batchNumber).
		Int("digests", len(digests)).
		Logger()

	sigs, err := s.TSS().SignBatch(ctx, digests, keysignHeight, keysignNonce, chainID)
	if err != nil {
		logger.Error().Err(err).Msg("batch keysign failed")
		return err
	}

	s.AddTSSSignatures(batch, sigs)

	s.SetSignedBatch(batchNumber, true)

	logger.Info().Msg("signed batch of digests")

	return nil
}

// AddKeysignInfo adds TSS keysign info to the cache.
func (s *Signer) AddKeysignInfo(nonce, cctxHeight uint64, digest ethcommon.Hash) {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger := s.Logger().Std.With().Uint64(logs.FieldNonce, nonce).Logger()

	info, found := s.tssKeysignInfoMap[nonce]
	if !found {
		s.tssKeysignInfoMap[nonce] = &tssKeysignInfo{
			digest:     digest,
			signature:  [65]byte{},
			cctxHeight: cctxHeight,
		}

		logger.Info().Msg("added digest to cache")
		return
	}

	// update digest if it has changed (e.g. gas price update)
	if info.digest != digest {
		info.digest = digest
		logger.Info().Msg("updated digest in cache")
	}
}

func (s *Signer) AddTSSSignatures(batch TSSKeysignBatch, signatures [][65]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		nonceLow  = batch.nonceLow
		nonceHigh = batch.nonceHigh
	)

	for nonce := nonceLow; nonce <= nonceHigh; nonce++ {
		info, found := s.tssKeysignInfoMap[nonce]
		if found {
			sigIndex := nonce - nonceLow
			info.signature = signatures[sigIndex]
		}
	}
}

// GetKeysignBatch returns the keysign batch to sign for given batchNonce.
func (s *Signer) GetKeysignBatch(batchNumber uint64) *TSSKeysignBatch {
	logger := s.Logger().Std.With().Uint64("batch_num", batchNumber).Logger()

	if s.IsBatchSigned(batchNumber) {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// sort all allNonces in ascending order
	allNonces := make([]uint64, 0, len(s.tssKeysignInfoMap))
	for nonce := range s.tssKeysignInfoMap {
		allNonces = append(allNonces, nonce)
	}
	slices.Sort(allNonces)

	var (
		// example batch ranges are: [0, 9], [10, 19], [20, 29], [30, 39], ...
		batchNonceLow  = batchNumber * batchSize
		batchNonceHigh = (batchNumber+1)*batchSize - 1

		keysignBatch = newtssKeysignBatch()
	)

	// collect digests for the keysign batch
	for _, nonce := range allNonces {
		if nonce < batchNonceLow {
			continue
		} else if nonce > batchNonceHigh {
			break
		}

		// early return if any one of the digests in range is not found
		info, found := s.tssKeysignInfoMap[nonce]
		if !found {
			logger.Info().Msgf("digest not found for nonce %d", nonce)
			return nil
		}
		keysignBatch.addKeysignInfo(nonce, *info)
	}

	// TODO: remove this check and allow partial batches
	if len(keysignBatch.digests) != batchSize {
		logger.Info().Int("digests_count", len(keysignBatch.digests)).Msg("waiting for full batch")
		return nil
	}

	return &keysignBatch
}
