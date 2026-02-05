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
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// constants for monitoring tx results
const (
	monitorInterval   = constant.ZetaBlockTime
	monitorRetryCount = 10

	// defaultInboundVoteMonitorTimeout is the default timeout for monitoring inbound vote tx result.
	// In our case, the upstream code ALWAYS sets a timeout in the context to override default value,
	// so this is just to keep the logic complete and avoid accidental missed timeout in the context.
	defaultInboundVoteMonitorTimeout = 2 * time.Minute

	// monitorDeadlineOffset is subtracted from the context deadline to ensure the tx result query loop
	// finishes earlier, allowing the caller to write errors to the error monitor channel before ctx.Done() closes.
	monitorDeadlineOffset = 10 * time.Second
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
	logger := c.logger.With().Str(logs.FieldZetaTx, zetaTxHash).Logger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error().Any("panic", r).Msg("recovered from panic")
		}
	}()

	call := func() error {
		return retry.Retry(c.monitorVoteOutboundResult(ctx, zetaTxHash, retryGasLimit, msg))
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval/2, monitorInterval)
	if err != nil {
		logger.Error().Err(err).Msg("unable to query tx result")
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

	logger := c.logger.With().
		Str(logs.FieldZetaTx, zetaTxHash).
		Str("outbound_raw_log", txResult.RawLog).
		Logger()

	raw := strings.ToLower(txResult.RawLog)
	switch {
	case strings.Contains(raw, "already voted"):
		// noop
	case strings.Contains(raw, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute
		// this shouldn't happen
		logger.Error().Msg("failed to execute vote")
	case strings.Contains(raw, "out of gas"):
		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		logger.Debug().Msg("out of gas")

		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			if _, _, err := c.PostVoteOutbound(ctx, retryGasLimit, 0, msg); err != nil {
				logger.Error().Err(err).Msg("failed to resend tx")
			} else {
				logger.Info().Msg("successfully resent tx")
			}
		}
	default:
		logger.Debug().Msg("successful")
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
	logger := c.logger.With().Str(logs.FieldZetaTx, zetaTxHash).Logger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error().Any("panic", r).Msg("recovered from panic")
		}
	}()

	// Use 1 iteration for chaos error checks to fail instantly
	err := c.monitorVoteInboundResult(ctx, zetaTxHash, retryGasLimit, msg, monitorErrCh)
	if err != nil {
		if strings.Contains(err.Error(), "chaos error") {
			logger.Warn().
				Err(err).
				Str(logs.FieldZetaTx, zetaTxHash).
				Str(logs.FieldBallotIndex, msg.Digest()).
				Msg("chaos error simulating mempool congestion, skipping retries for inbound vote to trigger add internal tracker")
			return err
		}

		call := func() error {
			return c.monitorVoteInboundResult(ctx, zetaTxHash, retryGasLimit, msg, monitorErrCh)
		}

		// extract the deadline (always provided) that is used by error monitor goroutine
		deadline := time.Now().Add(defaultInboundVoteMonitorTimeout)
		if ctxDeadline, ok := ctx.Deadline(); ok {
			// subtract deadline offset from the context deadline, so the following tx result query loop
			// can finish earlier, otherwise the caller won't be able to write the error to the error
			// monitor channel monitorErrCh because 'ctx.Done()' is already closed.
			deadline = ctxDeadline.Add(-monitorDeadlineOffset)
		}

		// query tx result with the same deadline used by the upstream error monitor goroutine
		start := time.Now()
		bo := backoff.NewConstantBackOff(monitorInterval)
		err = retry.DoWithDeadline(call, bo, deadline)
		if err != nil {
			// we only return an error if the tx result cannot be queried
			logger.Error().
				Err(err).
				Str(logs.FieldZetaTx, zetaTxHash).
				Str(logs.FieldBallotIndex, msg.Digest()).
				Stringer("timeout", deadline.Sub(start)).
				Msg("tx result query timed out")
			return err
		}
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
	txResult, err := c.self.QueryTxResult(zetaTxHash)
	if err != nil {
		return errors.Wrap(err, "failed to query tx result")
	}

	logger := c.logger.With().Str("inbound_raw_log", txResult.RawLog).Logger()

	// There is no error returned from here which mean the MonitorVoteInboundResult would return nil and no error is posted to monitorErrCh
	// However the channel is passed to the subsequent call, which can post an error to the channel if the "execute" vote fails.
	switch {
	case strings.Contains(txResult.RawLog, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute. this shouldn't happen
		logger.Error().Str(logs.FieldZetaTx, zetaTxHash).Msg("failed to execute inbound vote")

	case strings.Contains(txResult.RawLog, "out of gas"):
		metrics.InboundVotesWithOutOfGasErrorsTotal.WithLabelValues(c.chain.Name).Inc()
		// record this ready-to-execute ballot for future gas adjustment
		// The 500K is enough for regular inbound vote, out of gas error happens only on the finalizing vote
		c.addReadyToExecuteInboundBallot(msg.Digest(), txResult.GasWanted, zetaTxHash)

		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		logger.Debug().Str(logs.FieldZetaTx, zetaTxHash).Msg("out of gas")
		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			resentZetaTxHash, _, err := c.PostVoteInbound(ctx, retryGasLimit, 0, msg, monitorErrCh)
			if err != nil {
				logger.Error().Err(err).Str(logs.FieldZetaTx, zetaTxHash).Msg("failed to resend tx")
			} else {
				logger.Info().Str(logs.FieldZetaTx, resentZetaTxHash).Msg("successfully resent tx")
			}
		}
	default:
		// it is just a nice-to-have logic to reduce the memory usage and
		// we don't expect this cleanup is perfect and cover 100% of the cases
		c.removeReadyToExecuteInboundBallot(msg.Digest())
		logger.Debug().Str(logs.FieldZetaTx, zetaTxHash).Msg("successful")
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
