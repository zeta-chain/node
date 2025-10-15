package base

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/mode"
)

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

	// mu protects fields from concurrent access
	// Note: base signer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu sync.Mutex

	ClientMode mode.ClientMode
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

		s.logger.Std.Info().
			Bool("outbound_active", active).
			Str(logs.FieldOutboundID, outboundID).
			Time("outbound_timestamp", now).
			Int("outbound_total", len(s.activeOutbounds)).
			Msg("MarkOutbound")
	default:
		timeTaken := time.Since(startedAt)

		s.logger.Std.Info().
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
