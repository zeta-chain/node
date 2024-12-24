package orchestrator

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// V2 represents the orchestrator V2 while they co-exist with Orchestrator.
type V2 struct {
	zetacore interfaces.ZetacoreClient

	scheduler *scheduler.Scheduler

	chains map[int64]ObserverSigner

	logger zerolog.Logger
}

const schedulerGroup = scheduler.Group("orchestrator")

type ObserverSigner interface {
	Start(ctx context.Context) error
	Stop()
}

func NewV2(
	zetacore interfaces.ZetacoreClient,
	scheduler *scheduler.Scheduler,
	logger zerolog.Logger,
) *V2 {
	return &V2{
		zetacore:  zetacore,
		scheduler: scheduler,
		chains:    make(map[int64]ObserverSigner),
		logger:    logger.With().Str(logs.FieldModule, "orchestrator").Logger(),
	}
}

func (oc *V2) Start(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	contextUpdaterOpts := []scheduler.Opt{
		scheduler.GroupName(schedulerGroup),
		scheduler.Name("update_context"),
		scheduler.IntervalUpdater(func() time.Duration {
			return ticker.DurationFromUint64Seconds(app.Config().ConfigUpdateTicker)
		}),
	}

	oc.scheduler.Register(ctx, oc.UpdateContext, contextUpdaterOpts...)

	return nil
}

func (oc *V2) Stop() {
	oc.logger.Info().Msg("Stopping orchestrator")

	// stops *all* scheduler tasks
	oc.scheduler.Stop()
}

func (oc *V2) UpdateContext(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	err = UpdateAppContext(ctx, app, oc.zetacore, oc.logger)

	switch {
	case errors.Is(err, ErrUpgradeRequired):
		const msg = "Upgrade detected. Kill the process, " +
			"replace the binary with upgraded version, and restart zetaclientd"

		oc.logger.Warn().Str("upgrade", err.Error()).Msg(msg)

		// stop the orchestrator
		go oc.Stop()

		return nil
	case err != nil:
		return errors.Wrap(err, "unable to update app context")
	default:
		return nil
	}
}
