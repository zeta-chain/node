package outboundprocessor

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Processor struct {
	outboundStartTime  map[string]time.Time
	outboundEndTime    map[string]time.Time
	outboundActive     map[string]struct{}
	mu                 sync.Mutex
	Logger             zerolog.Logger
	numActiveProcessor int64
}

func NewOutboundProcessorManager(logger zerolog.Logger) *Processor {
	return &Processor{
		outboundStartTime:  make(map[string]time.Time),
		outboundEndTime:    make(map[string]time.Time),
		outboundActive:     make(map[string]struct{}),
		mu:                 sync.Mutex{},
		Logger:             logger.With().Str("module", "OutboundProcessorManager").Logger(),
		numActiveProcessor: 0,
	}
}

func (outboundManager *Processor) StartTryProcess(outboundID string) {
	outboundManager.mu.Lock()
	defer outboundManager.mu.Unlock()
	outboundManager.outboundStartTime[outboundID] = time.Now()
	outboundManager.outboundActive[outboundID] = struct{}{}
	outboundManager.numActiveProcessor++
	outboundManager.Logger.Info().Msgf("StartTryProcess %s, numActiveProcessor %d", outboundID, outboundManager.numActiveProcessor)
}

func (outboundManager *Processor) EndTryProcess(outboundID string) {
	outboundManager.mu.Lock()
	defer outboundManager.mu.Unlock()
	outboundManager.outboundEndTime[outboundID] = time.Now()
	delete(outboundManager.outboundActive, outboundID)
	outboundManager.numActiveProcessor--
	outboundManager.Logger.Info().Msgf("EndTryProcess %s, numActiveProcessor %d, time elapsed %s", outboundID, outboundManager.numActiveProcessor, time.Since(outboundManager.outboundStartTime[outboundID]))
}

func (outboundManager *Processor) IsOutboundActive(outboundID string) bool {
	outboundManager.mu.Lock()
	defer outboundManager.mu.Unlock()
	_, found := outboundManager.outboundActive[outboundID]
	return found
}

func (outboundManager *Processor) TimeInTryProcess(outboundID string) time.Duration {
	outboundManager.mu.Lock()
	defer outboundManager.mu.Unlock()
	if _, found := outboundManager.outboundActive[outboundID]; found {
		return time.Since(outboundManager.outboundStartTime[outboundID])
	}
	return 0
}

// ToOutboundID returns the outboundID for OutboundProcessorManager to track
func ToOutboundID(index string, receiverChainID int64, nonce uint64) string {
	return fmt.Sprintf("%s-%d-%d", index, receiverChainID, nonce)
}
