package signer

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

const (
	// SolanaTransactionTimeout is the timeout for waiting for an outbound to be confirmed
	// Transaction referencing a blockhash older than 150 blocks will expire and be rejected by Solana.
	SolanaTransactionTimeout = 2 * time.Minute
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
	// set being reported flag to avoid duplicate reporting
	alreadySet := signer.Signer.SetBeingReportedFlag(txSig.String())
	if alreadySet {
		logger.Info().
			Msgf("reportToOutboundTracker: outbound %s for chain %d nonce %d is being reported", txSig, chainID, nonce)
		return
	}
	signer.Logger().
		Std.Info().
		Msgf("reportToOutboundTracker: reporting outbound %s for chain %d nonce %d", txSig, chainID, nonce)

	// launch a goroutine to monitor tx confirmation status
	go func() {
		defer func() {
			signer.Signer.ClearBeingReportedFlag(txSig.String())
		}()

		start := time.Now()
		for {
			// Solana block time is 0.4~0.8 seconds; wait 5 seconds between each check
			time.Sleep(5 * time.Second)

			// give up if we know the tx is too old and already expired
			if time.Since(start) > SolanaTransactionTimeout {
				logger.Info().
					Msgf("reportToOutboundTracker: outbound %s expired for chain %d nonce %d", txSig, chainID, nonce)
				break
			}

			// get confirmed tx with metadata
			tx, err := signer.client.GetConfirmedTransactionWithOpts(ctx, txSig, &rpc.GetTransactionOpts{
				// commitment "CommitmentProcessed" seems to be a better choice but it's not supported
				// see: https://solana.com/docs/rpc/http/gettransaction
				Commitment: rpc.CommitmentConfirmed,
			})
			if err != nil {
				continue
			}

			// successfully included in a block; report and exit goroutine
			if tx.Meta.Err == nil {
				zetaHash, err := zetacoreClient.AddOutboundTracker(ctx, chainID, nonce, txSig.String(), nil, "", -1)
				if err != nil {
					logger.Err(err).
						Msgf("reportToOutboundTracker: error adding outbound %s for chain %d nonce %d", txSig, chainID, nonce)
				} else if zetaHash != "" {
					logger.Info().Msgf("reportToOutboundTracker: added outbound %s for chain %d nonce %d; zeta txhash %s", txSig, chainID, nonce, zetaHash)
				} else {
					// stop if the tracker already contains the outboundHash
					logger.Info().Msgf("reportToOutboundTracker: outbound %s already in tracker for chain %d nonce %d", txSig, chainID, nonce)
					return
				}
			}

			// it's included by failed (likely competing txs succeeded). exit goroutine.
			// FIXME: we should report this failed tx ONLY IF it failed not due to nonce mismatch error
			// FIXME: add a check for nonce mismatch error
			logger.Warn().Msgf("tx %s failed: %v", txSig, tx.Meta.Err)
			return
		}
	}()
}
