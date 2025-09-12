package zetacore

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/constant"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// constants for monitoring tx results
const (
	monitorInterval   = constant.ZetaBlockTime / 2
	monitorRetryCount = 10
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
				Any("panic", r).
				Str(logs.FieldZetaTx, zetaTxHash).
				Str(logs.FieldMethod, "monitorVoteOutboundResult").
				Msg("Recovered from panic")
		}
	}()

	call := func() error {
		return retry.Retry(c.monitorVoteOutboundResult(ctx, zetaTxHash, retryGasLimit, msg))
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval/2, monitorInterval)
	if err != nil {
		c.logger.Error().Err(err).
			Str(logs.FieldZetaTx, zetaTxHash).
			Str(logs.FieldMethod, "monitorVoteOutboundResult").
			Msg("Unable to query tx result")

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
		logs.FieldZetaTx:   zetaTxHash,
		"outbound.raw_log": txResult.RawLog,
	}

	raw := strings.ToLower(txResult.RawLog)

	switch {
	case strings.Contains(raw, "already voted"):
		// noop
	case strings.Contains(raw, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute
		// this shouldn't happen
		c.logger.Error().Fields(logFields).Msg("monitorVoteOutboundResult: failed to execute vote")
	case strings.Contains(raw, "out of gas"):
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
	monitorErrCh chan<- zetaerrors.ErrTxMonitor,
) error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().
				Any("panic", r).
				Str(logs.FieldZetaTx, zetaTxHash).
				Msg("monitorVoteInboundResult: recovered from panic")
		}
	}()

	call := func() error {
		err := c.monitorVoteInboundResult(ctx, zetaTxHash, retryGasLimit, msg, monitorErrCh)

		// force retry on err
		return retry.Retry(err)
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval, monitorInterval*2)
	if err != nil {
		// All errors are forced to be retryable, we only return an error if the tx result cannot be queri
		c.logger.Error().Err(err).
			Str(logs.FieldZetaTx, zetaTxHash).
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
	monitorErrCh chan<- zetaerrors.ErrTxMonitor,
) error {
	// query tx result from ZetaChain
	txResult, err := c.QueryTxResult(zetaTxHash)
	if err != nil {
		return errors.Wrap(err, "failed to query tx result")
	}

	logFields := map[string]any{
		logs.FieldZetaTx:  zetaTxHash,
		"inbound.raw_log": txResult.RawLog,
	}

	// There is no error returned from here which mean the MonitorVoteInboundResult would return nil and no error is posted to monitorErrCh
	// However the channel is passed to the subsequent call, which can post an error to the channel if the "execute" vote fails.
	switch {
	case strings.Contains(txResult.RawLog, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute. this shouldn't happen
		c.logger.Error().Fields(logFields).Msg("monitorVoteInboundResult: failed to execute vote")

	case strings.Contains(txResult.RawLog, "out of gas"):
		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		c.logger.Debug().Fields(logFields).Msg("monitorVoteInboundResult: out of gas")

		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			resentZetaTxHash, _, err := c.PostVoteInbound(ctx, retryGasLimit, 0, msg, monitorErrCh)
			if err != nil {
				c.logger.Error().Err(err).Fields(logFields).Msg("monitorVoteInboundResult: failed to resend tx")
			} else {
				logFields[logs.FieldZetaTx] = resentZetaTxHash
				c.logger.Info().Fields(logFields).Msgf("monitorVoteInboundResult: successfully resent tx")
			}
		}
	default:
		c.logger.Debug().Fields(logFields).Msg("monitorVoteInboundResult: successful")
	}

	return nil
}

func retryWithBackoff(call func() error, attempts uint64, minInternal, maxInterval time.Duration) error {
	if attempts == 0 {
		return errors.New("attempts must be positive")
	}

	bo := backoff.WithMaxRetries(
		backoff.NewExponentialBackOff(
			backoff.WithInitialInterval(minInternal),
			backoff.WithMaxInterval(maxInterval),
		),
		attempts,
	)

	return retry.DoWithBackoff(call, bo)
}
