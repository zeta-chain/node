// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/bg"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// reportToOutboundTracker reports outboundHash to tracker only when tx receipt is available
func (signer *Signer) reportToOutboundTracker(
	ctx context.Context,
	zetaRepo *zrepo.ZetaRepo,
	chainID int64,
	nonce uint64,
	outboundHash string,
	logger zerolog.Logger,
) {
	// prepare logger
	logger = logger.With().
		Int64(logs.FieldChain, chainID).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, outboundHash).
		Logger()

	// set being reported flag to avoid duplicate reporting
	alreadySet := signer.SetBeingReportedFlag(outboundHash)
	if alreadySet {
		logger.Info().Msg("outbound is being reported to tracker")
		return
	}

	// launch a goroutine to monitor tx confirmation status
	bg.Work(ctx, func(ctx context.Context) error {
		metrics.NumTrackerReporters.WithLabelValues(signer.Chain().Name).Inc()

		defer func() {
			metrics.NumTrackerReporters.WithLabelValues(signer.Chain().Name).Dec()
			signer.ClearBeingReportedFlag(outboundHash)
		}()

		// try monitoring tx inclusion status for 20 minutes
		tStart := time.Now()
		for {
			// take a rest between each check
			time.Sleep(10 * time.Second)

			// give up (forget about the tx) after 20 minutes of monitoring, there are 2 reasons:
			// 1. the gas stability pool should have kicked in and replaced the tx by then.
			// 2. even if there is a chance that the tx is included later, most likely it's going to be a false tx hash (either replaced or dropped).
			// 3. we prefer missed tx hash over potentially invalid txhash.
			if time.Since(tStart) > common.OutboundInclusionTimeout {
				logger.Info().Msg("timeout waiting outbound inclusion")
				return nil
			}

			// stop if the CCTX is already finalized for optimization purposes:
			// 1. all monitoring goroutines should stop and release resources if the CCTX is finalized
			// 2. especially reduces the lifetime of goroutines that monitor "nonce too low" tx hashes
			cctx, err := zetaRepo.GetCCTX(ctx, nonce)
			if err != nil {
				logger.Err(err).Send()
			} else if !crosschainkeeper.IsPending(cctx) {
				logger.Info().Msg("CCTX is already finalized")
				return nil
			}

			// check tx confirmation status
			confirmed, err := signer.evmClient.IsTxConfirmed(ctx, outboundHash, common.ReorgProtectBlockCount)
			if err != nil {
				logger.Err(err).Msg("unable to check confirmation status of outbound")
				continue
			}
			if !confirmed {
				continue
			}

			// report outbound hash to tracker
			zhash, err := zetaRepo.PostOutboundTracker(ctx, logger, nonce, outboundHash)
			if zhash == "" && err == nil {
				// exit goroutine only when the tracker contains the hash (reported by either this
				// or other signers)
				logger.Info().Msg("outbound now exists in tracker")
				return nil
			}
		}
	}, bg.WithName("TrackerReporterEVM"), bg.WithLogger(logger))
}
