package types

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type DynamicTicker struct {
	name     string
	interval uint64
	impl     *time.Ticker
}

func NewDynamicTicker(name string, interval uint64) (*DynamicTicker, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("non-positive ticker interval %d for %s", interval, name)
	}

	return &DynamicTicker{
		name:     name,
		interval: interval,
		impl:     time.NewTicker(time.Duration(interval) * time.Second),
	}, nil
}

func (t *DynamicTicker) C() <-chan time.Time {
	return t.impl.C
}

func (t *DynamicTicker) UpdateInterval(newInterval uint64, logger zerolog.Logger) {
	if newInterval > 0 && t.interval != newInterval {
		t.impl.Stop()
		oldInterval := t.interval
		t.interval = newInterval
		t.impl = time.NewTicker(time.Duration(t.interval) * time.Second)
		logger.Info().Msgf("%s ticker interval changed from %d to %d", t.name, oldInterval, newInterval)
	}
}

func (t *DynamicTicker) Stop() {
	t.impl.Stop()
}
