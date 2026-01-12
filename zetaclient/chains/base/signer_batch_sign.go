package base

import (
	"bytes"
	"context"
	"encoding/hex"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// collectBatchBackoff is the backoff duration for retrying batch collection
	collectBatchBackoff = 4 * time.Second

	// collectBatchRetries is the maximum number of retries for batch collection
	collectBatchRetries = 2
)

// CheckBlockEvent checks if the block event is stale and returns (zeta_height, is_stale, error).
func (s *Signer) CheckBlockEvent(ctx context.Context, zetaRepo *zrepo.ZetaRepo) (int64, bool, error) {
	zetaBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return 0, false, errors.Wrap(err, "unable to get block event from context")
	}

	// get real-time zeta height
	zetaHeight, err := zetaRepo.GetBlockHeight(ctx)
	if err != nil {
		return 0, false, errors.Wrap(err, "unable to get zeta height")
	}

	// real-time zeta height are the signals to trigger TSS keysign on the exact same time,
	// so we need to ensure the block event is up to date (not a stale one accumulated in the channel)
	if zetaBlock.Block.Height < zetaHeight {
		s.Logger().
			Std.Info().
			Int64("zeta_height", zetaHeight).
			Int64("event_block", zetaBlock.Block.Height).
			Msg("stale block event")
		return zetaHeight, true, nil
	}

	// if the event is not stale, it applies the keysign delay configured in operational flags
	time.Sleep(delay)

	return zetaHeight, false, nil
}

// IsTimeToKeysign checks if it's time to perform keysign.
func (s *Signer) IsTimeToKeysign(
	p types.PendingNonces,
	nextTSSNonce uint64,
	zetaHeight int64,
	scheduleInterval int64,
) bool {
	// keysign happens only when zeta height is a multiple of the schedule interval
	if zetaHeight%scheduleInterval != 0 {
		return false
	}

	// return false if no pending cctx to sign
	if p.NonceLow >= p.NonceHigh {
		return false
	}

	// return false if TSS nonce is already ahead, it means that outbounds
	// were processed by external chain but don't have enough confirmations
	// #nosec G115 - always positive
	if nextTSSNonce >= uint64(p.NonceHigh) {
		return false
	}

	return true
}

// GetKeysignBatch returns the keysign batch to for given batch number.
func (s *Signer) GetKeysignBatch(ctx context.Context, zetaRepo *zrepo.ZetaRepo, batchNumber uint64) *TSSKeysignBatch {
	logger := s.Logger().Std.With().Uint64("batch_num", batchNumber).Logger()

	// return nil if batch number is not ready to sign
	ready, untilNonce, err := s.isBatchReadyToSign(ctx, zetaRepo, batchNumber)
	if err != nil {
		logger.Error().Err(err).Msg("unable to check batch readiness")
		return nil
	} else if !ready {
		return nil
	}

	// batch collector function
	var batch *TSSKeysignBatch
	fnCollect := func() error {
		batch, err = s.collectKeysignBatch(batchNumber, untilNonce)
		// force retry on error
		return retry.Retry(err)
	}

	// collect batch with retries
	bo := backoff.NewConstantBackOff(collectBatchBackoff)
	boWithMaxRetries := backoff.WithMaxRetries(bo, collectBatchRetries)
	if err := retry.DoWithBackoff(fnCollect, boWithMaxRetries); err != nil {
		logger.Error().Err(err).Msg("unable to collect keysign batch")
	}

	return batch
}

// SignBatch signs a batch of digests and adds the signatures to the cache.
func (s *Signer) SignBatch(ctx context.Context, batch TSSKeysignBatch, zetaHeight int64) error {
	var (
		chainID      = s.Chain().ChainId
		digests      = batch.Digests()
		keysignNonce = batch.NonceHigh()
		logger       = s.batchLogger(batch)
	)

	// calculate keysign height
	keysignHeight, err := KeysignHeight(chainID, zetaHeight)
	if err != nil {
		return errors.Wrap(err, "unable to calculate keysign height")
	}

	logger = logger.With().Int64("height", zetaHeight).Uint64("keysign_height", keysignHeight).Logger()
	logger.Info().Msg("signing batch of digests")

	// sign batch
	sigs, err := s.TSS().SignBatch(ctx, digests, keysignHeight, keysignNonce, chainID)
	if err != nil {
		logger.Error().Err(err).Msg("batch keysign failed")
		return err
	}
	logger.Info().Msg("signed batch of digests")

	// add signatures to cache
	s.AddBatchSignatures(batch, sigs)

	return nil
}

