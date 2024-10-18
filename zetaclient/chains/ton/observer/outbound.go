package observer

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/ticker"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

func (ob *Observer) watchOutbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	var (
		initialInterval = ticker.SecondsFromUint64(ob.ChainParams().OutboundTicker)
		sampledLogger   = ob.Logger().Inbound.Sample(&zerolog.BasicSampler{N: 10})
	)

	task := func(ctx context.Context, t *ticker.Ticker) error {
		if !app.IsOutboundObservationEnabled() {
			sampledLogger.Info().Msg("WatchOutbound: outbound observation is disabled")
			return nil
		}

		if err := ob.observeOutbound(ctx); err != nil {
			ob.Logger().Outbound.Err(err).Msg("WatchOutbound: observeOutbound error")
		}

		newInterval := ticker.SecondsFromUint64(ob.ChainParams().OutboundTicker)
		t.SetInterval(newInterval)

		return nil
	}

	return ticker.Run(
		ctx,
		initialInterval,
		task,
		ticker.WithStopChan(ob.StopChannel()),
		ticker.WithLogger(ob.Logger().Outbound, "WatchOutbound"),
	)
}

func (ob *Observer) observeOutbound(_ context.Context) error {
	// todo implement
	return nil
}
