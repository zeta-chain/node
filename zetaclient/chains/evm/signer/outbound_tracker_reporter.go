// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/rpc"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/logs"
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
	alreadySet := signer.Signer.SetBeingReportedFlag(outboundHash)
	if alreadySet {
		logger.Info().
			Msgf("outbound %s for chain %d nonce %d is being reported", outboundHash, chainID, nonce)
		return
	}

	// launch a goroutine to monitor tx confirmation status
	go func() {
		defer func() {
			signer.Signer.ClearBeingReportedFlag(outboundHash)
		}()

		// try monitoring tx inclusion status for 20 minutes
		tStart := time.Now()
		for {
			// take a rest between each check
			time.Sleep(10 * time.Second)

			// give up (forget about the tx) after 20 minutes of monitoring, there are 2 reasons:
			// 1. the gas stability pool should have kicked in and replaced the tx by then.
			// 2. even if there is a chance that the tx is included later, most likely it's going to be a false tx hash (either replaced or dropped).
			if time.Since(tStart) > evm.OutboundInclusionTimeout {
				logger.Info().
					Msgf("timeout waiting outbound %s inclusion for chain %d nonce %d", outboundHash, chainID, nonce)
				return
			}

			// check tx confirmation status
			confirmed, err := rpc.IsTxConfirmed(ctx, signer.client, outboundHash, evm.ReorgProtectBlockCount)
			if err != nil {
				logger.Err(err).
					Msgf("unable to check confirmation status for chain %d nonce %d outbound %s", chainID, nonce, outboundHash)
			}
			if !confirmed {
				continue
			}

			// report outbound hash to tracker
			zetaHash, err := zetacoreClient.AddOutboundTracker(ctx, chainID, nonce, outboundHash, nil, "", -1)
			if err != nil {
				logger.Err(err).
					Msgf("error adding outbound %s to tracker for chain %d nonce %d", outboundHash, chainID, nonce)
			} else if zetaHash != "" {
				logger.Info().Msgf("added outbound %s to tracker for chain %d nonce %d; zeta txhash %s", outboundHash, chainID, nonce, zetaHash)
			} else {
				// exit goroutine until the tracker contains the hash (reported by either this or other signers)
				logger.Info().Msgf("outbound %s now exists in tracker for chain %d nonce %d", outboundHash, chainID, nonce)
				return
			}
		}
	}()
}
