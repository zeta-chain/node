// Package maintenance provides maintenance functionalities for the zetaclient.
package maintenance

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/retry"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const tssListenerTicker = 5 * time.Second

// TSSListener is a struct that listens for TSS updates, new keygen, and new TSS key generation.
type TSSListener struct {
	client ZetacoreClient
	logger zerolog.Logger
}

// NewTSSListener creates a new TSSListener.
func NewTSSListener(client ZetacoreClient, logger zerolog.Logger) *TSSListener {
	log := logger.With().Str("module", "tss_listener").Logger()

	return &TSSListener{
		client: client,
		logger: log,
	}
}

// Listen listens for any maintenance regarding TSS and calls action specified. Works in the background.
//
// Both watchers key off the TSS itself: the address changing, or a new key landing in history.
// The keygen record is deliberately not watched. It is reset to "pending at block MaxInt64" on
// any observer set change, and restarting on that reset is what turns a routine validator
// unbonding into an outage — every signer shuts down at once and none of them can start again.
//
// Note this also removes the only trigger that restarted zetaclient for a scheduled keygen, so
// while a finalized key exists a rotation ceremony will not start. That is deliberate; see the
// step 5 comment in zetaclient/tss/setup.go.
func (tl *TSSListener) Listen(ctx context.Context, action func()) {
	var (
		withLogger = bg.WithLogger(tl.logger)
		onComplete = bg.OnComplete(action)
	)

	bg.Work(ctx, tl.waitForUpdate, bg.WithName("tss.wait_for_update"), withLogger, onComplete)
	bg.Work(ctx, tl.waitForNewKeyGeneration, bg.WithName("tss.wait_for_generation"), withLogger, onComplete)
}

// waitForUpdate listens for TSS updates. Returns `nil` when the TSS address is updated
func (tl *TSSListener) waitForUpdate(ctx context.Context) error {
	// Initial TSS retrieval
	tss, err := retry.DoTypedWithBackoffAndRetry(
		func() (observertypes.TSS, error) { return tl.client.GetTSS(ctx) },
		retry.DefaultConstantBackoff(),
	)

	if err != nil {
		return errors.Wrap(err, "unable to get initial tss")
	}

	ticker := time.NewTicker(tssListenerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tssNew, err := tl.client.GetTSS(ctx)
			if err != nil {
				tl.logger.Warn().Err(err).Msg("unable to get new TSS")
				continue
			}
			// If the TSS address is not updated, continue loop
			if tssNew.TssPubkey == tss.TssPubkey {
				continue
			}

			tl.logger.Info().
				Str("tss_old", tss.TssPubkey).
				Str("tss_new", tssNew.TssPubkey).
				Msg("updated the TSS address")

			return nil
		case <-ctx.Done():
			tl.logger.Info().Msg("stopped waiting for updates in the TSS listener")
			return nil
		}
	}
}

// waitForNewKeyGeneration waits for new TSS key generation; it returns when a new key is generated
// It uses the length of the TSS list to determine if a new key is generated
func (tl *TSSListener) waitForNewKeyGeneration(ctx context.Context) error {
	// Initial TSS history retrieval
	tssHistoricalList, err := retry.DoTypedWithBackoffAndRetry(
		func() ([]observertypes.TSS, error) { return tl.client.GetTSSHistory(ctx) },
		retry.DefaultConstantBackoff(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get initial TSS history")
	}

	tssLen := len(tssHistoricalList)

	ticker := time.NewTicker(tssListenerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tssHistoricalListNew, err := tl.client.GetTSSHistory(ctx)
			if err != nil {
				continue
			}

			tssLenUpdated := len(tssHistoricalListNew)
			// New tss key has not been added to list , continue loop
			if tssLenUpdated <= tssLen {
				continue
			}

			tl.logger.Info().
				Int("from_length", tssLen).
				Int("to_length", tssLenUpdated).
				Msg("updated the TSS list")
			return nil
		case <-ctx.Done():
			tl.logger.Info().Msg("stopped waiting for new key generation in the TSS listener")
			return nil
		}
	}
}
