package observer

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	zctx "github.com/zeta-chain/node/zetaclient/context"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// WatchInboundTracker watches zetacore for Solana inbound trackers
func (ob *Observer) WatchInboundTracker(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchInboundTracker_%d", ob.Chain().ChainId),
		ob.GetChainParams().InboundTicker,
	)
	if err != nil {
		ob.Logger().Inbound.Err(err).Msg("error creating ticker")
		return err
	}
	defer ticker.Stop()

	ob.Logger().Inbound.Info().Msgf("WatchInboundTracker started for chain %d", ob.Chain().ChainId)
	for {
		select {
		case <-ticker.C():
			if !app.IsInboundObservationEnabled() {
				continue
			}
			err := ob.ProcessInboundTrackers(ctx)
			if err != nil {
				ob.Logger().Inbound.Error().
					Err(err).
					Msgf("WatchInboundTracker: error ProcessInboundTrackers for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.Logger().Inbound)
		case <-ob.StopChannel():
			ob.Logger().Inbound.Info().Msgf("WatchInboundTracker stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// ProcessInboundTrackers processes inbound trackers
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.ZetacoreClient().GetInboundTrackersForChain(ctx, chainID)
	if err != nil {
		return err
	}

	// process inbound trackers
	for _, tracker := range trackers {
		signature := solana.MustSignatureFromBase58(tracker.TxHash)
		txResult, err := ob.solClient.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{
			Commitment: rpc.CommitmentFinalized,
		})
		if err != nil {
			return errors.Wrapf(err, "error GetTransaction for chain %d sig %s", chainID, signature)
		}

		// filter inbound events and vote
		err = ob.FilterInboundEventsAndVote(ctx, txResult)
		if err != nil {
			return errors.Wrapf(err, "error FilterInboundEventAndVote for chain %d sig %s", chainID, signature)
		}
	}

	return nil
}
