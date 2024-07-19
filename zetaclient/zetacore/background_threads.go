package zetacore

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/zeta-chain/zetacore/pkg/retry"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
)

func (c *Client) HandleTSSUpdate(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get app context")
	}

	logger := app.Logger().With().Str("module", "HandleTSSUpdate").Logger()

	bo := backoff.NewConstantBackOff(5 * time.Second)
	backoff.WithMaxRetries(bo, 10)

	// Initial TSS retrieval
	tss, err := retry.DoTypedWithBackoffAndRetry[observertypes.TSS](func() (observertypes.TSS, error) {
		return c.GetTSS(ctx)
	}, bo)
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
				tssNew, err := retry.DoTypedWithBackoffAndRetry[observertypes.TSS](func() (observertypes.TSS, error) {
					return c.GetTSS(ctx)
				}, bo)
				if err != nil {
					logger.Warn().Err(err).Msg("unable to get new tss")
					continue
				}

				if tssNew.TssPubkey == tss.TssPubkey {
					continue
				}
				tss = tssNew
				logger.Info().Msgf("tss address is updated from %s to %s", tss.TssPubkey, tssNew.TssPubkey)
				logger.Info().Msg("restarting zetaclient to update tss address")
				return nil
			}
		case <-ctx.Done():
			{
				return errors.Wrap(ctx.Err(), "context done")
			}
		}
	}
}

func (c *Client) HandleNewTSSKeyGeneration(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get app context")
	}

	logger := app.Logger().With().Str("module", "HandleNewTSSKeyGeneration").Logger()

	bo := backoff.NewConstantBackOff(5 * time.Second)
	backoff.WithMaxRetries(bo, 10)

	// Initial TSS retrieval
	tssHistoricalList, err := retry.DoTypedWithBackoffAndRetry[[]observertypes.TSS](func() ([]observertypes.TSS, error) {
		return c.GetTSSHistory(ctx)
	}, bo)
	if err != nil {
		return errors.Wrap(err, "failed to get initial tss history")
	}
	tssLen := len(tssHistoricalList)

	fmt.Println("tssLen old: ", tssLen)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			{
				tssHistoricalListNew, err := retry.DoTypedWithBackoffAndRetry[[]observertypes.TSS](func() ([]observertypes.TSS, error) {
					return c.GetTSSHistory(ctx)
				}, bo)
				if err != nil {
					continue
				}
				tssLenUpdated := len(tssHistoricalListNew)
				fmt.Println("tssLen updated: ", tssLenUpdated)

				if tssLenUpdated == tssLen {
					continue
				}
				if tssLenUpdated < tssLen {
					tssLen = tssLenUpdated
					continue
				}
				logger.Info().Msgf("tss list updated from %d to %d", tssLen, tssLenUpdated)
				tssLen = tssLenUpdated
				logger.Info().Msg("restarting zetaclient to update tss list")
				return nil
			}
		case <-ctx.Done():
			{
				return errors.Wrap(ctx.Err(), "context done")
			}
		}
	}
}
func (c *Client) HandleNewKeygen(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}
	logger := app.Logger().With().Str("module", "HandleNewKeygen").Logger()

	bo := backoff.NewConstantBackOff(5 * time.Second)
	backoff.WithMaxRetries(bo, 10)

	// Initial TSS retrieval
	keygen, err := retry.DoTypedWithBackoffAndRetry[*observertypes.Keygen](func() (*observertypes.Keygen, error) {
		return c.GetKeyGen(ctx)
	}, bo)
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

				keygen = keygenUpdated
				logger.Info().Msgf("got new keygen at block %d", keygen.BlockNumber)
				return nil
			}
		case <-ctx.Done():
			{
				return errors.Wrap(ctx.Err(), "context done")
			}
		}
	}
}
