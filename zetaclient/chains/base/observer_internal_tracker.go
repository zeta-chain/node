package base

import (
	"context"

	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// GetInboundInternalTrackers returns internal inbound trackers
func (ob *Observer) GetInboundInternalTrackers(ctx context.Context) []crosschaintypes.InboundTracker {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		finalizedBallots = make([]string, 0)
		internalTrackers = make([]crosschaintypes.InboundTracker, 0, len(ob.internalInboundTrackers))
	)

	// collect up to MaxInternalTrackersPerScan trackers
	for ballot, tracker := range ob.internalInboundTrackers {
		// skip those that are already finalized
		if _, err := ob.ZetacoreClient().GetCctxByHash(ctx, ballot); err == nil {
			finalizedBallots = append(finalizedBallots, ballot)
			continue
		}
		internalTrackers = append(internalTrackers, tracker)
	}

	// remove trackers for finalized ballots
	for _, ballot := range finalizedBallots {
		ob.removeInternalInboundTracker(ballot)
	}

	return internalTrackers
}

// AddInternalInboundTracker adds an internal inbound tracker for given inbound vote.
func (ob *Observer) AddInternalInboundTracker(msg *crosschaintypes.MsgVoteInbound) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		ballot  = msg.Digest()
		tracker = msg.InboundTracker()
	)

	if _, found := ob.internalInboundTrackers[ballot]; !found {
		ob.internalInboundTrackers[ballot] = tracker
		ob.logger.Inbound.Info().
			Str(logs.FieldBallotIndex, ballot).
			Str(logs.FieldTx, tracker.TxHash).
			Str(logs.FieldCoinType, tracker.CoinType.String()).
			Msg("added internal inbound tracker")
	}
}

// removeInternalInboundTracker removes an internal inbound tracker for given ballot.
func (ob *Observer) removeInternalInboundTracker(ballot string) {
	if tracker, found := ob.internalInboundTrackers[ballot]; found {
		delete(ob.internalInboundTrackers, ballot)
		ob.logger.Inbound.Info().
			Str(logs.FieldBallotIndex, ballot).
			Str(logs.FieldTx, tracker.TxHash).
			Str(logs.FieldCoinType, tracker.CoinType.String()).
			Msg("removed internal inbound tracker")
	}
}

func (ob *Observer) handleMonitoringError(
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
