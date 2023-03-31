package zetaclient

import (
	"github.com/rs/zerolog"
	"sync"
	"syscall"
	"time"
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

// suicide whole zetaclient if keysign appears deadlocked.
func (outTxMan *OutTxProcessorManager) StartMonitorHealth() {
	outTxMan.Logger.Info().Msgf("StartMonitorHealth")
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		count := 0
		for outTxID := range outTxMan.outTxActive {
			if outTxMan.TimeInTryProcess(outTxID).Minutes() > 2 {
				count++
			}
		}
		if count > 0 {
			outTxMan.Logger.Warn().Msgf("Health: %d OutTx are more than 2min in process!", count)
		} else {
			outTxMan.Logger.Info().Msgf("Monitor: healthy; numActiveProcessor %d", outTxMan.numActiveProcessor)
		}
		if count > 10 {
			// suicide:
			outTxMan.Logger.Error().Msgf("suicide zetaclient because keysign appears deadlocked; kill this process and the process supervisor should restart it")
			outTxMan.Logger.Info().Msgf("numActiveProcessor: %d", outTxMan.numActiveProcessor)
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}
}
