package keeper

import (
	"context"
	"fmt"
	"strings"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MaxOutTxTrackerHashes is the maximum number of hashes that can be stored in the outbound transaction tracker
const MaxOutTxTrackerHashes = 2

// AddToOutTxTracker adds a new record to the outbound transaction tracker.
// only the admin policy account and the observer validators are authorized to broadcast this message without proof.
// If no pending cctx is found, the tracker is removed, if there is an existed tracker with the nonce & chainID.
func (k msgServer) AddToOutTxTracker(goCtx context.Context, msg *types.MsgAddToOutTxTracker) (*types.MsgAddToOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check the chain is supported
	chain := k.GetObserverKeeper().GetSupportedChainFromChainID(ctx, msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	// the cctx must exist
	cctx, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
		ChainID: msg.ChainId,
		Nonce:   msg.Nonce,
	})
	if err != nil {
		return nil, cosmoserrors.Wrap(types.ErrCannotFindCctx, err.Error())
	}
	if cctx == nil || cctx.CrossChainTx == nil {
		return nil, cosmoserrors.Wrapf(types.ErrCannotFindCctx, "no corresponding cctx found for chain %d, nonce %d", msg.ChainId, msg.Nonce)
	}

	// tracker submission is only allowed when the cctx is pending
	if !IsPending(cctx.CrossChainTx) {
		// garbage tracker (for any reason) is harmful to outTx observation and should be removed if it exists
		// it if does not exist, RemoveOutTxTracker is a no-op
		k.RemoveOutTxTracker(ctx, msg.ChainId, msg.Nonce)
		return &types.MsgAddToOutTxTrackerResponse{IsRemoved: true}, nil
	}

	isEmergencyGroup := false
	err = k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err == nil {
		isEmergencyGroup = true
	}
	isObserver := k.GetObserverKeeper().IsNonTombstonedObserver(ctx, msg.Creator)
	isProven := false

	// only emergency group and observer can submit tracker without proof
	// if the sender is not from the emergency group or observer, the outbound proof must be provided
	if !(isEmergencyGroup || isObserver) {
		if msg.Proof == nil {
			return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
		}
		// verify proof when it is provided
		if err := verifyProofAndOutTxBody(ctx, k, msg); err != nil {
			return nil, err
		}

		isProven = true
	}

	// fetch the tracker
	// if the tracker does not exist, initialize a new one
	tracker, found := k.GetOutTxTracker(ctx, msg.ChainId, msg.Nonce)
	hash := types.TxHashList{
		TxHash:   msg.TxHash,
		TxSigner: msg.Creator,
		Proved:   isProven,
	}
	if !found {
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    "",
			ChainId:  msg.ChainId,
			Nonce:    msg.Nonce,
			HashList: []*types.TxHashList{&hash},
		})
		return &types.MsgAddToOutTxTrackerResponse{}, nil
	}

	// check if the hash is already in the tracker
	for i, hash := range tracker.HashList {
		hash := hash
		if strings.EqualFold(hash.TxHash, msg.TxHash) {
			// if the hash is already in the tracker but we have a proof, mark it as proven and only keep this one in the list
			if isProven {
				tracker.HashList[i].Proved = true
				k.SetOutTxTracker(ctx, tracker)
			}
			return &types.MsgAddToOutTxTrackerResponse{}, nil
		}
	}

	// check if max hashes are reached
	if len(tracker.HashList) >= MaxOutTxTrackerHashes {
		return nil, types.ErrMaxTxOutTrackerHashesReached.Wrapf(
			"max hashes reached for chain %d, nonce %d, hash number: %d",
			msg.ChainId,
			msg.Nonce,
			len(tracker.HashList),
		)
	}

	// add the tracker to the list
	tracker.HashList = append(tracker.HashList, &hash)
	k.SetOutTxTracker(ctx, tracker)
	return &types.MsgAddToOutTxTrackerResponse{}, nil
}

// verifyProofAndOutTxBody verifies the proof and outbound tx body
// Precondition: the proof must be non-nil
func verifyProofAndOutTxBody(ctx sdk.Context, k msgServer, msg *types.MsgAddToOutTxTracker) error {
	txBytes, err := k.lightclientKeeper.VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
	if err != nil {
		return types.ErrProofVerificationFail.Wrapf(err.Error())
	}

	// get tss address
	var bitcoinChainID int64
	if chains.IsBitcoinChain(msg.ChainId) {
		bitcoinChainID = msg.ChainId
	}

	tss, err := k.GetObserverKeeper().GetTssAddress(ctx, &observertypes.QueryGetTssAddressRequest{
		BitcoinChainId: bitcoinChainID,
	})
	if err != nil {
		return observertypes.ErrTssNotFound.Wrapf(err.Error())
	}
	if tss == nil {
		return observertypes.ErrTssNotFound.Wrapf("tss address nil")
	}

	if err := types.VerifyOutTxBody(*msg, txBytes, *tss); err != nil {
		return types.ErrTxBodyVerificationFail.Wrapf(err.Error())
	}

	return nil
}
