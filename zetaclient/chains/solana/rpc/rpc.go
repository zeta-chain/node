package rpc

import (
	"context"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

const (
	// defaultPageLimit is the default number of signatures to fetch in one GetSignaturesForAddressWithOpts call
	DefaultPageLimit = 1000

	// RPCAlertLatency is the default threshold for RPC latency to be considered unhealthy and trigger an alert.
	// The 'HEALTH_CHECK_SLOT_DISTANCE' is default to 150 slots, which is 150 * 0.4s = 60s
	RPCAlertLatency = time.Duration(60) * time.Second

	// see: https://github.com/solana-labs/solana/blob/master/rpc/src/rpc.rs#L7276
	errorCodeUnsupportedTransactionVersion = "-32015"
)

// ErrUnsupportedTxVersion is returned when the transaction version is not supported by zetaclient
var ErrUnsupportedTxVersion = errors.New("unsupported tx version")

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
	_, err := GetTransaction(ctx, client, untilSig)
	if err != nil && !errors.Is(err, ErrUnsupportedTxVersion) {
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

// GetTransaction fetches a transaction with the given signature.
// Note that it might return ErrUnsupportedTxVersion (for tx that we don't support yet).
func GetTransaction(
	ctx context.Context,
	client interfaces.SolanaRPCClient,
	sig solana.Signature,
) (*rpc.GetTransactionResult, error) {
	txResult, err := client.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0,
	})

	switch {
	case err != nil && strings.Contains(err.Error(), errorCodeUnsupportedTransactionVersion):
		return nil, ErrUnsupportedTxVersion
	case err != nil:
		return nil, err
	default:
		return txResult, nil
	}
}

// HealthCheck returns the last block time
func HealthCheck(ctx context.Context, client interfaces.SolanaRPCClient) (time.Time, error) {
	// query latest slot
	slot, err := client.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get latest slot")
	}

	// query latest block time
	blockTime, err := client.GetBlockTime(ctx, slot)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get latest block time")
	}

	return blockTime.Time(), nil
}
