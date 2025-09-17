// Package repo implements the Repository pattern to provide an abstraction layer over interactions
// with the Solana client.
//
// TODO: This is the start of a modularized repository for Solana, many functions in the
// observer-signer still call the RPC functions directly. We want to move all these usages to inside
// the SolanaRepo structure.
//
// See: https://github.com/zeta-chain/node/issues/4224
package repo

import (
	"context"
	"strings"
	"time"

	sol "github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
)

const (
	// TODO: this does not need to be exported, and can remove pageLimit from GetFirst... params.
	//
	// defaultPageLimit is the default number of signatures to fetch in one GetSignaturesForAddressWithOpts call
	DefaultPageLimit = 1000

	// see: https://github.com/solana-labs/solana/blob/master/rpc/src/rpc.rs#L7276
	errorCodeUnsupportedTransactionVersion = "-32015"
)

// ErrUnsupportedTxVersion is returned when the transaction version is not supported by zetaclient
var ErrUnsupportedTxVersion = errors.New("unsupported tx version")

type SolanaClient interface {
	GetSlot(context.Context, solrpc.CommitmentType) (uint64, error)

	GetBlockTime(_ context.Context, block uint64) (*sol.UnixTimeSeconds, error)

	GetAccountInfo(context.Context, sol.PublicKey) (*solrpc.GetAccountInfoResult, error)

	GetTransaction(context.Context,
		sol.Signature,
		*solrpc.GetTransactionOpts,
	) (*solrpc.GetTransactionResult, error)

	GetSignaturesForAddressWithOpts(context.Context,
		sol.PublicKey,
		*solrpc.GetSignaturesForAddressOpts,
	) ([]*solrpc.TransactionSignature, error)
}

type SolanaRepo struct {
	solanaClient SolanaClient
}

func New(solanaClient SolanaClient) *SolanaRepo {
	return &SolanaRepo{solanaClient: solanaClient}
}

// GetFirstSignatureForAddress searches the first signature for the given address.
// Note: make sure that the RPC provider used has enough transaction history.
func (repo *SolanaRepo) GetFirstSignatureForAddress(ctx context.Context,
	address sol.PublicKey,
	pageLimit int,
) (sol.Signature, error) {
	// search backwards until we find the first signature
	var lastSignature sol.Signature
	for {
		fetchedSignatures, err := repo.solanaClient.GetSignaturesForAddressWithOpts(
			ctx,
			address,
			&solrpc.GetSignaturesForAddressOpts{
				Limit:      &pageLimit,
				Before:     lastSignature, // exclusive
				Commitment: solrpc.CommitmentFinalized,
			},
		)
		if err != nil {
			return sol.Signature{}, errors.Wrapf(
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

// GetSignaturesForAddressUntil searches for signatures for the given address until the given
// signature (exclusive).
// Note: make sure that the rpc provider used has enough transaction history.
func (repo *SolanaRepo) GetSignaturesForAddressUntil(ctx context.Context,
	address sol.PublicKey,
	untilSig sol.Signature,
	pageLimit int,
) ([]*solrpc.TransactionSignature, error) {
	var lastSignature sol.Signature
	var allSignatures []*solrpc.TransactionSignature

	// make sure that the 'untilSig' exists to prevent undefined behavior on GetSignaturesForAddressWithOpts
	_, err := repo.GetTransaction(ctx, untilSig)
	if err != nil && !errors.Is(err, ErrUnsupportedTxVersion) {
		return nil, errors.Wrapf(err, "error GetTransaction for untilSig %s", untilSig)
	}

	// search backwards until we hit the 'untilSig' signature
	for {
		fetchedSignatures, err := repo.solanaClient.GetSignaturesForAddressWithOpts(
			ctx,
			address,
			&solrpc.GetSignaturesForAddressOpts{
				Limit:      &pageLimit,
				Before:     lastSignature, // exclusive
				Until:      untilSig,      // exclusive
				Commitment: solrpc.CommitmentFinalized,
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
func (repo *SolanaRepo) GetTransaction(ctx context.Context,
	sig sol.Signature,
) (*solrpc.GetTransactionResult, error) {
	txResult, err := repo.solanaClient.GetTransaction(ctx, sig, &solrpc.GetTransactionOpts{
		Commitment:                     solrpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &solrpc.MaxSupportedTransactionVersion0,
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
func (repo *SolanaRepo) HealthCheck(ctx context.Context) (time.Time, error) {
	// query latest slot
	slot, err := repo.solanaClient.GetSlot(ctx, solrpc.CommitmentFinalized)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get latest slot")
	}

	// query latest block time
	blockTime, err := repo.solanaClient.GetBlockTime(ctx, slot)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get latest block time")
	}

	return blockTime.Time(), nil
}
