package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	cometbft "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/ticker"
)

// intervalTicker wrapper for ticker.Ticker.
type intervalTicker struct {
	ticker *ticker.Ticker
}

func newIntervalTicker(
	task Executable,
	interval time.Duration,
	intervalUpdater func() time.Duration,
	taskName string,
	logger zerolog.Logger,
) *intervalTicker {
	wrapper := func(ctx context.Context, t *ticker.Ticker) error {
		if err := task(ctx); err != nil {
			logger.Error().Err(err).Msgf("Task %s failed", taskName)
		}

		if intervalUpdater != nil {
			// noop if interval is not changed
			t.SetInterval(normalizeInterval(intervalUpdater()))
		}

		return nil
	}

	tt := ticker.New(interval, wrapper, ticker.WithLogger(logger, taskName))

	return &intervalTicker{ticker: tt}
}

func (t *intervalTicker) Start(ctx context.Context) error {
	return t.ticker.Start(ctx)
}

func (t *intervalTicker) Stop() {
	t.ticker.StopBlocking()
}

// blockTicker represents custom ticker implementation that ticks on new Zeta block events.
// Pass blockTicker ONLY by pointer.
type blockTicker struct {
	exec Executable

	// block channel that will be used to receive new blocks
	blockChan <-chan cometbft.EventDataNewBlock

	// stopChan is used to stop the ticker
	stopChan chan struct{}

	// doneChan is used to signal that the ticker has stopped (i.e. "blocking stop")
	doneChan chan struct{}

	isRunning bool
	mu        sync.Mutex

	logger zerolog.Logger
}

func newBlockTicker(task Executable, blockChan <-chan cometbft.EventDataNewBlock, logger zerolog.Logger) *blockTicker {
	return &blockTicker{
		exec:      task,
		blockChan: blockChan,
		logger:    logger,
	}
}

func (t *blockTicker) Start(ctx context.Context) error {
	if err := t.init(); err != nil {
		return err
	}

	defer t.cleanup()

	// release Stop() blocking
	defer func() { close(t.doneChan) }()

	for {
		select {
		case block, ok := <-t.blockChan:
			// channel closed
			if !ok {
				t.logger.Warn().Msg("Block channel closed")
				return nil
			}

			ctx := WithBlockEvent(ctx, block)

			if err := t.exec(ctx); err != nil {
				t.logger.Error().Err(err).Msg("Task error")
			}
		case <-ctx.Done():
			t.logger.Warn().Err(ctx.Err()).Msg("Content error")
			return nil
		case <-t.stopChan:
			// caller invoked t.stop()
			return nil
		}
	}
}

func (t *blockTicker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// noop
	if !t.isRunning {
		return
	}

	// notify async loop to stop
	close(t.stopChan)

	// wait for the loop to stop
	<-t.doneChan

	t.isRunning = false
}

func (t *blockTicker) init() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isRunning {
		return fmt.Errorf("ticker already started")
	}

	t.stopChan = make(chan struct{})
	t.doneChan = make(chan struct{})
	t.isRunning = true

	return nil
}

// if ticker was stopped NOT by Stop() method, we want to make a cleanup
func (t *blockTicker) cleanup() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// noop
	if !t.isRunning {
		return
	}

	t.isRunning = false
	close(t.stopChan)
}
