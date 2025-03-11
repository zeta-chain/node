package base

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// DefaultTSSSignatureCacheSize is the default number of signatures that the signer will keep in cache.
	// Caching 200 recent transactions signatures is good enough for most chains because zetaclients don't
	// look ahead for more than 200 outbound transactions.
	DefaultTSSSignatureCacheSize = 200

	// DefaultTSSSignatureExpiration is the default expiration time for cached TSS signatures.
	DefaultTSSSignatureExpiration = time.Minute * 30
)

// Signer is the base structure for grouping the common logic between chain signers.
// The common logic includes: chain, chainParams, contexts, tss, metrics, loggers etc.
type Signer struct {
	// chain contains static information about the external chain
	chain chains.Chain

	// tss is the TSS signer
	tss interfaces.TSSSigner

	// logger contains the loggers used by signer
	logger Logger

	// outboundBeingReported is a map of outbound being reported to tracker
	outboundBeingReported map[string]bool

	activeOutbounds map[string]time.Time

	// tssSignatureCache stores cached TSS signatures
	tssSignatureCache *TSSSignatureCache

	// mu protects fields from concurrent access
	// Note: base signer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu sync.Mutex
}

// NewSigner creates a new base signer.
func NewSigner(
	chain chains.Chain,
	tss interfaces.TSSSigner,
	sigCacheSize int,
	sigExpiration time.Duration,
	logger Logger,
) (*Signer, error) {
	withLogFields := func(log zerolog.Logger) zerolog.Logger {
		return log.With().
			Int64(logs.FieldChain, chain.ChainId).
			Str(logs.FieldModule, "signer").
			Logger()
	}

	tssSignatureCache, err := NewTSSSignatureCache(sigCacheSize, sigExpiration)
	if err != nil {
		return nil, errors.Wrap(err, "error creating tss signature cache")
	}

	return &Signer{
		chain:                 chain,
		tss:                   tss,
		outboundBeingReported: make(map[string]bool),
		activeOutbounds:       make(map[string]time.Time),
		tssSignatureCache:     tssSignatureCache,
		logger: Logger{
			Std:        withLogFields(logger.Std),
			Compliance: withLogFields(logger.Compliance),
		},
	}, nil
}

// Chain returns the chain for the signer.
func (s *Signer) Chain() chains.Chain {
	return s.chain
}

// TSS returns the tss signer for the signer.
func (s *Signer) TSS() interfaces.TSSSigner {
	return s.tss
}

// Logger returns the logger for the signer.
func (s *Signer) Logger() *Logger {
	return &s.logger
}

// TSSSign signs a given digest with TSS.
func (s *Signer) TSSSign(ctx context.Context, digest []byte, height, nonce uint64) (sig65B [65]byte, err error) {
	// get cached signature if available
	pkBech32 := s.tss.PubKey().Bech32String()
	sig65B, found := s.tssSignatureCache.Get(pkBech32, digest)
	if found {
		return sig65B, nil
	}

	// sign the digest with TSS
	sig65B, err = s.TSS().Sign(ctx, digest, height, nonce, s.Chain().ChainId)
	if err != nil {
		return [65]byte{}, errors.Wrap(err, "tss Sign failed")
	}

	// add signature to the cache
	s.tssSignatureCache.Add(pkBech32, digest, sig65B)
	s.Logger().Std.Info().
		Str(logs.FieldMethod, "Sign").
		Str("tss_addr", s.TSS().PubKey().AddressEVM().Hex()).
		Uint64("height", height).
		Uint64(logs.FieldNonce, nonce).
		Str("digest", hex.EncodeToString(digest)).
		Str("signature", hex.EncodeToString(sig65B[:])).
		Msg("add new tss signature to cache")

	return sig65B, nil
}

// TSSSignBatch signs a batch of digests with TSS.
func (s *Signer) TSSSignBatch(
	ctx context.Context,
	digests [][]byte,
	height, nonce uint64,
) (sig65Bs [][65]byte, err error) {
	// get cached signatures if available
	pkBech32 := s.tss.PubKey().Bech32String()
	sig65Bs, found := s.tssSignatureCache.GetBatch(pkBech32, digests)
	if found {
		return sig65Bs, nil
	}

	// sign the digests with TSS
	sig65Bs, err = s.TSS().SignBatch(ctx, digests, height, nonce, s.Chain().ChainId)
	if err != nil {
		return nil, errors.Wrap(err, "tss SignBatch failed")
	}

	// add signatures to the cache
	err = s.tssSignatureCache.AddBatch(pkBech32, digests, sig65Bs)
	if err != nil {
		return nil, err
	}
	s.Logger().Std.Info().
		Str(logs.FieldMethod, "SignBatch").
		Str("tss_addr", s.TSS().PubKey().AddressEVM().Hex()).
		Uint64("height", height).
		Uint64(logs.FieldNonce, nonce).
		Int("batch_size", len(digests)).
		Msg("add new tss signatures to cache")

	return sig65Bs, nil
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
		// noop
	case active:
		now := time.Now().UTC()
		s.activeOutbounds[outboundID] = now

		s.logger.Std.Info().
			Bool("outbound.active", active).
			Str("outbound.id", outboundID).
			Time("outbound.timestamp", now).
			Int("outbound.total", len(s.activeOutbounds)).
			Msg("MarkOutbound")
	default:
		timeTaken := time.Since(startedAt)

		s.logger.Std.Info().
			Bool("outbound.active", active).
			Str("outbound.id", outboundID).
			Float64("outbound.time_taken", timeTaken.Seconds()).
			Int("outbound.total", len(s.activeOutbounds)).
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

// OutboundID returns the outbound ID.
func OutboundID(index string, receiverChainID int64, nonce uint64) string {
	return fmt.Sprintf("%s-%d-%d", index, receiverChainID, nonce)
}

// OutboundIDFromCCTX returns the outbound ID from the cctx.
func OutboundIDFromCCTX(cctx *types.CrossChainTx) string {
	index, params := cctx.GetIndex(), cctx.GetCurrentOutboundParam()
	return OutboundID(index, params.ReceiverChainId, params.TssNonce)
}
