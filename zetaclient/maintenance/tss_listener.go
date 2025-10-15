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
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
)

const tssListenerTicker = 5 * time.Second

// TSSListener is a struct that listens for TSS updates, new keygen, and new TSS key generation.
type TSSListener struct {
	client zrepo.ZetacoreClient
	logger zerolog.Logger
}

// NewTSSListener creates a new TSSListener.
func NewTSSListener(client zrepo.ZetacoreClient, logger zerolog.Logger) *TSSListener {
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

// waitForNewKeygen is a background thread that listens for new keygen; it returns when a new keygen is set
func (tl *TSSListener) waitForNewKeygen(ctx context.Context) error {
	// Initial Keygen retrieval
	keygen, err := tl.client.GetKeyGen(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get initial TSS history")
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
			// Keygen is not pending it has already been successfully generated, continue loop
			case keygenUpdated.Status == observertypes.KeygenStatus_KeyGenSuccess:
				continue
			// Keygen failed we to need to wait until a new keygen is set, continue loop
			case keygenUpdated.Status == observertypes.KeygenStatus_KeyGenFailed:
				continue
			// Keygen is pending but block number is not updated, continue loop.
			// Most likely the zetaclient is waiting for the keygen block to arrive.
			case keygenUpdated.Status == observertypes.KeygenStatus_PendingKeygen &&
				keygenUpdated.BlockNumber <= keygen.BlockNumber:
				continue
			}

			// Trigger restart only when the following conditions are met:
			// 1. Keygen is pending
			// 2. Block number is updated

			tl.logger.Info().Int64("block_number", keygenUpdated.BlockNumber).Msg("got new keygen")
			return nil
		case <-ctx.Done():
			tl.logger.Info().Msg("stopped waiting for new keygen in the TSS listener")
			return nil
		}
	}
}
