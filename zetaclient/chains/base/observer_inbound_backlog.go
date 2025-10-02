package base

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// GetInboundTrackersWithBacklog returns combined inbound trackers from zetacore and the backlog.
func (ob *Observer) GetInboundTrackersWithBacklog(ctx context.Context) ([]crosschaintypes.InboundTracker, error) {
	trackers, err := ob.ZetacoreClient().GetInboundTrackersForChain(ctx, ob.chain.ChainId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get inbound trackers")
	}

	// create a map to look up trackers by inbound hash
	trackersSet := make(map[string]bool)
	for _, tracker := range trackers {
		trackersSet[strings.ToLower(tracker.TxHash)] = true
	}

	ob.mu.Lock()
	defer ob.mu.Unlock()

	// combine trackers with the backlog and skip those that are already finalized
	finalizedBallots := make([]string, 0)
	for ballot, tracker := range ob.failedInboundBacklog {
		if _, err := ob.ZetacoreClient().GetCctxByHash(ctx, ballot); err == nil {
			finalizedBallots = append(finalizedBallots, ballot)
			continue
		}

		if _, found := trackersSet[strings.ToLower(tracker.TxHash)]; !found {
			trackers = append(trackers, tracker)
		}
	}

	// remove finalized ballots from the backlog
	for _, ballot := range finalizedBallots {
		ob.removeFailedInbound(ballot)
	}

	return trackers, nil
}

// AddFailedInbound adds a failed inbound to the backlog.
func (ob *Observer) AddFailedInbound(msg *crosschaintypes.MsgVoteInbound) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var (
		ballot  = msg.Digest()
		tracker = msg.InboundTracker()
	)

	if _, found := ob.failedInboundBacklog[ballot]; !found {
		ob.failedInboundBacklog[ballot] = tracker
		ob.logger.Inbound.Info().
			Str(logs.FieldBallotIndex, ballot).
			Str(logs.FieldTx, tracker.TxHash).
			Str(logs.FieldCoinType, tracker.CoinType.String()).
			Msg("added failed inbound to backlog")
	}
}

// removeFailedInbound removes a failed inbound from the backlog.
func (ob *Observer) removeFailedInbound(ballot string) {
	if tracker, found := ob.failedInboundBacklog[ballot]; found {
		delete(ob.failedInboundBacklog, ballot)
		ob.logger.Inbound.Info().
			Str(logs.FieldBallotIndex, ballot).
			Str(logs.FieldTx, tracker.TxHash).
			Str(logs.FieldCoinType, tracker.CoinType.String()).
			Msg("removed failed inbound from backlog")
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

			ob.AddFailedInbound(&monitorErr.Msg)
		}
	case <-ctx.Done():
		logger.Debug().
			Str(logs.FieldZetaTx, zetaHash).
			Msg("no monitoring error received, the inbound vote likely succeeded")
	}
}
