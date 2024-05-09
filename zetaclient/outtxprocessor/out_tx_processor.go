package outtxprocessor

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Processor struct {
	outTxStartTime     map[string]time.Time
	outTxEndTime       map[string]time.Time
	outTxActive        map[string]struct{}
	mu                 sync.Mutex
	Logger             zerolog.Logger
	numActiveProcessor int64
}

func NewProcessor(logger zerolog.Logger) *Processor {
	return &Processor{
		outTxStartTime:     make(map[string]time.Time),
		outTxEndTime:       make(map[string]time.Time),
		outTxActive:        make(map[string]struct{}),
		mu:                 sync.Mutex{},
		Logger:             logger.With().Str("module", "OutTxProcessorManager").Logger(),
		numActiveProcessor: 0,
	}
}

func (outTxProc *Processor) StartTryProcess(outTxID string) {
	outTxProc.mu.Lock()
	defer outTxProc.mu.Unlock()
	outTxProc.outTxStartTime[outTxID] = time.Now()
	outTxProc.outTxActive[outTxID] = struct{}{}
	outTxProc.numActiveProcessor++
	outTxProc.Logger.Info().Msgf("StartTryProcess %s, numActiveProcessor %d", outTxID, outTxProc.numActiveProcessor)
}

func (outTxProc *Processor) EndTryProcess(outTxID string) {
	outTxProc.mu.Lock()
	defer outTxProc.mu.Unlock()
	outTxProc.outTxEndTime[outTxID] = time.Now()
	delete(outTxProc.outTxActive, outTxID)
	outTxProc.numActiveProcessor--
	outTxProc.Logger.Info().Msgf("EndTryProcess %s, numActiveProcessor %d, time elapsed %s", outTxID, outTxProc.numActiveProcessor, time.Since(outTxProc.outTxStartTime[outTxID]))
}

func (outTxProc *Processor) IsOutTxActive(outTxID string) bool {
	outTxProc.mu.Lock()
	defer outTxProc.mu.Unlock()
	_, found := outTxProc.outTxActive[outTxID]
	return found
}

func (outTxProc *Processor) TimeInTryProcess(outTxID string) time.Duration {
	outTxProc.mu.Lock()
	defer outTxProc.mu.Unlock()
	if _, found := outTxProc.outTxActive[outTxID]; found {
		return time.Since(outTxProc.outTxStartTime[outTxID])
	}
	return 0
}

// ToOutTxID returns the outTxID for OutTxProcessorManager to track
func ToOutTxID(index string, receiverChainID int64, nonce uint64) string {
	return fmt.Sprintf("%s-%d-%d", index, receiverChainID, nonce)
}
