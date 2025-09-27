package signer

import (
	"context"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// reportOutboundTracker queries the tx and sends its digest to the outbound tracker
// for further processing by the Observer.
func (s *Signer) reportOutboundTracker(ctx context.Context, nonce uint64, digest string) error {
	metrics.NumTrackerReporters.WithLabelValues(s.Chain().Name).Inc()
	defer metrics.NumTrackerReporters.WithLabelValues(s.Chain().Name).Dec()

	// approx Sui checkpoint interval
	const interval = 3 * time.Second

	// some sanity timeout
	const maxTimeout = time.Minute

	// prepare logger
	logger := zerolog.Ctx(ctx).With().Str(logs.FieldTx, digest).Logger()

	alreadySet := s.SetBeingReportedFlag(digest)
	if alreadySet {
		logger.Info().Msg("outbound is already being observed for the tracker")
		return nil
	}

	start := time.Now()
	attempts := 0

	// request tx with effects as we want to see its status
	req := models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
	}

	defer s.ClearBeingReportedFlag(digest)

	for {
		switch {
		case time.Since(start) > maxTimeout:
			return errors.Errorf("timeout reached (%s)", maxTimeout.String())
		case attempts == 0:
			// best case we'd be able to report the tx ~immediately
			time.Sleep(interval / 2)
		default:
			time.Sleep(interval)
		}
		attempts++

		res, err := s.client.SuiGetTransactionBlock(ctx, req)
		switch {
		case ctx.Err() != nil:
			return errors.Wrap(ctx.Err(), "failed to get transaction block")
		case err != nil:
			logger.Error().Err(err).Msg("failed to get transaction block")
			continue
		case res.Effects.Status.Status == client.TxStatusFailure:
			// failed outbound should be ignored as it cannot increment the gateway nonce.
			// Sui transaction status is one of ["success", "failure"]
			// see: https://github.com/MystenLabs/sui/blob/615516edb0ed55e45d599f042f9570b493ce9643/crates/sui-json-rpc-types/src/sui_transaction.rs#L1345
			return errors.Errorf("tx failed with error: %s", res.Effects.Status.Error)
		case res.Effects.Status.Status == client.TxStatusSuccess && res.Checkpoint != "":
			return s.postTrackerVote(ctx, nonce, digest)
		default:
			// otherwise, hold on until the tx status can be clearly determined.
			// we prefer missed tracker hash over potentially invalid hash.
			continue
		}
	}
}

// note that at this point we don't care whether tx was successful or not.
func (s *Signer) postTrackerVote(ctx context.Context, nonce uint64, digest string) error {
	_, err := s.zetacore.PostOutboundTracker(ctx, s.Chain().ChainId, nonce, digest)
	return err
}
