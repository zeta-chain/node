package observer

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

// WatchInboundTracker watches zetacore for Solana inbound trackers
func (ob *Observer) WatchInboundTracker() {
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchInboundTracker_%d", ob.Chain().ChainId),
		ob.GetChainParams().InboundTicker,
	)
	if err != nil {
		ob.Logger().Inbound.Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	ob.Logger().Inbound.Info().Msgf("WatchInboundTracker started for chain %d", ob.Chain().ChainId)
	for {
		select {
		case <-ticker.C():
			if !ob.AppContext().IsInboundObservationEnabled(ob.GetChainParams()) {
				continue
			}
			err := ob.ProcessInboundTrackers()
			if err != nil {
				ob.Logger().Inbound.Error().
					Err(err).
					Msgf("WatchInboundTracker: error ProcessInboundTrackers for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.Logger().Inbound)
		case <-ob.StopChannel():
			ob.Logger().Inbound.Info().Msgf("WatchInboundTracker stopped for chain %d", ob.Chain().ChainId)
			return
		}
	}
}

// ProcessInboundTrackers processes inbound trackers
func (ob *Observer) ProcessInboundTrackers() error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.ZetacoreClient().GetInboundTrackersForChain(chainID)
	if err != nil {
		return err
	}

	// process inbound trackers
	for _, tracker := range trackers {
		signature := solana.MustSignatureFromBase58(tracker.TxHash)
		txResult, err := ob.solClient.GetTransaction(context.TODO(), signature, &rpc.GetTransactionOpts{
			Commitment: rpc.CommitmentFinalized,
		})
		if err != nil {
			return errors.Wrapf(err, "error GetTransaction for chain %d sig %s", chainID, signature)
		}

		// filter inbound event and vote
		err = ob.FilterInboundEventAndVote(txResult)
		if err != nil {
			return errors.Wrapf(err, "error FilterInboundEventAndVote for chain %d sig %s", chainID, signature)
		}
	}

	return nil
}
