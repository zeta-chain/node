package base

import (
	"sync"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Signer is the base structure for grouping the common logic between chain signers.
// The common logic includes: chain, chainParams, contexts, tss, metrics, loggers etc.
type Signer struct {
	// chain contains static information about the external chain
	chain chains.Chain

	// tss is the TSS signer
	tss interfaces.TSSSigner

	// ts is the telemetry server for metrics
	ts *metrics.TelemetryServer

	// logger contains the loggers used by signer
	logger Logger

	// outboundBeingReported is a map of outbound being reported to tracker
	outboundBeingReported map[string]bool

	// mu protects fields from concurrent access
	// Note: base signer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu sync.Mutex
}

// NewSigner creates a new base signer.
func NewSigner(chain chains.Chain, tss interfaces.TSSSigner, ts *metrics.TelemetryServer, logger Logger) *Signer {
	return &Signer{
		chain: chain,
		tss:   tss,
		ts:    ts,
		logger: Logger{
			Std: logger.Std.With().
				Int64(logs.FieldChain, chain.ChainId).
				Str(logs.FieldModule, "signer").
				Logger(),
			Compliance: logger.Compliance,
		},
		outboundBeingReported: make(map[string]bool),
	}
}

// Chain returns the chain for the signer.
func (s *Signer) Chain() chains.Chain {
	return s.chain
}

// WithChain attaches a new chain to the signer.
func (s *Signer) WithChain(chain chains.Chain) *Signer {
	s.chain = chain
	return s
}

// Tss returns the tss signer for the signer.
func (s *Signer) TSS() interfaces.TSSSigner {
	return s.tss
}

// WithTSS attaches a new tss signer to the signer.
func (s *Signer) WithTSS(tss interfaces.TSSSigner) *Signer {
	s.tss = tss
	return s
}

// TelemetryServer returns the telemetry server for the signer.
func (s *Signer) TelemetryServer() *metrics.TelemetryServer {
	return s.ts
}

// WithTelemetryServer attaches a new telemetry server to the signer.
func (s *Signer) WithTelemetryServer(ts *metrics.TelemetryServer) *Signer {
	s.ts = ts
	return s
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
