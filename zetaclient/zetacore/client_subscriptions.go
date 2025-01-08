package zetacore

import (
	"context"

	ctypes "github.com/cometbft/cometbft/types"

	"github.com/zeta-chain/node/pkg/fanout"
)

// NewBlockSubscriber subscribes to comet bft new block events.
// Subscribes share the same websocket connection but their channels are independent (fanout)
func (c *Client) NewBlockSubscriber(ctx context.Context) (chan ctypes.EventDataNewBlock, error) {
	blockSubscriber, err := c.resolveBlockSubscriber()
	if err != nil {
		return nil, err
	}

	// we need a "proxy" chan instead of directly returning blockSubscriber.Add()
	// to support context cancellation
	blocksChan := make(chan ctypes.EventDataNewBlock)

	go func() {
		consumer := blockSubscriber.Add()

		for {
			select {
			case <-ctx.Done():
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
	if blocksFanout, ok := c.getBlockFanoutChan(); ok {
		c.logger.Info().Msg("Resolved existing block subscriber")
		return blocksFanout, nil
	}

	// Subscribe to comet bft events
	eventsChan, err := c.cometBFTClient.Subscribe(context.Background(), "", ctypes.EventQueryNewBlock.String())
	if err != nil {
		return nil, err
	}

	c.logger.Info().Msg("Subscribed to new block events")

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

			blockChan <- newBlockEvent
		}
	}()

	// Create a fanout
	// It allows a "global" chan (i.e. blockChan) to stream to multiple consumers independently.
	c.mu.Lock()
	defer c.mu.Unlock()
	c.blocksFanout = fanout.New[ctypes.EventDataNewBlock](blockChan, fanout.DefaultBuffer)

	c.blocksFanout.Start()

	return c.blocksFanout, nil
}

func (c *Client) getBlockFanoutChan() (*fanout.FanOut[ctypes.EventDataNewBlock], bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.blocksFanout, c.blocksFanout != nil
}
