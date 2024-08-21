package rpc

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

const (
	// defaultPageLimit is the default number of signatures to fetch in one GetSignaturesForAddressWithOpts call
	DefaultPageLimit = 1000

	// RPCAlertLatency is the default threshold for RPC latency to be considered unhealthy and trigger an alert.
	// The 'HEALTH_CHECK_SLOT_DISTANCE' is default to 150 slots, which is 150 * 0.4s = 60s
	RPCAlertLatency = time.Duration(60) * time.Second
)

// GetFirstSignatureForAddress searches the first signature for the given address.
// Note: make sure that the rpc provider used has enough transaction history.
func GetFirstSignatureForAddress(
	ctx context.Context,
	client interfaces.SolanaRPCClient,
	address solana.PublicKey,
	pageLimit int,
) (solana.Signature, error) {
	// search backwards until we find the first signature
	var lastSignature solana.Signature
	for {
		fetchedSignatures, err := client.GetSignaturesForAddressWithOpts(
			ctx,
			address,
			&rpc.GetSignaturesForAddressOpts{
				Limit:      &pageLimit,
				Before:     lastSignature, // exclusive
				Commitment: rpc.CommitmentFinalized,
			},
		)
		if err != nil {
			return solana.Signature{}, errors.Wrapf(
				err,
				"error GetSignaturesForAddressWithOpts for address %s",
				address,
			)
		}

		// no more signatures, stop searching
		if len(fetchedSignatures) == 0 {
			break
		}

		// update last signature for next search
		lastSignature = fetchedSignatures[len(fetchedSignatures)-1].Signature
	}

	// there is no signature for the given address
	if lastSignature.IsZero() {
		return lastSignature, errors.Errorf("no signatures found for address %s", address)
	}

	return lastSignature, nil
}

// GetSignaturesForAddressUntil searches for signatures for the given address until the given signature (exclusive).
// Note: make sure that the rpc provider used has enough transaction history.
func GetSignaturesForAddressUntil(
	ctx context.Context,
	client interfaces.SolanaRPCClient,
	address solana.PublicKey,
	untilSig solana.Signature,
	pageLimit int,
) ([]*rpc.TransactionSignature, error) {
	var lastSignature solana.Signature
	var allSignatures []*rpc.TransactionSignature

	// make sure that the 'untilSig' exists to prevent undefined behavior on GetSignaturesForAddressWithOpts
	_, err := client.GetTransaction(
		ctx,
		untilSig,
		&rpc.GetTransactionOpts{Commitment: rpc.CommitmentFinalized},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error GetTransaction for untilSig %s", untilSig)
	}

	// search backwards until we hit the 'untilSig' signature
	for {
		fetchedSignatures, err := client.GetSignaturesForAddressWithOpts(
			ctx,
			address,
			&rpc.GetSignaturesForAddressOpts{
				Limit:      &pageLimit,
				Before:     lastSignature, // exclusive
				Until:      untilSig,      // exclusive
				Commitment: rpc.CommitmentFinalized,
			},
		)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"error GetSignaturesForAddressWithOpts for address %s",
				address,
			)
		}

		// no more signatures, stop searching
		if len(fetchedSignatures) == 0 {
			break
		}

		// update last signature for next search
		lastSignature = fetchedSignatures[len(fetchedSignatures)-1].Signature

		// append fetched signatures
		allSignatures = append(allSignatures, fetchedSignatures...)
	}

	return allSignatures, nil
}

// CheckRPCStatus checks the RPC status of the solana chain
func CheckRPCStatus(
	ctx context.Context,
	client interfaces.SolanaRPCClient,
	alertLatency time.Duration,
	logger zerolog.Logger,
) error {
	// query solana health (always return "ok" unless --trusted-validator is provided)
	_, err := client.GetHealth(ctx)
	if err != nil {
		return errors.Wrap(err, "GetHealth error: RPC down?")
	}

	// query latest slot
	slot, err := client.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return errors.Wrap(err, "GetSlot error: RPC down?")
	}

	// query latest block time
	blockTime, err := client.GetBlockTime(ctx, slot)
	if err != nil {
		return errors.Wrap(err, "GetBlockTime error: RPC down?")
	}

	// use default alert latency if not provided
	if alertLatency <= 0 {
		alertLatency = RPCAlertLatency
	}

	// latest block should not be too old
	elapsedTime := time.Since(blockTime.Time())
	if elapsedTime > alertLatency {
		return errors.Errorf(
			"Latest slot %d is %.0fs old, RPC stale or chain stuck (check explorer)?",
			slot,
			elapsedTime.Seconds(),
		)
	}

	logger.Info().
		Msgf("RPC Status [OK]: latest slot %d, timestamp %s (%.0fs ago)", slot, blockTime.String(), elapsedTime.Seconds())
	return nil
}
