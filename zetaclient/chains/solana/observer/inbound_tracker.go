package observer

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	solanarpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
)

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
		txResult, err := solanarpc.GetTransaction(ctx, ob.solClient, signature)
		switch {
		case errors.Is(err, solanarpc.ErrUnsupportedTxVersion):
			ob.Logger().Inbound.Warn().Stringer("tx.signature", signature).Msg("skip inbound tracker hash")
			continue
		case err != nil:
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
