package zetacore

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	appcontext "github.com/zeta-chain/zetacore/zetaclient/context"
)

var logSampler = &zerolog.BasicSampler{N: 10}

// UpdateZetacoreContextWorker is a polling goroutine that checks and updates zetacore context at every height.
// todo implement graceful shutdown and work group
func (c *Client) UpdateZetacoreContextWorker(ctx context.Context, app *appcontext.AppContext) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().Interface("panic", r).Msg("UpdateZetacoreContextWorker: recovered from panic")
		}
	}()

	var (
		updateEvery = time.Duration(app.Config().ConfigUpdateTicker) * time.Second
		ticker      = time.NewTicker(updateEvery)
		logger      = c.logger.Sample(logSampler)
	)

	c.logger.Info().Msg("UpdateZetacoreContextWorker started")

	for {
		select {
		case <-ticker.C:
			c.logger.Debug().Msg("UpdateZetacoreContextWorker invocation")
			if err := c.UpdateZetacoreContext(ctx, app, false, logger); err != nil {
				c.logger.Err(err).Msg("UpdateZetacoreContextWorker failed to update config")
			}
		case <-c.stop:
			c.logger.Info().Msg("UpdateZetacoreContextWorker stopped")
			return
		}
	}
}
