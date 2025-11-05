// Package repo implements the Repository pattern to provide an abstraction layer over interactions
// with the Solana client.
package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	alt "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	zetamath "github.com/zeta-chain/node/pkg/math"
)

const (
	// see: https://github.com/solana-labs/solana/blob/master/rpc/src/rpc.rs#L7276
	errorCodeUnsupportedTransactionVersion = "-32015"

	// broadcastOutboundCommitment is the commitment level for broadcasting solana outbound.
	// Commitment "finalized" eliminate all risk but the tradeoff is pretty severe and effectively
	// reduces the expiration of tx by about 13 seconds. The "confirmed" level has very low risk of
	// belonging to a dropped fork.
	// see: https://solana.com/developers/guides/advanced/confirmation#use-an-appropriate-preflight-commitment-level
	broadcastOutboundCommitment = rpc.CommitmentConfirmed
)

type SolanaRepo struct {
	client SolanaClient

	// gatewayID is the program ID of the gateway program on Solana chain.
	gatewayID solana.PublicKey
}

func New(client SolanaClient, gatewayID solana.PublicKey) *SolanaRepo {
	return &SolanaRepo{client, gatewayID}
}

// HealthCheck checks the health of the RPC client by querying for the latest block time.
func (repo *SolanaRepo) HealthCheck(ctx context.Context) (*time.Time, error) {
	// Get last finalized slot.
	slot, err := repo.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	// Get latest block time.
	blockTime, err := repo.client.GetBlockTime(ctx, slot)
	if err != nil {
		return nil, newClientError(ErrClientGetBlockTime, err)
	}

	time := blockTime.Time()
	return &time, nil
}

// TODO
func (repo *SolanaRepo) GetLatestBlockHash(ctx context.Context) (blockhash solana.Hash, _ error) {
	result, err := repo.client.GetLatestBlockhash(ctx, broadcastOutboundCommitment)
	if err != nil {
		return blockhash, newClientError(ErrClientGetLatestBlockhash, err)
	}
	return result.Value.Blockhash, nil
}

// GetSlot returns the the most recent slot of the given commitment type.
func (repo *SolanaRepo) GetSlot(ctx context.Context,
	commitment rpc.CommitmentType,
) (uint64, error) {
	slot, err := repo.client.GetSlot(ctx, commitment)
	if err != nil {
		errSlot := fmt.Errorf("%w - %s", ErrClientGetSlot, commitment)
		return 0, newClientError(errSlot, err)
	}

	return slot, nil
}

// TODO
func (repo *SolanaRepo) GetAccountInfo(ctx context.Context,
	account solana.PublicKey,
	commitment rpc.CommitmentType,
) (*rpc.GetAccountInfoResult, error) {
	opts := rpc.GetAccountInfoOpts{Commitment: commitment}
	result, err := repo.client.GetAccountInfoWithOpts(ctx, account, &opts)
	if err != nil {
		return nil, newClientError(ErrClientGetAccountInfo, err)
	}

	return result, nil
}

// TODO
func (repo *SolanaRepo) GetBalance(ctx context.Context,
	account solana.PublicKey,
	commitment rpc.CommitmentType,
) (uint64, error) {
	result, err := repo.client.GetBalance(ctx, account, commitment)
	if err != nil {
		return 0, newClientError(ErrClientGetBalance, err)
	}

	return result.Value, nil
}

// GetTransaction returns a transaction with the given signature.
// It might return ErrUnsupportedTxVersion for transactions that we do not support yet.
func (repo *SolanaRepo) GetTransaction(ctx context.Context,
	sig solana.Signature,
	commitment rpc.CommitmentType,
) (*rpc.GetTransactionResult, error) {
	opts := &rpc.GetTransactionOpts{
		Commitment:                     commitment,
		MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion0,
	}

	result, err := repo.client.GetTransaction(ctx, sig, opts)
	if err != nil {
		if strings.Contains(err.Error(), errorCodeUnsupportedTransactionVersion) {
			err = fmt.Errorf("%w: %w", ErrUnsupportedTxVersion, err)
		} else {
			err = newClientError(ErrClientGetTransaction, err)
		}
		return nil, fmt.Errorf("%w (signature: %s)", err, sig.String())
	}

	return result, nil
}

