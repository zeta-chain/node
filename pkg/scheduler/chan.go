package scheduler

import (
	"context"
	"fmt"
	"sync/atomic"

	cometbft "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"
)

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

	isRunning atomic.Bool

	logger zerolog.Logger
}

type blockCtxKey struct{}

func newBlockTicker(task Executable, blockChan <-chan cometbft.EventDataNewBlock, logger zerolog.Logger) *blockTicker {
	return &blockTicker{
		exec:      task,
		blockChan: blockChan,
		stopChan:  make(chan struct{}),
		doneChan:  nil,
		logger:    logger,
	}
}

func withBlockEvent(ctx context.Context, event cometbft.EventDataNewBlock) context.Context {
	return context.WithValue(ctx, blockCtxKey{}, event)
}

// BlockFromContext returns cometbft.EventDataNewBlock from the context or false.
func BlockFromContext(ctx context.Context) (cometbft.EventDataNewBlock, bool) {
	blockEvent, ok := ctx.Value(blockCtxKey{}).(cometbft.EventDataNewBlock)
	return blockEvent, ok
}

func (t *blockTicker) Start(ctx context.Context) error {
	if !t.setRunning(true) {
		return fmt.Errorf("ticker already started")
	}

	t.doneChan = make(chan struct{})
	defer func() {
		close(t.doneChan)

		// closes stopChan if it's not closed yet
		if t.setRunning(false) {
			close(t.stopChan)
		}
	}()

	for {
		select {
		case block, ok := <-t.blockChan:
			// channel closed
			if !ok {
				t.logger.Warn().Msg("Block channel closed")
				return nil
			}

			ctx := withBlockEvent(ctx, block)

			if err := t.exec(ctx); err != nil {
				t.logger.Warn().Err(err).Msg("Task error")
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
	// noop
	if !t.isRunning.Load() {
		return
	}

	// notify async loop to stop
	close(t.stopChan)

	// wait for the loop to stop
	<-t.doneChan
	t.setRunning(false)
}

func (t *blockTicker) setRunning(running bool) (changed bool) {
	if running {
		return t.isRunning.CompareAndSwap(false, true)
	}

	return t.isRunning.CompareAndSwap(true, false)
}
