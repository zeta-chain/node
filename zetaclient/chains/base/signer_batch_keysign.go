package base

import (
	"bytes"
	"context"
	"slices"

	mathpkg "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// batchSize is the number of digests in a keysign batch
	// signing a 10-digest batch takes about 3~4 seconds on average
	batchSize = 10
)

// TSSKeysignInfo represents a record of TSS keysign information.
type TSSKeysignInfo struct {
	// cctxHeight is the zeta block height when the cctx was created
	cctxHeight uint64

	// digest is the digest of the outbound transaction
	digest []byte

	// signature is the TSS signature of the digest
	signature [65]byte
}

// TSSKeysignBatch contains a batch of TSS keysign information.
type TSSKeysignBatch struct {
	// digests is a list of digests to sign
	digests [][]byte

	// nonceLow is the lowest cctx nonce in the batch
	nonceLow uint64

	// nonceHigh is the highest cctx nonce in the batch
	nonceHigh uint64

	// heightLow is the lowest cctx height in the batch
	heightLow uint64

	// heightHigh is the highest cctx height in the batch
	heightHigh uint64
}

// newTSSKeysignBatch creates a new TSS keysign batch.
func newTSSKeysignBatch() TSSKeysignBatch {
	return TSSKeysignBatch{
		digests: make([][]byte, 0),
	}
}

// addKeysignInfo adds one TSS keysign info to the batch and updates the nonce and height.
func (b *TSSKeysignBatch) addKeysignInfo(nonce uint64, info TSSKeysignInfo) {
	b.digests = append(b.digests, info.digest)

	// initialize on first record
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

// KeysignHeight calculates an artificial keysign height tweaked with chainID.
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

// SetSignedFlag sets the given batch number as signed.
func (s *Signer) SetSignedFlag(batchNumber uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signedBatchNumbers[batchNumber] = true
}

func (s *Signer) IsBatchSigned(batchNumber uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.signedBatchNumbers[batchNumber]
}

// GetKeysignBatch returns the keysign batch to for given batch number.
func (s *Signer) GetKeysignBatch(batchNumber uint64) *TSSKeysignBatch {
	if s.IsBatchSigned(batchNumber) {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	logger := s.Logger().Std.With().Uint64("batch_num", batchNumber).Logger()

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

		keysignBatch = newTSSKeysignBatch()
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
		logger.Info().Int("digests", len(keysignBatch.digests)).Msg("waiting for full batch")
		return nil
	}

	return &keysignBatch
}

// RemoveKeysignInfo removes keysign info for all nonces before the given nonce.
// This function is used to clean up stale keysign info in the cache.
func (s *Signer) RemoveKeysignInfo(beforeNonce uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for nonce := range s.tssKeysignInfoMap {
		if nonce < beforeNonce {
			delete(s.tssKeysignInfoMap, nonce)
		}
	}
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
		Uint64("batch_num", batchNumber).
		Int("digests", len(digests)).
		Logger()

	sigs, err := s.TSS().SignBatch(ctx, digests, keysignHeight, keysignNonce, chainID)
	if err != nil {
		logger.Error().Err(err).Msg("batch keysign failed")
		return err
	}

	s.AddBatchSignatures(batch, sigs)

	s.SetSignedFlag(batchNumber)

	logger.Info().Msg("signed batch of digests")

	return nil
}

// GetSignatureOrAddDigest returns cached signature for given nonce and digest, or adds digest to cache if not found.
func (s *Signer) GetSignatureOrAddDigest(nonce uint64, cctxHeight uint64, digest []byte) ([65]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		batchNumber = NonceToBatchNumber(nonce)
		logger      = s.Logger().Std.With().Uint64(logs.FieldNonce, nonce).Uint64("batch_num", batchNumber).Logger()
	)

	info, found := s.tssKeysignInfoMap[nonce]
	if !found {
		s.tssKeysignInfoMap[nonce] = &TSSKeysignInfo{
			digest:     digest,
			signature:  [65]byte{},
			cctxHeight: cctxHeight,
		}
		logger.Info().Msg("added digest to cache")

		return [65]byte{}, false
	}

	// if digest has changed (e.g. increased gas price),
	// it means the signature is no longer valid. Update
	// digest and mark the batch as unsigned, return false
	if !bytes.Equal(info.digest, digest) {
		info.digest = digest
		delete(s.signedBatchNumbers, batchNumber)
		logger.Info().Msg("updated digest in cache")

		return [65]byte{}, false
	}

	return info.signature, info.signature != [65]byte{}
}

// AddBatchSignatures adds TSS signatures to the cache.
func (s *Signer) AddBatchSignatures(batch TSSKeysignBatch, sigs [][65]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		nonceLow  = batch.nonceLow
		nonceHigh = batch.nonceHigh

		batchNumber = batch.BatchNumber()
		logger      = s.Logger().Std.With().Uint64("batch_num", batchNumber).Logger()
	)

	for nonce := nonceLow; nonce <= nonceHigh; nonce++ {
		info, found := s.tssKeysignInfoMap[nonce]
		if found {
			if info.signature == [65]byte{} {
				logger.Info().Uint64(logs.FieldNonce, nonce).Msg("add signature to cache")
			} else {
				logger.Info().Uint64(logs.FieldNonce, nonce).Msg("update signature in cache")
			}

			sigIndex := nonce - nonceLow
			info.signature = sigs[sigIndex]
		}
	}
}

// NonceToBatchNumber maps a nonce to a batch number.
// For example:
// - nonce 1 falls into batch 0
// - nonce 10 falls into batch 1
// - nonce 19 falls into batch 1
// - nonce 20 falls into batch 2
// - ...
func NonceToBatchNumber(nonce uint64) uint64 {
	return nonce / batchSize
}
