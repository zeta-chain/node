package zetacore

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeta-chain/go-tss/blame"

	"github.com/zeta-chain/node/pkg/chains"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/crosschain/types"
	observerclient "github.com/zeta-chain/node/x/observer/client/cli"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

// PostVoteGasPrice posts a gas price vote. Returns txHash and error.
func (c *Client) PostVoteGasPrice(
	ctx context.Context,
	chain chains.Chain,
	gasPrice uint64, priorityFee, blockNum uint64,
) (string, error) {
	// get gas price multiplier for the chain
	multiplier := GasPriceMultiplier(chain)

	// #nosec G115 always in range
	gasPrice = uint64(float64(gasPrice) * multiplier)
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteGasPrice(signerAddress, chain.ChainId, gasPrice, priorityFee, blockNum)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	hash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, PostGasPriceGasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to broadcast vote gas price")
	}

	return hash, nil
}

// PostVoteTSS sends message to vote TSS. Returns txHash and error.
func (c *Client) PostVoteTSS(
	ctx context.Context,
	tssPubKey string,
	keyGenZetaHeight int64,
	status chains.ReceiveStatus,
) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := observertypes.NewMsgVoteTSS(signerAddress, tssPubKey, keyGenZetaHeight, status)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	zetaTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, PostTSSGasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to broadcast vote for setting tss")
	}

	return zetaTxHash, nil
}

// PostVoteBlameData posts blame data message to zetacore. Returns txHash and error.
func (c *Client) PostVoteBlameData(
	ctx context.Context,
	blame *blame.Blame,
	chainID int64,
	index string,
) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	zetaBlame := observertypes.Blame{
		Index:         index,
		FailureReason: blame.FailReason,
		Nodes:         observerclient.ConvertNodes(blame.BlameNodes),
	}
	msg := observertypes.NewMsgVoteBlameMsg(signerAddress, chainID, zetaBlame)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	var gasLimit uint64 = PostBlameDataGasLimit

	zetaTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, gasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to broadcast blame data")
	}

	return zetaTxHash, nil
}

// PostVoteOutbound posts a vote on an observed outbound tx from a MsgVoteOutbound.
// Returns tx hash, ballotIndex, and error.
func (c *Client) PostVoteOutbound(
	ctx context.Context,
	gasLimit, retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) (string, string, error) {
	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to wrap message with authz")
	}

	// don't post confirmation if it  already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := c.HasVoted(ctx, ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(
			err,
			"hasVoted check failed for ballot %s voter %s",
			ballotIndex,
			msg.Creator,
		)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}

	zetaTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, gasLimit, authzMsg, authzSigner)
	})
	if err != nil {
		return "", ballotIndex, errors.Wrap(err, "unable to broadcast vote outbound")
	}

	go func() {
		ctxForWorker := zctx.Copy(ctx, context.Background())
		err := c.MonitorVoteOutboundResult(ctxForWorker, zetaTxHash, retryGasLimit, msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to monitor vote outbound result")
		}
	}()

	return zetaTxHash, ballotIndex, nil
}

// PostVoteInbound posts a vote on an observed inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// it is used when the ballot is finalized and the inbound tx needs to be processed
func (c *Client) PostVoteInbound(
	ctx context.Context,
	gasLimit, retryGasLimit uint64,
	msg *types.MsgVoteInbound,
	monitorErrCh chan<- zetaerrors.ErrTxMonitor,
) (string, string, error) {
	// zetaclient patch
	// force use SAFE mode for all inbound votes (both fast and slow votes)
	msg.ConfirmationMode = types.ConfirmationMode_SAFE

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", err
	}

	// don't post send if has already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := c.HasVoted(ctx, ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(err,
			"PostVoteInbound: unable to check if already voted for ballot %s voter %s",
			ballotIndex,
			msg.Creator,
		)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}

	zetaTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, gasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", ballotIndex, errors.Wrap(err, "unable to broadcast vote inbound")
	}

	go func() {
		// Use the passed context directly instead of creating a new one
		// This ensures the monitoring goroutine respects the same timeout as the error handler
		errMonitor := c.MonitorVoteInboundResult(ctx, zetaTxHash, retryGasLimit, msg, monitorErrCh)
		if errMonitor != nil {
			c.logger.Error().Err(errMonitor).Msg("failed to monitor vote inbound result")

			if monitorErrCh != nil {
				select {
				case monitorErrCh <- zetaerrors.ErrTxMonitor{
					Err:                errMonitor,
					InboundBlockHeight: msg.InboundBlockHeight,
					ZetaTxHash:         zetaTxHash,
					BallotIndex:        ballotIndex,
				}:
				case <-ctx.Done():
					c.logger.Error().Msg("context cancelled: timeout")
				}
			}
		}
	}()

	return zetaTxHash, ballotIndex, nil
}
