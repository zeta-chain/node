package observer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// ProcessInboundTrackers processes inbound trackers
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	trackers, err := ob.ZetaRepo().GetInboundTrackers(ctx)
	if err != nil {
		return err
	}

	return ob.observeInboundTrackers(ctx, trackers, false)
}

// ProcessInternalTrackers processes internal inbound trackers
func (ob *Observer) ProcessInternalTrackers(ctx context.Context) error {
	trackers := ob.GetInboundInternalTrackers(ctx)
	if len(trackers) > 0 {
		ob.Logger().Inbound.Info().Int("total_count", len(trackers)).Msg("processing internal trackers")
	}

	return ob.observeInboundTrackers(ctx, trackers, true)
}

// observeInboundTrackers observes given inbound trackers
func (ob *Observer) observeInboundTrackers(
	ctx context.Context,
	trackers []types.InboundTracker,
	isInternal bool,
) error {
	// take at most MaxInternalTrackersPerScan for each scan
	if len(trackers) > config.MaxInboundTrackersPerScan {
		trackers = trackers[:config.MaxInboundTrackersPerScan]
	}

	for _, tracker := range trackers {
		ob.logger.Inbound.Info().
			Str(logs.FieldTx, tracker.TxHash).
			Stringer(logs.FieldCoinType, tracker.CoinType).
			Bool("is_internal", isInternal).
			Msg("processing inbound tracker")
		if _, err := ob.CheckReceiptAndPostVoteForBtcTxHash(ctx, tracker.TxHash, true); err != nil {
			return err
		}
	}

	return nil
}

// CheckReceiptAndPostVoteForBtcTxHash checks the receipt for a btc tx hash
func (ob *Observer) CheckReceiptAndPostVoteForBtcTxHash(ctx context.Context, txHash string, vote bool) (string, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return "", errors.Wrap(err, "error parsing btc tx hash")
	}

	tx, err := ob.bitcoinClient.GetRawTransactionVerbose(ctx, hash)
	if err != nil {
		return "", errors.Wrap(err, "error getting btc raw tx verbose")
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return "", errors.Wrap(err, "error parsing btc block hash")
	}

	blockVb, err := ob.bitcoinClient.GetBlockVerbose(ctx, blockHash)
	if err != nil {
		return "", errors.Wrap(err, "error getting btc block verbose")
	}

	if len(blockVb.Tx) <= 1 {
		return "", fmt.Errorf("block %d has no transactions", blockVb.Height)
	}

	tss, err := ob.ZetaRepo().GetBTCTSSAddress(ctx)
	if err != nil {
		return "", err
	}

	// check confirmation
	// #nosec G115 block height always positive
	if !ob.IsBlockConfirmedForInboundSafe(uint64(blockVb.Height)) {
		return "", fmt.Errorf("block %d is not confirmed yet", blockVb.Height)
	}

	// #nosec G115 always positive
	event, err := GetBtcEventWithWitness(
		ctx,
		ob.bitcoinClient,
		*tx,
		tss,
		uint64(blockVb.Height),
		ob.logger.Inbound,
		ob.netParams,
		common.CalcDepositorFee,
	)
	if err != nil {
		return "", errors.Wrap(err, "error getting btc event")
	}

	if event == nil {
		return "", errors.New("no btc deposit event found")
	}

	msg := ob.GetInboundVoteFromBtcEvent(event)
	if msg == nil {
		return "", errors.New("no message built for btc sent to TSS")
	}

	if !vote {
		return msg.Digest(), nil
	}

	metrics.InboundObservationsTrackerTotal.WithLabelValues(ob.Chain().Name, strconv.FormatBool(false)).Inc()
	return ob.ZetaRepo().VoteInbound(ctx,
		ob.logger.Inbound,
		msg,
		zetacore.PostVoteInboundExecutionGasLimit,
		ob.WatchMonitoringError,
	)
}
