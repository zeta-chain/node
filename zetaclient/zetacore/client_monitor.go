package zetacore

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/pkg/retry"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// MonitorVoteOutboundResult monitors the result of a vote outbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (c *Client) MonitorVoteOutboundResult(
	ctx context.Context,
	zetaTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().
				Interface("panic", r).
				Str("outbound.hash", zetaTxHash).
				Msg("monitorVoteOutboundResult: recovered from panic")
		}
	}()

	call := func() error {
		return retry.Retry(c.monitorVoteOutboundResult(ctx, zetaTxHash, retryGasLimit, msg))
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval/2, monitorInterval)
	if err != nil {
		c.logger.Error().Err(err).
			Str("outbound.hash", zetaTxHash).
			Msg("monitorVoteOutboundResult: unable to query tx result")

		return err
	}

	return nil
}

func (c *Client) monitorVoteOutboundResult(
	ctx context.Context,
	zetaTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) error {
	// query tx result from ZetaChain
	txResult, err := c.QueryTxResult(zetaTxHash)
	if err != nil {
		return errors.Wrap(err, "failed to query tx result")
	}

	logFields := map[string]any{
		"outbound.hash":    zetaTxHash,
		"outbound.raw_log": txResult.RawLog,
	}

	switch {
	case strings.Contains(txResult.RawLog, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute
		// this shouldn't happen
		c.logger.Error().Fields(logFields).Msg("monitorVoteOutboundResult: failed to execute vote")
	case strings.Contains(txResult.RawLog, "out of gas"):
		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		c.logger.Debug().Fields(logFields).Msg("monitorVoteOutboundResult: out of gas")

		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			if _, _, err := c.PostVoteOutbound(ctx, retryGasLimit, 0, msg); err != nil {
				c.logger.Error().Err(err).Fields(logFields).Msg("monitorVoteOutboundResult: failed to resend tx")
			} else {
				c.logger.Info().Fields(logFields).Msg("monitorVoteOutboundResult: successfully resent tx")
			}
		}
	default:
		c.logger.Debug().Fields(logFields).Msg("monitorVoteOutboundResult: successful")
	}

	return nil
}

// MonitorVoteInboundResult monitors the result of a vote inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (c *Client) MonitorVoteInboundResult(
	ctx context.Context,
	zetaTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteInbound,
) error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().
				Interface("panic", r).
				Str("inbound.hash", zetaTxHash).
				Msg("monitorVoteInboundResult: recovered from panic")
		}
	}()

	call := func() error {
		return retry.Retry(c.monitorVoteInboundResult(ctx, zetaTxHash, retryGasLimit, msg))
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval/2, monitorInterval)
	if err != nil {
		c.logger.Error().Err(err).
			Str("inbound.hash", zetaTxHash).
			Msg("monitorVoteInboundResult: unable to query tx result")

		return err
	}

	return nil
}

func (c *Client) monitorVoteInboundResult(
	ctx context.Context,
	zetaTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteInbound,
) error {
	// query tx result from ZetaChain
	txResult, err := c.QueryTxResult(zetaTxHash)
	if err != nil {
		return errors.Wrap(err, "failed to query tx result")
	}

	logFields := map[string]any{
		"inbound.hash":    zetaTxHash,
		"inbound.raw_log": txResult.RawLog,
	}

	switch {
	case strings.Contains(txResult.RawLog, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute
		// this shouldn't happen
		c.logger.Error().Fields(logFields).Msg("monitorVoteInboundResult: failed to execute vote")

	case strings.Contains(txResult.RawLog, "out of gas"):
		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		c.logger.Debug().Fields(logFields).Msg("monitorVoteInboundResult: out of gas")

		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			if _, _, err := c.PostVoteInbound(ctx, retryGasLimit, 0, msg); err != nil {
				c.logger.Error().Err(err).Fields(logFields).Msg("monitorVoteInboundResult: failed to resend tx")
			} else {
				c.logger.Info().Fields(logFields).Msg("monitorVoteInboundResult: successfully resent tx")
			}
		}

	default:
		c.logger.Debug().Fields(logFields).Msgf("monitorVoteInboundResult: successful")
	}

	return nil
}

func retryWithBackoff(call func() error, attempts int, minInternal, maxInterval time.Duration) error {
	if attempts < 1 {
		return errors.New("attempts must be positive")
	}

	bo := backoff.WithMaxRetries(
		backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(minInternal),
			backoff.WithMaxInterval(maxInterval),
		),
		// #nosec G115 always positive
		uint64(attempts),
	)

	return retry.DoWithBackoff(call, bo)
}
