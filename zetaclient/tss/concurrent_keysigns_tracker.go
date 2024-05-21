package tss

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

// ConcurrentKeysignsTracker keeps track of concurrent keysigns performed by go-tss
type ConcurrentKeysignsTracker struct {
	numActiveMsgSigns int64
	mu                sync.Mutex
	Logger            zerolog.Logger
}

// NewKeysignsTracker - constructor
func NewKeysignsTracker(logger zerolog.Logger) *ConcurrentKeysignsTracker {
	return &ConcurrentKeysignsTracker{
		numActiveMsgSigns: 0,
		mu:                sync.Mutex{},
		Logger:            logger.With().Str("submodule", "ConcurrentKeysignsTracker").Logger(),
	}
}

// StartMsgSign is incrementing the number of active signing ceremonies as well as updating the prometheus metric
func (k *ConcurrentKeysignsTracker) StartMsgSign() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.numActiveMsgSigns++
	metrics.NumActiveMsgSigns.Inc()
	k.Logger.Debug().Msgf("Start TSS message sign, numActiveMsgSigns: %d", k.numActiveMsgSigns)
}

// EndMsgSign is decrementing the number of active signing ceremonies as well as updating the prometheus metric
func (k *ConcurrentKeysignsTracker) EndMsgSign() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.numActiveMsgSigns > 0 {
		k.numActiveMsgSigns--
		metrics.NumActiveMsgSigns.Dec()
	}
	k.Logger.Debug().Msgf("End TSS message sign, numActiveMsgSigns: %d", k.numActiveMsgSigns)
}

// GetNumActiveMessageSigns gets the current number of active signing ceremonies
func (k *ConcurrentKeysignsTracker) GetNumActiveMessageSigns() int64 {
	return k.numActiveMsgSigns
}
