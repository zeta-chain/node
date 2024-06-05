package base

import (
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

// Signer is the base chain signer for external chains
type Signer struct {
	// the external chain
	chain chains.Chain

	// zetacore context
	zetacoreContext *context.ZetacoreContext

	// tss signer
	tss interfaces.TSSSigner

	// telemetry server
	ts *metrics.TelemetryServer

	// the standard logger
	logger zerolog.Logger

	// the compliance logger
	loggerCompliance zerolog.Logger
}

// NewSigner creates a new base signer
func NewSigner(
	chain chains.Chain,
	zetacoreContext *context.ZetacoreContext,
	tss interfaces.TSSSigner,
	logger zerolog.Logger,
	loggerCompliance zerolog.Logger,
	ts *metrics.TelemetryServer,
) *Signer {
	return &Signer{
		chain:            chain,
		zetacoreContext:  zetacoreContext,
		tss:              tss,
		logger:           logger,
		loggerCompliance: loggerCompliance,
		ts:               ts,
	}
}
