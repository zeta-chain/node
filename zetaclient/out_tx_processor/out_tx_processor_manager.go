package out_tx_processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type OutTxProcessorManager struct {
	outTxStartTime     map[string]time.Time
	outTxEndTime       map[string]time.Time
	outTxActive        map[string]struct{}
	mu                 sync.Mutex
	Logger             zerolog.Logger
	numActiveProcessor int64
}

func NewOutTxProcessorManager(logger zerolog.Logger) *OutTxProcessorManager {
	return &OutTxProcessorManager{
		outTxStartTime:     make(map[string]time.Time),
		outTxEndTime:       make(map[string]time.Time),
		outTxActive:        make(map[string]struct{}),
		mu:                 sync.Mutex{},
		Logger:             logger.With().Str("module", "OutTxProcessorManager").Logger(),
		numActiveProcessor: 0,
	}
}

func (outTxMan *OutTxProcessorManager) StartTryProcess(outTxID string) {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	outTxMan.outTxStartTime[outTxID] = time.Now()
	outTxMan.outTxActive[outTxID] = struct{}{}
	outTxMan.numActiveProcessor++
	outTxMan.Logger.Info().Msgf("StartTryProcess %s, numActiveProcessor %d", outTxID, outTxMan.numActiveProcessor)
}

func (outTxMan *OutTxProcessorManager) EndTryProcess(outTxID string) {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	outTxMan.outTxEndTime[outTxID] = time.Now()
	delete(outTxMan.outTxActive, outTxID)
	outTxMan.numActiveProcessor--
	outTxMan.Logger.Info().Msgf("EndTryProcess %s, numActiveProcessor %d, time elapsed %s", outTxID, outTxMan.numActiveProcessor, time.Since(outTxMan.outTxStartTime[outTxID]))
}

func (outTxMan *OutTxProcessorManager) IsOutTxActive(outTxID string) bool {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	_, found := outTxMan.outTxActive[outTxID]
	return found
}

func (outTxMan *OutTxProcessorManager) TimeInTryProcess(outTxID string) time.Duration {
	outTxMan.mu.Lock()
	defer outTxMan.mu.Unlock()
	if _, found := outTxMan.outTxActive[outTxID]; found {
		return time.Since(outTxMan.outTxStartTime[outTxID])
	}
	return 0
}

// ToOutTxID returns the outTxID for OutTxProcessorManager to track
func ToOutTxID(index string, receiverChainID int64, nonce uint64) string {
	return fmt.Sprintf("%s-%d-%d", index, receiverChainID, nonce)
}
