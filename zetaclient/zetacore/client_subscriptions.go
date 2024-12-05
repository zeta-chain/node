package zetacore

import (
	"context"

	cometbft_types "github.com/cometbft/cometbft/types"
)

// NewBlockSubscriber subscribes to cometbft new block events
func (c *Client) NewBlockSubscriber(ctx context.Context) (chan cometbft_types.EventDataNewBlock, error) {
	rawBlockEventChan, err := c.cometBFTClient.Subscribe(ctx, "", cometbft_types.EventQueryNewBlock.String())
	if err != nil {
		return nil, err
	}

	blockEventChan := make(chan cometbft_types.EventDataNewBlock)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-rawBlockEventChan:
				newBlockEvent, ok := event.Data.(cometbft_types.EventDataNewBlock)
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
