// Package maintenance provides maintenance functionalities for the zetaclient.
package maintenance

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/bg"
	"github.com/zeta-chain/zetacore/pkg/retry"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

const tssListenerTicker = 5 * time.Second

// TSSListener is a struct that listens for TSS updates, new keygen, and new TSS key generation.
type TSSListener struct {
	client interfaces.ZetacoreClient
	logger zerolog.Logger
}

// NewTSSListener creates a new TSSListener.
func NewTSSListener(client interfaces.ZetacoreClient, logger zerolog.Logger) *TSSListener {
	log := logger.With().Str("module", "tss_listener").Logger()

	return &TSSListener{
		client: client,
		logger: log,
	}
}

// Listen listens for any maintenance regarding TSS and calls action specified. Works in the background.
func (tl *TSSListener) Listen(ctx context.Context, action func()) {
	var (
		withLogger = bg.WithLogger(tl.logger)
		onComplete = bg.OnComplete(action)
	)

	bg.Work(ctx, tl.waitForUpdate, bg.WithName("tss.wait_for_update"), withLogger, onComplete)
	bg.Work(ctx, tl.waitForNewKeyGeneration, bg.WithName("tss.wait_for_generation"), withLogger, onComplete)
	bg.Work(ctx, tl.waitForNewKeygen, bg.WithName("tss.wait_for_keygen"), withLogger, onComplete)
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
				tl.logger.Warn().Err(err).Msg("unable to get new tss")
				continue
			}

			if tssNew.TssPubkey == tss.TssPubkey {
				continue
			}

			tl.logger.Info().
				Str("tss.old", tss.TssPubkey).
				Str("tss.new", tssNew.TssPubkey).
				Msgf("TSS address is updated")

			return nil
		case <-ctx.Done():
			tl.logger.Info().Msg("waitForTSSUpdate stopped")
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
		return errors.Wrap(err, "failed to get initial tss history")
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

			if tssLenUpdated <= tssLen {
				continue
			}

			tl.logger.Info().Msgf("tss list updated from %d to %d", tssLen, tssLenUpdated)
			return nil
		case <-ctx.Done():
			tl.logger.Info().Msg("waitForNewKeyGeneration stopped")
			return nil
		}
	}
}

// waitForNewKeygen is a background thread that listens for new keygen; it returns when a new keygen is set
func (tl *TSSListener) waitForNewKeygen(ctx context.Context) error {
	// Initial Keygen retrieval
	keygen, err := retry.DoTypedWithBackoffAndRetry(
		func() (*observertypes.Keygen, error) { return tl.client.GetKeyGen(ctx) },
		retry.DefaultConstantBackoff(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get initial tss history")
	}

	ticker := time.NewTicker(tssListenerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			keygenUpdated, err := tl.client.GetKeyGen(ctx)
			switch {
			case err != nil:
				tl.logger.Warn().Err(err).Msg("unable to get keygen")
				continue
			case keygenUpdated == nil:
				continue
			case keygenUpdated.Status == observertypes.KeygenStatus_PendingKeygen:
				continue
			case keygen.BlockNumber == keygenUpdated.BlockNumber:
				continue
			}

			tl.logger.Info().Msgf("got new keygen at block %d", keygen.BlockNumber)
			return nil
		case <-ctx.Done():
			tl.logger.Info().Msg("waitForNewKeygen stopped")
			return nil
		}
	}
}
