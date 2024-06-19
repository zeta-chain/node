package base

import (
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

// Signer is the base structure for grouping the common logic between chain signers.
// The common logic includes: chain, chainParams, contexts, tss, metrics, loggers etc.
type Signer struct {
	// chain contains static information about the external chain
	chain chains.Chain

	// zetacoreContext is the Zetacore client to interact with ZetaChain
	zetacoreContext *context.ZetacoreContext

	// tss is the TSS signer
	tss interfaces.TSSSigner

	// ts is the telemetry server for metrics
	ts *metrics.TelemetryServer

	// logger contains the loggers used by signer
	logger Logger
}

// NewSigner creates a new base signer
func NewSigner(
	chain chains.Chain,
	zetacoreContext *context.ZetacoreContext,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
	logger Logger,
) *Signer {
	return &Signer{
		chain:           chain,
		zetacoreContext: zetacoreContext,
		tss:             tss,
		ts:              ts,
		logger: Logger{
			Std:        logger.Std.With().Int64("chain", chain.ChainId).Str("module", "signer").Logger(),
			Compliance: logger.Compliance,
		},
	}
}

// Chain returns the chain for the signer
func (s *Signer) Chain() chains.Chain {
	return s.chain
}

// WithChain attaches a new chain to the signer
func (s *Signer) WithChain(chain chains.Chain) *Signer {
	s.chain = chain
	return s
}

// ZetacoreContext returns the zetacore context for the signer
func (s *Signer) ZetacoreContext() *context.ZetacoreContext {
	return s.zetacoreContext
}

// WithZetacoreContext attaches a new zetacore context to the signer
func (s *Signer) WithZetacoreContext(context *context.ZetacoreContext) *Signer {
	s.zetacoreContext = context
	return s
}

// Tss returns the tss signer for the signer
func (s *Signer) TSS() interfaces.TSSSigner {
	return s.tss
}

// WithTSS attaches a new tss signer to the signer
func (s *Signer) WithTSS(tss interfaces.TSSSigner) *Signer {
	s.tss = tss
	return s
}

// TelemetryServer returns the telemetry server for the signer
func (s *Signer) TelemetryServer() *metrics.TelemetryServer {
	return s.ts
}

// WithTelemetryServer attaches a new telemetry server to the signer
func (s *Signer) WithTelemetryServer(ts *metrics.TelemetryServer) *Signer {
	s.ts = ts
	return s
}

// Logger returns the logger for the signer
func (s *Signer) Logger() *Logger {
	return &s.logger
}
