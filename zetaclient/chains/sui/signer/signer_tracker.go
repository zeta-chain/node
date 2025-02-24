package signer

import (
	"context"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// reportOutboundTracker queries the tx and sends its digest to the outbound tracker
// for further processing by the Observer.
func (s *Signer) reportOutboundTracker(ctx context.Context, nonce uint64, digest string) error {
	// approx Sui checkpoint interval
	const interval = 3 * time.Second

	// some sanity timeout
	const maxTimeout = time.Minute

	logger := zerolog.Ctx(ctx)

	alreadySet := s.SetBeingReportedFlag(digest)
	if alreadySet {
		logger.Info().Msg("Outbound is already being observed for the tracker")
		return nil
	}

	start := time.Now()
	attempts := 0

	req := models.SuiGetTransactionBlockRequest{Digest: digest}

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
		case err != nil:
			logger.Error().Err(err).Msg("Failed to get transaction block")
			continue
		case res.Checkpoint == "":
			// should not happen
			logger.Error().Msg("Checkpoint is empty")
			continue
		default:
			return s.postTrackerVote(ctx, nonce, digest)
		}
	}
}

// note that at this point we don't care whether tx was successful or not.
func (s *Signer) postTrackerVote(ctx context.Context, nonce uint64, digest string) error {
	_, err := s.zetacore.PostOutboundTracker(ctx, s.Chain().ChainId, nonce, digest)
	return err
}
