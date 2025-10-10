package base

import (
	"context"
	"time"

	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// internalTrackerRetryInterval is the minimum interval between retries for each internal tracker
	// The vote tx may still be pending in mempool after each retry during mempool congestion, we should
	// give some time for it to be processed, without retrying too often and causing unnecessary mempool spam.
	internalTrackerRetryInterval = 5 * time.Minute
)

type InternalInboundTracker struct {
	// CreatedAt is the time when the tracker is created
	CreatedAt time.Time

	// LastRetry is the time when the tracker was last retried
	LastRetry time.Time

	// Tracker is the inbound tracker struct
	Tracker crosschaintypes.InboundTracker
}

// GetInboundInternalTrackers returns internal inbound trackers
func (ob *Observer) GetInboundInternalTrackers(
	ctx context.Context,
	retryTime time.Time,
) []crosschaintypes.InboundTracker {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		voterAddress     = ob.ZetaRepo().GetOperatorAddress()
		trackersToRemove = make([]string, 0)
		internalTrackers = make([]crosschaintypes.InboundTracker, 0, len(ob.internalInboundTrackers))
	)

	// collect up to MaxInternalTrackersPerScan trackers
	for ballot, tracker := range ob.internalInboundTrackers {
		// skip those that are already finalized
		if exist, err := ob.ZetaRepo().CCTXExists(ctx, ballot); err == nil && exist {
			ob.logger.Inbound.Info().Msgf("removing ballot %s from internal tracker as it is already finalized", ballot)
			trackersToRemove = append(trackersToRemove, ballot)
			continue
		}

		// skip those that have already voted
		if hasVoted, err := ob.ZetaRepo().HasVoted(ctx, ballot, voterAddress); err == nil && hasVoted {
			ob.logger.Inbound.Info().Msgf("removing ballot %s from internal tracker as it is already voted", ballot)
			trackersToRemove = append(trackersToRemove, ballot)
			continue
		}

		// skip those that have already been retried recently
		if retryTime.Sub(tracker.LastRetry) < internalTrackerRetryInterval {
			continue
		}

		// update last retry timestamp
		tracker.LastRetry = retryTime
		internalTrackers = append(internalTrackers, tracker.Tracker)
	}

	// remove trackers for finalized ballots
	for _, ballot := range trackersToRemove {
		ob.removeInternalInboundTracker(ballot)
	}

	return internalTrackers
}

// AddInternalInboundTracker adds an internal inbound tracker for given inbound vote.
func (ob *Observer) AddInternalInboundTracker(msg *crosschaintypes.MsgVoteInbound) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		timeNow = time.Now()
		ballot  = msg.Digest()
		tracker = &InternalInboundTracker{
			CreatedAt: timeNow,
			LastRetry: timeNow,
			Tracker:   msg.InboundTracker(),
		}
	)

	if _, found := ob.internalInboundTrackers[ballot]; !found {
		ob.internalInboundTrackers[ballot] = tracker
		ob.logger.Inbound.Info().
			Str(logs.FieldBallotIndex, ballot).
			Str(logs.FieldTx, tracker.Tracker.TxHash).
			Str(logs.FieldCoinType, tracker.Tracker.CoinType.String()).
			Msg("added internal inbound tracker")
	}
}

// removeInternalInboundTracker removes an internal inbound tracker for given ballot.
func (ob *Observer) removeInternalInboundTracker(ballot string) {
	if tracker, found := ob.internalInboundTrackers[ballot]; found {
		delete(ob.internalInboundTrackers, ballot)
		ob.logger.Inbound.Info().
			Str(logs.FieldBallotIndex, ballot).
			Str(logs.FieldTx, tracker.Tracker.TxHash).
			Str(logs.FieldCoinType, tracker.Tracker.CoinType.String()).
			Stringer("time_elapsed", time.Since(tracker.CreatedAt)).
			Msg("removed internal inbound tracker")
	}
}

func (ob *Observer) WatchMonitoringError(
	ctx context.Context,
	monitorErrCh <-chan zetaerrors.ErrTxMonitor,
	zetaHash string,
) {
	logger := ob.logger.Inbound

	select {
	case monitorErr := <-monitorErrCh:
		if monitorErr.Err != nil {
			logger.Error().
				Err(monitorErr).
				Str(logs.FieldTx, monitorErr.Msg.InboundHash).
				Str(logs.FieldZetaTx, monitorErr.ZetaTxHash).
				Str(logs.FieldBallotIndex, monitorErr.Msg.Digest()).
				Msg("error monitoring inbound vote")

			ob.AddInternalInboundTracker(&monitorErr.Msg)
		}
	case <-ctx.Done():
		logger.Debug().
			Str(logs.FieldZetaTx, zetaHash).
			Msg("no monitoring error received, the inbound vote likely succeeded")
	}
}
