package zetacore

import (
	"context"

	cometbfttypes "github.com/cometbft/cometbft/types"
)

// NewBlockSubscriber subscribes to cometbft new block events
func (c *Client) NewBlockSubscriber(ctx context.Context) (chan cometbfttypes.EventDataNewBlock, error) {
	rawBlockEventChan, err := c.cometBFTClient.Subscribe(ctx, "", cometbfttypes.EventQueryNewBlock.String())
	if err != nil {
		return nil, err
	}

	blockEventChan := make(chan cometbfttypes.EventDataNewBlock)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-rawBlockEventChan:
				newBlockEvent, ok := event.Data.(cometbfttypes.EventDataNewBlock)
				if !ok {
					c.logger.Error().Msgf("expecting new block event, got %T", event.Data)
					continue
				}
				blockEventChan <- newBlockEvent
			}
		}
	}()

	return blockEventChan, nil
}