// GetSignatureOrAddDigest returns cached signature for given nonce and digest, or adds digest to cache if not found.
func (s *Signer) GetSignatureOrAddDigest(nonce uint64, digest []byte) ([65]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		batchNumber = NonceToBatchNumber(nonce)
		digestHex   = hex.EncodeToString(digest)
		logger      = s.Logger().
				Std.With().
				Uint64("batch_num", batchNumber).
				Uint64(logs.FieldNonce, nonce).
				Str("digest", digestHex).
				Logger()
	)

	info, found := s.tssKeysignInfoMap[nonce]
	if !found {
		s.tssKeysignInfoMap[nonce] = NewTSSKeysignInfo(digest, [65]byte{})
		logger.Info().Msg("added digest to cache")

		return [65]byte{}, false
	}

	// if digest has changed (e.g. increased gas price),
	// it means the signature is no longer valid. Update
	// digest and clear the old signature, then return false.
	if !bytes.Equal(info.digest, digest) {
		oldDigestHex := hex.EncodeToString(info.digest)
		info.digest = digest
		info.signature = [65]byte{}
		logger.Info().Str("old_digest", oldDigestHex).Msg("updated digest in cache")

		return [65]byte{}, false
	}

	return info.signature, info.signature != [65]byte{}
}

// IsBatchSigned returns true if the given batch was already signed before.
func (s *Signer) IsBatchSigned(batch *TSSKeysignBatch) bool {
	return s.TSS().IsSignatureCached(s.Chain().ChainId, batch.Digests())
}

// isBatchReadyToSign returns (true, untilNonce) if the given batch number is ready to sign.
// A batch is ready when it covers a range of nonces that overlaps with the pending nonces.
func (s *Signer) isBatchReadyToSign(
	ctx context.Context,
	zetaRepo *zrepo.ZetaRepo,
	batchNumber uint64,
) (bool, uint64, error) {
	p, err := zetaRepo.GetPendingNonces(ctx)
	if err != nil {
		return false, 0, errors.Wrapf(err, "unable to get pending nonces for chain %d", s.Chain().ChainId)
	}

	// prepare logger
	logger := s.Logger().
		Std.With().
		Uint64("batch_num", batchNumber).
		Int64("nonce_low", p.NonceLow).
		Int64("nonce_high", p.NonceHigh).
		Logger()

	// return false if no pending cctx
	if p.NonceLow >= p.NonceHigh {
		logger.Info().Msg("no pending cctx to sign")
		return false, 0, nil
	}

	// #nosec G115 - always positive
	cctxNonceLow, cctxNonceHigh := uint64(p.NonceLow), uint64(p.NonceHigh)-1
	batchNonceLow, batchNonceHigh := BatchNumberToRange(batchNumber)

	// calculate the overlap range
	overlap := cctxNonceLow <= batchNonceHigh && batchNonceLow <= cctxNonceHigh
	if !overlap {
		logger.Info().Msg("batch is not ready to sign")
		return false, 0, nil
	}
	untilNonce := min(cctxNonceHigh, batchNonceHigh)

	return true, untilNonce, nil
}

// collectKeysignBatch collects keysign batch for the given batch number until the given nonce.
func (s *Signer) collectKeysignBatch(batchNumber uint64, untilNonce uint64) (*TSSKeysignBatch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var (
		batch                         = NewTSSKeysignBatch()
		batchNonceLow, batchNonceHigh = BatchNumberToRange(batchNumber)
	)

	// collect digests within the batch's range
	for nonce := batchNonceLow; nonce <= batchNonceHigh; nonce++ {
		info, found := s.tssKeysignInfoMap[nonce]
		if found {
			batch.AddKeysignInfo(nonce, *info)
		}
	}

	switch {
	case batch.IsEmpty():
		return nil, errors.New("waiting for digests")
	case !batch.IsSequential():
		// if batch contains gaps, wait for the digests to be added to cache
		return nil, errors.New("waiting for digests gaps")
	case !batch.ContainsNonce(untilNonce):
		// wait for all nonces until the given nonce to be added to cache
		return nil, errors.Errorf("waiting for digests until nonce %d", untilNonce)
	default:
		return batch, nil
	}
}

// AddBatchSignatures adds TSS signatures to the cache.
func (s *Signer) AddBatchSignatures(batch TSSKeysignBatch, sigs [][65]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		nonceLow  = batch.NonceLow()
		nonceHigh = batch.NonceHigh()
		logger    = s.batchLogger(batch)
	)

	for nonce := nonceLow; nonce <= nonceHigh; nonce++ {
		info, found := s.tssKeysignInfoMap[nonce]
		if !found {
			continue
		}

		// TSS service ensures the order of signatures matches the order of digests
		sigIndex := nonce - nonceLow

		// skip signature if digest has changed (e.g. increased gas price)
		if !bytes.Equal(info.digest, batch.Digests()[sigIndex]) {
			logger.Info().Uint64(logs.FieldNonce, nonce).Msg("skipping signature")
			continue
		}

		// set signature
		info.signature = sigs[sigIndex]
		logger.Info().Uint64(logs.FieldNonce, nonce).Msg("added signature to cache")
	}
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

// batchLogger returns the logger for keysign batch.
func (s *Signer) batchLogger(batch TSSKeysignBatch) zerolog.Logger {
	return s.Logger().
		Std.With().
		Uint64("batch_num", batch.BatchNumber()).
		Uint64("nonce_low", batch.NonceLow()).
		Uint64("nonce_high", batch.NonceHigh()).
		Logger()
}
