package tss

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

type KeySignManager struct {
	numActiveMsgSigns int64
	mu                sync.Mutex
	Logger            zerolog.Logger
}

func NewKeySignManager(logger zerolog.Logger) *KeySignManager {
	return &KeySignManager{
		numActiveMsgSigns: 0,
		mu:                sync.Mutex{},
		Logger:            logger.With().Str("module", "KeySignManager").Logger(),
	}
}

func (k *KeySignManager) StartMsgSign() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.numActiveMsgSigns++
	metrics.NumActiveMsgSigns.Inc()
	k.Logger.Debug().Msgf("Start TSS message sign, numActiveMsgSigns: %d", k.numActiveMsgSigns)
}

func (k *KeySignManager) EndMsgSign() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.numActiveMsgSigns > 0 {
		k.numActiveMsgSigns--
	}
	metrics.NumActiveMsgSigns.Dec()
	k.Logger.Debug().Msgf("End TSS message sign, numActiveMsgSigns: %d", k.numActiveMsgSigns)
}

func (k *KeySignManager) GetNumActiveMessageSigns() int64 {
	return k.numActiveMsgSigns
}
