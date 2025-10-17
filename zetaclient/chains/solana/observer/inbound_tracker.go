package observer

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// ProcessInboundTrackers processes inbound trackers
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	trackers, err := ob.ZetaRepo().GetInboundTrackers(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get inbound trackers")
	}

	return ob.observeInboundTrackers(ctx, trackers, false)
}

// ProcessInternalTrackers processes internal inbound trackers
func (ob *Observer) ProcessInternalTrackers(ctx context.Context) error {
	trackers := ob.GetInboundInternalTrackers(ctx, time.Now())
	if len(trackers) > 0 {
		ob.Logger().Inbound.Info().Int("total_count", len(trackers)).Msg("processing internal inbound trackers")
	}

	return ob.observeInboundTrackers(ctx, trackers, true)
}

// observeInboundTrackers observes given inbound trackers
func (ob *Observer) observeInboundTrackers(
	ctx context.Context,
	trackers []types.InboundTracker,
	isInternal bool,
) error {
	chainID := ob.Chain().ChainId

	// take at most MaxInternalTrackersPerScan for each scan
	if len(trackers) > config.MaxInboundTrackersPerScan {
		trackers = trackers[:config.MaxInboundTrackersPerScan]
	}

	// process inbound trackers
	for _, tracker := range trackers {
		signature := solana.MustSignatureFromBase58(tracker.TxHash)
		txResult, err := ob.solanaRepo.GetTransaction(ctx, signature)
		switch {
		case errors.Is(err, repo.ErrUnsupportedTxVersion):
			ob.Logger().Inbound.Warn().
				Stringer(logs.FieldTx, signature).
				Bool("is_internal", isInternal).
				Msg("skip inbound tracker hash")
			continue
		case err != nil:
			return errors.Wrapf(err, "error GetTransaction for chain %d sig %s", chainID, signature)
		}

		// filter inbound events
		events, err := FilterInboundEvents(txResult, ob.gatewayID, ob.Chain().ChainId, ob.Logger().Inbound)
		if err != nil {
			return errors.Wrapf(err, "error FilterInboundEvents for chain %d sig %s", chainID, signature)
		}

		// vote inbound events
		if err := ob.VoteInboundEvents(ctx, events); err != nil {
			// return error to retry this transaction
			return errors.Wrapf(err, "error VoteInboundEvents for chain %d sig %s", chainID, signature)
		}
	}

	return nil
}
