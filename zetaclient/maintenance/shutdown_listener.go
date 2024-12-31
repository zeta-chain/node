package maintenance

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/retry"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

const restartListenerTicker = 10 * time.Second

// ShutdownListener is a struct that listens for scheduled shutdown notices via the observer
// operational flags
type ShutdownListener struct {
	client interfaces.ZetacoreClient
	logger zerolog.Logger

	lastRestartHeightMissed int64
}

// NewShutdownListener creates a new ShutdownListener.
func NewShutdownListener(client interfaces.ZetacoreClient, logger zerolog.Logger) *ShutdownListener {
	log := logger.With().Str("module", "shutdown_listener").Logger()
	return &ShutdownListener{
		client: client,
		logger: log,
	}
}

func (o *ShutdownListener) Listen(ctx context.Context, action func()) {
	var (
		withLogger = bg.WithLogger(o.logger)
		onComplete = bg.OnComplete(action)
	)

	bg.Work(ctx, o.waitForUpdate, bg.WithName("shutdown_listener.wait_for_update"), withLogger, onComplete)
}

func (o *ShutdownListener) waitForUpdate(ctx context.Context) error {
	operationalFlags, err := retry.DoTypedWithBackoffAndRetry(
		func() (observertypes.OperationalFlags, error) { return o.client.GetOperationalFlags(ctx) },
		retry.DefaultConstantBackoff(),
	)
	if err != nil {
		return errors.Wrap(err, "unable to get initial operational flags")
	}
	if o.handleNewFlags(ctx, operationalFlags) {
		return nil
	}

	ticker := time.NewTicker(restartListenerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			operationalFlags, err = o.client.GetOperationalFlags(ctx)
			if err != nil {
				return errors.Wrap(err, "unable to get operational flags")
			}
			if o.handleNewFlags(ctx, operationalFlags) {
				return nil
			}
		case <-ctx.Done():
			o.logger.Info().Msg("waitForUpdate (shutdown listener) stopped")
			return nil
		}
	}
}

// handleNewFlags processes the flags and returns true if a shutdown should be signaled
func (o *ShutdownListener) handleNewFlags(ctx context.Context, f observertypes.OperationalFlags) bool {
	if f.RestartHeight < 1 {
		return false
	}

	currentHeight, err := o.client.GetBlockHeight(ctx)
	if err != nil {
		o.logger.Error().Err(err).Msg("unable to get block height")
		return false
	}

	if f.RestartHeight < currentHeight {
		// only log restart height misseed once
		if o.lastRestartHeightMissed != f.RestartHeight {
			o.logger.Error().
				Int64("restart_height", f.RestartHeight).
				Int64("current_height", currentHeight).
				Msg("restart height missed")
		}
		o.lastRestartHeightMissed = f.RestartHeight
		return false
	}

	o.logger.Warn().
		Int64("restart_height", f.RestartHeight).
		Int64("current_height", currentHeight).
		Msg("restart scheduled")

	newBlockChan, err := o.client.NewBlockSubscriber(ctx)
	if err != nil {
		o.logger.Error().Err(err).Msg("unable to subscribe to new blocks")
		return false
	}
	for {
		select {
		case newBlock := <-newBlockChan:
			if newBlock.Block.Height >= f.RestartHeight {
				o.logger.Warn().
					Int64("restart_height", f.RestartHeight).
					Int64("current_height", newBlock.Block.Height).
					Msg("restart height reached")
				return true
			}
		case <-ctx.Done():
			return false
		}
	}
}
