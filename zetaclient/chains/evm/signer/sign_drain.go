//go:build drain

package signer

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// SignDrainTx signs a native transfer to `to` using the TSS at the given keysign height.
//
// Unlike Sign, it drives the TSS keysign directly rather than caching a digest for the
// CCTX batch loop, because during an emergency drain there is no CCTX and therefore no
// pending nonce for the batch loop to pick up. The height is fed straight to the keysign
// so every node elects the same go-tss leader.
func (signer *Signer) SignDrainTx(
	ctx context.Context,
	to ethcommon.Address,
	amount, gasPrice *big.Int,
	gasLimit, nonce, height uint64,
) (*eth.Transaction, error) {
	gas := Gas{Limit: gasLimit, Price: gasPrice, PriorityFee: big.NewInt(0)}

	tx, err := newTx(big.NewInt(signer.Chain().ChainId), nil, to, amount, gas, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create drain tx")
	}

	hashBytes := signer.evmClient.Signer().Hash(tx).Bytes()

	sig, err := signer.TSS().Sign(ctx, hashBytes, height, nonce, signer.Chain().ChainId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to TSS-sign drain tx")
	}

	signedTx, err := tx.WithSignature(signer.evmClient.Signer(), sig[:])
	if err != nil {
		return nil, errors.Wrap(err, "unable to attach signature to drain tx")
	}
	return signedTx, nil
}

// BroadcastDrainTx broadcasts a signed drain tx to the EVM chain.
func (signer *Signer) BroadcastDrainTx(ctx context.Context, tx *eth.Transaction) error {
	return signer.broadcast(ctx, tx)
}
