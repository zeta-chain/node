// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/bg"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// reportToOutboundTracker reports outboundHash to tracker only when tx receipt is available
func (signer *Signer) reportToOutboundTracker(
	ctx context.Context,
	zetacoreClient interfaces.ZetacoreClient,
	chainID int64,
	nonce uint64,
	outboundHash string,
	logger zerolog.Logger,
) {
	// prepare logger
	logger = logger.With().
		Str(logs.FieldMethod, "reportToOutboundTracker").
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
		defer func() {
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
				logger.Info().Msgf("timeout waiting outbound inclusion")
				return nil
			}

			// check tx confirmation status
			confirmed, err := signer.client.IsTxConfirmed(ctx, outboundHash, common.ReorgProtectBlockCount)
			if err != nil {
				logger.Err(err).Msg("unable to check confirmation status of outbound")
				continue
			}
			if !confirmed {
				continue
			}

			// stop if the cctx is already finalized
			cctx, err := zetacoreClient.GetCctxByNonce(ctx, chainID, nonce)
			if err != nil {
				logger.Err(err).Msg("unable to get cctx for outbound")
				continue
			} else if !crosschainkeeper.IsPending(cctx) {
				logger.Info().Msg("cctx is finalized")
				return nil
			}

			// report outbound hash to tracker
			zetaHash, err := zetacoreClient.PostOutboundTracker(ctx, chainID, nonce, outboundHash)
			if err != nil {
				logger.Err(err).Msg("error adding outbound to tracker")
			} else if zetaHash != "" {
				logger.Info().Msgf("added outbound to tracker; zeta txhash %s", zetaHash)
			} else {
				// exit goroutine until the tracker contains the hash (reported by either this or other signers)
				logger.Info().Msg("outbound now exists in tracker")
				return nil
			}
		}
	}, bg.WithName("TrackerReporterEVM"), bg.WithLogger(logger))
}
