package zetacore

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	appcontext "github.com/zeta-chain/node/zetaclient/context"
)

var logSampler = &zerolog.BasicSampler{N: 10}

// UpdateAppContextWorker is a polling goroutine that checks and updates AppContext at every height.
// todo implement graceful shutdown and work group
func (c *Client) UpdateAppContextWorker(ctx context.Context, app *appcontext.AppContext) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().Interface("panic", r).Msg("UpdateAppContextWorker: recovered from panic")
		}
	}()

	var (
		updateEvery = time.Duration(app.Config().ConfigUpdateTicker) * time.Second
		ticker      = time.NewTicker(updateEvery)
		logger      = c.logger.Sample(logSampler)
	)

	c.logger.Info().Msg("UpdateAppContextWorker started")

	for {
		select {
		case <-ticker.C:
			c.logger.Debug().Msg("UpdateAppContextWorker invocation")
			if err := c.UpdateAppContext(ctx, app, logger); err != nil {
				c.logger.Err(err).Msg("UpdateAppContextWorker failed to update config")
			}
		case <-c.stop:
			c.logger.Info().Msg("UpdateAppContextWorker stopped")
			return
		}
	}
}
