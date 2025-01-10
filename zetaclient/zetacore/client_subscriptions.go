package zetacore

import (
	"context"

	"cosmossdk.io/errors"
	ctypes "github.com/cometbft/cometbft/types"

	"github.com/zeta-chain/node/pkg/fanout"
)

// NewBlockSubscriber subscribes to comet bft new block events.
// Subscribes share the same websocket connection but their channels are independent (fanout)
func (c *Client) NewBlockSubscriber(ctx context.Context) (chan ctypes.EventDataNewBlock, error) {
	blockSubscriber, err := c.resolveBlockSubscriber()
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve block subscriber")
	}

	// we need a "proxy" chan instead of directly returning blockSubscriber.Add()
	// to support context cancellation
	blocksChan := make(chan ctypes.EventDataNewBlock)

	go func() {
		consumer := blockSubscriber.Add()

		for {
			select {
			case <-ctx.Done():
				// fixme: MEMORY LEAK: this might be dangerous because the consumer is not closed.
				// Fanout will spawn "zombie" goroutines to push to the chan, but nobody is reading from it,
				// Will be addressed in future orchestrator V2 PRs (not urgent as of now)
				return
			case block := <-consumer:
				blocksChan <- block
			}
		}
	}()

	return blocksChan, nil
}

// resolveBlockSubscriber returns the block subscriber channel
// or subscribes to it for the first time.
func (c *Client) resolveBlockSubscriber() (*fanout.FanOut[ctypes.EventDataNewBlock], error) {
	// noop
	if blocksFanout := c.blockFanOutThreadSafe(); blocksFanout != nil {
		c.logger.Info().Msg("Resolved existing block subscriber")
		return blocksFanout, nil
	}

	// we need this lock to prevent 2 Subscribe calls at the same time
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info().Msg("Subscribing to block events")

	// Subscribe to comet bft events
	eventsChan, err := c.cometBFTClient.Subscribe(context.Background(), "", ctypes.EventQueryNewBlock.String())
	if err != nil {
		return nil, errors.Wrap(err, "unable to subscribe to new block events")
	}

	c.logger.Info().Msg("Subscribed to block events")

	// Create block chan
	blockChan := make(chan ctypes.EventDataNewBlock)

	// Spin up a pipeline to forward block events to the blockChan
	go func() {
		for event := range eventsChan {
			newBlockEvent, ok := event.Data.(ctypes.EventDataNewBlock)
			if !ok {
				c.logger.Error().Msgf("expecting new block event, got %T", event.Data)
				continue
			}

			c.logger.Info().Int64("height", newBlockEvent.Block.Height).Msg("Received new block event")

			blockChan <- newBlockEvent
		}
	}()

	// Create a fanout
	// It allows a "global" chan (i.e. blockChan) to stream to multiple consumers independently.
	fo := fanout.New[ctypes.EventDataNewBlock](blockChan, fanout.DefaultBuffer)
	fo.Start()

	c.blocksFanout = fo

	return fo, nil
}

func (c *Client) blockFanOutThreadSafe() *fanout.FanOut[ctypes.EventDataNewBlock] {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.blocksFanout
}
