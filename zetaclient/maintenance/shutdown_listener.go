package maintenance

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"
	"golang.org/x/mod/semver"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/retry"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// restartListenerTicker is the duration between checks used by the shutdown listener
// this is currently used in both waitForUpdate and waitUntilSyncing
const restartListenerTicker = 10 * time.Second

// waitForSyncing is the duration to allow the zetacorednode to sync up before signaling the shutdown of zetaclient
const waitForSyncing = 10 * time.Minute

// ShutdownListener is a struct that listens for scheduled shutdown notices via the observer
// operational flags
type ShutdownListener struct {
	client ZetacoreClient
	logger zerolog.Logger

	lastRestartHeightMissed int64
	// get the current version of zetaclient
	getVersion func() string

	restartListenerTicker time.Duration
	waitForSyncing        time.Duration
}

// NewShutdownListener creates a new ShutdownListener.
func NewShutdownListener(client ZetacoreClient, logger zerolog.Logger) *ShutdownListener {
	log := logger.With().Str("module", "shutdown_listener").Logger()
	return &ShutdownListener{
		client:                client,
		logger:                log,
		getVersion:            getVersionDefault,
		restartListenerTicker: restartListenerTicker,
		waitForSyncing:        waitForSyncing,
	}
}

// RunPreStartCheck runs any checks that must run before fully starting zetaclient.
// Specifically this should be run before any TSS P2P is started.
func (o *ShutdownListener) RunPreStartCheck(ctx context.Context) error {
	operationalFlags, err := o.getOperationalFlagsWithRetry(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get initial operational flags")
	}
	return o.checkMinimumVersion(operationalFlags)
}

func (o *ShutdownListener) Listen(ctx context.Context, action func()) {
	var (
		withLogger = bg.WithLogger(o.logger)
		onComplete = bg.OnComplete(action)
	)

	bg.Work(ctx, o.waitForUpdate, bg.WithName("shutdown_listener.wait_for_update"), withLogger, onComplete)
	bg.Work(ctx, o.waitUntilSyncing, bg.WithName("shutdown_listener.wait_until_syncing"), withLogger, onComplete)
}

// waitUntilSyncing checks if the node is syncing
// if it is syncing it returns nil which completes the bg task and calls onComplete
func (o *ShutdownListener) waitUntilSyncing(ctx context.Context) error {
	ticker := time.NewTicker(o.restartListenerTicker)

	defer ticker.Stop()

	var syncDetectedAt time.Time
	syncDetected := false

	for {
		select {
		case <-ticker.C:
			isSyncing, err := o.client.GetSyncStatus(ctx)
			if err != nil {
				return errors.Wrap(err, "unable to get sync status")
			}

			if isSyncing {
				if !syncDetected {
					syncDetectedAt = time.Now()
					syncDetected = true
					o.logger.Info().Msgf("Node syncing detected, waiting %s before shutdown", o.waitForSyncing.String())
				} else {
					if time.Since(syncDetectedAt) >= o.waitForSyncing {
						o.logger.Info().Msgf("Node still syncing after %s proceeding with shutdown", o.waitForSyncing.String())
						return nil
					}
				}
			} else {
				if syncDetected {
					syncDetected = false
				}
			}
		case <-ctx.Done():
			o.logger.Info().Msg("waitUntilSyncing (shutdown listener) stopped")
			return nil
		}
	}
}

func (o *ShutdownListener) waitForUpdate(ctx context.Context) error {
	operationalFlags, err := o.getOperationalFlagsWithRetry(ctx)
	if err != nil {
		return errors.Wrap(err, "get initial operational flags")
	}
	if o.handleNewFlags(ctx, operationalFlags) {
		return nil
	}

	ticker := time.NewTicker(o.restartListenerTicker)
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

func (o *ShutdownListener) getOperationalFlagsWithRetry(ctx context.Context) (observertypes.OperationalFlags, error) {
	return retry.DoTypedWithBackoffAndRetry(
		func() (observertypes.OperationalFlags, error) { return o.client.GetOperationalFlags(ctx) },
		retry.DefaultConstantBackoff(),
	)
}

// handleNewFlags processes the flags and returns true if a shutdown should be signaled
func (o *ShutdownListener) handleNewFlags(ctx context.Context, f observertypes.OperationalFlags) bool {
	if err := o.checkMinimumVersion(f); err != nil {
		o.logger.Error().Err(err).Any("operational_flags", f).Msg("minimum version check")
		return true
	}
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

func (o *ShutdownListener) checkMinimumVersion(f observertypes.OperationalFlags) error {
	if f.MinimumVersion != "" {
		currentVersion := o.getVersion()
		if semver.Compare(currentVersion, f.MinimumVersion) == -1 {
			return fmt.Errorf(
				"current version (%s) is less than minimum version (%s)",
				currentVersion,
				f.MinimumVersion,
			)
		}
	}
	return nil
}

func getVersionDefault() string {
	return constant.GetNormalizedVersion()
}
