package signer

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// reportToOutboundTracker launch a go routine with timeout to check for tx confirmation;
// it reports tx to outbound tracker only if it's confirmed by the Solana network.
func (signer *Signer) reportToOutboundTracker(
	ctx context.Context,
	zetacoreClient interfaces.ZetacoreClient,
	chainID int64,
	nonce uint64,
	txSig solana.Signature,
	logger zerolog.Logger,
) {
	// prepare logger
	logger = logger.With().
		Str(logs.FieldMethod, "reportToOutboundTracker").
		Int64(logs.FieldChain, chainID).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txSig.String()).
		Logger()

	// set being reported flag to avoid duplicate reporting
	alreadySet := signer.Signer.SetBeingReportedFlag(txSig.String())
	if alreadySet {
		logger.Info().Msg("Outbound is already reported to tracker")
		return
	}

	// launch a goroutine to monitor tx confirmation status
	bg.Work(ctx, func(ctx context.Context) error {
		metrics.NumTrackerReporters.WithLabelValues(signer.Chain().Name).Inc()

		defer func() {
			metrics.NumTrackerReporters.WithLabelValues(signer.Chain().Name).Dec()
			signer.Signer.ClearBeingReportedFlag(txSig.String())
		}()

		start := time.Now()
		for {
			// Solana block time is 0.4~0.8 seconds; wait 5 seconds between each check
			time.Sleep(5 * time.Second)

			// give up if we know the tx is too old and already expired
			if time.Since(start) > solanaTransactionTimeout {
				logger.Info().Msg("Outbound is expired")
				return nil
			}

			// query tx using optimistic commitment level "confirmed"
			tx, err := signer.client.GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
				// commitment "processed" seems to be a better choice but it's not supported
				// see: https://solana.com/docs/rpc/http/gettransaction
				Commitment: rpc.CommitmentConfirmed,
			})
			if err != nil {
				continue
			}

			// exit goroutine if tx failed.
			if tx.Meta.Err != nil {
				// unlike Ethereum, Solana doesn't have protocol-level nonce; the nonce is enforced by the gateway program.
				// a failed outbound (e.g. signature err, balance err) will never be able to increment the gateway program nonce.
				// a good/valid candidate of outbound tracker hash must come with a successful tx.
				logger.Warn().Any("tx_error", tx.Meta.Err).Msg("Outbound is failed")
				return nil
			}

			// report outbound hash to zetacore
			zetaHash, err := zetacoreClient.PostOutboundTracker(ctx, chainID, nonce, txSig.String())
			if err != nil {
				logger.Err(err).Msg("Error adding outbound to tracker")
			} else if zetaHash != "" {
				logger.Info().Str(logs.FieldZetaTx, zetaHash).Msg("Added outbound to tracker")
			} else {
				// exit goroutine until the tracker contains the hash (reported by either this or other signers)
				logger.Info().Msg("Outbound now exists in tracker")
				return nil
			}
		}
	}, bg.WithName("TrackerReporterSolana"), bg.WithLogger(logger))
}
