package zetacore

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/bg"
	"github.com/zeta-chain/zetacore/pkg/retry"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
)

// startBackgroundRoutines: This function will start background threads.
// These threads are responsible for handling TSS updates, new keygen, and new TSS key generation.
// These threads are provided with a cancel function which is used to restart the main thread based on the outcome of the background task.
func (c *Client) StartTssMigrationRoutines(
	ctx context.Context,
	cancelFunc context.CancelCauseFunc,
	masterLogger zerolog.Logger,
) context.CancelFunc {
	backgroundContext, cancel := context.WithCancel(ctx)
	bg.Work(
		backgroundContext,
		c.HandleTSSUpdate,
		bg.WithName("HandleTSSUpdate"),
		bg.WithLogger(masterLogger),
		bg.WithCancel(cancelFunc),
	)
	bg.Work(
		backgroundContext,
		c.HandleNewKeygen,
		bg.WithName("HandleNewKeygen"),
		bg.WithLogger(masterLogger),
		bg.WithCancel(cancelFunc),
	)
	bg.Work(
		backgroundContext,
		c.HandleNewTSSKeyGeneration,
		bg.WithName("HandleNewTSSKeyGeneration"),
		bg.WithLogger(masterLogger),
		bg.WithCancel(cancelFunc),
	)
	return cancel
}

// HandleTSSUpdate is a background thread that listens for TSS updates; it returns when the TSS address is updated
func (c *Client) HandleTSSUpdate(ctx context.Context) error {
	appCtx, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get app context")
	}

	logger := appCtx.Logger().With().Str("module", "HandleTSSUpdate").Logger()

	// Initial TSS retrieval
	tss, err := retry.DoTypedWithBackoffAndRetry[observertypes.TSS](func() (observertypes.TSS, error) {
		return c.GetTSS(ctx)
	}, retry.DefaultConstantBackoff())
	if err != nil {
		logger.Warn().Err(err).Msg("unable to get initial tss")
		return err
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			{
				tssNew, err := c.GetTSS(ctx)
				if err != nil {
					logger.Warn().Err(err).Msg("unable to get new tss")
					continue
				}

				if tssNew.TssPubkey == tss.TssPubkey {
					continue
				}
				logger.Info().Msgf("tss address is updated from %s to %s", tss.TssPubkey, tssNew.TssPubkey)
				return nil
			}
		case <-ctx.Done():
			{
				logger.Info().Msg("HandleTSSUpdate stopped")
				return nil
			}
		}
	}
}

// HandleNewTSSKeyGeneration is a background thread that listens for new TSS key generation; it returns when a new key is generated
// It uses the length of the TSS list to determine if a new key is generated
func (c *Client) HandleNewTSSKeyGeneration(ctx context.Context) error {
	appCtx, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get app context")
	}

	logger := appCtx.Logger().With().Str("module", "HandleNewTSSKeyGeneration").Logger()

	// Initial TSS history retrieval
	tssHistoricalList, err := retry.DoTypedWithBackoffAndRetry[[]observertypes.TSS](
		func() ([]observertypes.TSS, error) {
			return c.GetTSSHistory(ctx)
		},
		retry.DefaultConstantBackoff(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get initial tss history")
	}
	tssLen := len(tssHistoricalList)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			{
				tssHistoricalListNew, err := c.GetTSSHistory(ctx)
				if err != nil {
					continue
				}
				tssLenUpdated := len(tssHistoricalListNew)

				if tssLenUpdated <= tssLen {
					continue
				}
				logger.Info().Msgf("tss list updated from %d to %d", tssLen, tssLenUpdated)
				return nil
			}
		case <-ctx.Done():
			{
				logger.Info().Msg("HandleNewTSSKeyGeneration stopped")
				return nil
			}
		}
	}
}

// HandleNewKeygen is a background thread that listens for new keygen; it returns when a new keygen is set
func (c *Client) HandleNewKeygen(ctx context.Context) error {
	appCtx, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}
	logger := appCtx.Logger().With().Str("module", "HandleNewKeygen").Logger()

	// Initial Keygen retrieval
	keygen, err := retry.DoTypedWithBackoffAndRetry[*observertypes.Keygen](func() (*observertypes.Keygen, error) {
		return c.GetKeyGen(ctx)
	}, retry.DefaultConstantBackoff())
	if err != nil {
		return errors.Wrap(err, "failed to get initial tss history")
	}

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			{
				keygenUpdated, err := c.GetKeyGen(ctx)
				if err != nil {
					logger.Warn().Err(err).Msg("unable to get keygen")
					continue
				}
				if keygenUpdated == nil {
					logger.Warn().Err(err).Msg("keygen is nil")
					continue
				}
				if keygenUpdated.Status != observertypes.KeygenStatus_PendingKeygen {
					continue
				}

				if keygen.BlockNumber == keygenUpdated.BlockNumber {
					continue
				}

				logger.Info().Msgf("got new keygen at block %d", keygen.BlockNumber)
				return nil
			}
		case <-ctx.Done():
			{
				logger.Info().Msg("HandleNewKeygen stopped")
				return nil
			}
		}
	}
}
