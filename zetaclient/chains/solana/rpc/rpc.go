package rpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

const (
	// defaultPageLimit is the default number of signatures to fetch in one GetSignaturesForAddressWithOpts call
	DefaultPageLimit = 1000
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