// GetPriorityFee queries recent priority fees and returns their median (in micro lamports).
func (repo *SolanaRepo) GetPriorityFee(ctx context.Context) (uint64, error) {
	// Get recent priority fees.
	fees, err := repo.client.GetRecentPrioritizationFees(ctx, nil)
	if err != nil {
		return 0, newClientError(ErrClientGetRecentPrioritizationFees, err)
	}

	// Compute median priority fee (excludes zeroes).
	positiveFees := make([]uint64, 0, len(fees))
	for _, fee := range fees {
		if fee.PrioritizationFee > 0 {
			positiveFees = append(positiveFees, fee.PrioritizationFee)
		}
	}
	priorityFee := zetamath.SliceMedianValue(positiveFees, true)

	return priorityFee, nil
}

// GetFirstSignature returns the first signature for the gateway address.
func (repo *SolanaRepo) GetFirstSignature(ctx context.Context) (solana.Signature, error) {
	opts := rpc.GetSignaturesForAddressOpts{Commitment: rpc.CommitmentFinalized}
	_, err := repo.getSignatures(ctx, &opts, false)
	if err != nil {
		return solana.Signature{}, err
	}

	sig := opts.Before

	if sig.IsZero() {
		return sig, ErrFoundNoSignatures
	}

	return sig, nil
}

// GetSignaturesSince returns the signatures finalized since the given signature (exclusive).
func (repo *SolanaRepo) GetSignaturesSince(ctx context.Context,
	sig solana.Signature,
) ([]*rpc.TransactionSignature, error) {
	// Make sure that limit signature exists to prevent undefined behavior.
	_, err := repo.GetTransaction(ctx, sig, rpc.CommitmentFinalized)
	if err != nil && !errors.Is(err, ErrUnsupportedTxVersion) {
		return nil, err
	}

	opts := rpc.GetSignaturesForAddressOpts{
		Until:      sig, // exclusive
		Commitment: rpc.CommitmentFinalized,
	}

	return repo.getSignatures(ctx, &opts, true)
}

// TODO
func (repo *SolanaRepo) SendTransaction(ctx context.Context,
	tx *solana.Transaction,
) (solana.Signature, error) {
	opts := rpc.TransactionOpts{PreflightCommitment: broadcastOutboundCommitment}

	sig, err := repo.client.SendTransactionWithOpts(ctx, tx, opts)
	if err != nil {
		return sig, newClientError(ErrClientSendTransaction, err)
	}

	return sig, nil
}

func (repo *SolanaRepo) GetAddressLookupTableState(ctx context.Context,
	address solana.PublicKey,
) (*alt.AddressLookupTableState, error) {
	client, ok := repo.client.(*rpc.Client)
	if !ok {
		return nil, errors.New("solana AddressLookupTable requires *rpc.Client")
	}

	opts := rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentProcessed}
	altState, err := alt.GetAddressLookupTableStateWithOpts(ctx, client, address, &opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetAddressLookupTableState, err)
	}

	return altState, nil
}

// ------------------------------------------------------------------------------------------------

// getSignatures queries the last signatures given the parameters in opts.
// It returns the signatures if and only if store is set to true.
// It mutates opts, so opts.Before holds the last observed signature.
//
// NOTE: make sure that the RPC provider has enough transaction history.
func (repo SolanaRepo) getSignatures(ctx context.Context,
	opts *rpc.GetSignaturesForAddressOpts,
	store bool,
) ([]*rpc.TransactionSignature, error) {
	var all []*rpc.TransactionSignature

	// Search backwards until we hit the given signature.
	for {
		sigs, err := repo.client.GetSignaturesForAddressWithOpts(ctx, repo.gatewayID, opts)
		if err != nil {
			return nil, newClientError(ErrClientGetSignaturesForAddress, err)
		}

		// Stop if there are no more signatures.
		if len(sigs) == 0 {
			break
		}

		opts.Before = sigs[len(sigs)-1].Signature
		if store {
			all = append(all, sigs...)
		}
	}

	return all, nil
}
