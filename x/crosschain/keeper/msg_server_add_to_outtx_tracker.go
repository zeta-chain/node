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

// AddToOutTxTracker adds a new record to the outbound transaction tracker.
// only the admin policy account and the observer validators are authorized to broadcast this message without proof.
// If no pending cctx is found, the tracker is removed, if there is an existed tracker with the nonce & chainID.
//
// Authorized: admin policy group 1, observer.
func (k msgServer) AddToOutTxTracker(goCtx context.Context, msg *types.MsgAddToOutTxTracker) (*types.MsgAddToOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	// the cctx must exist
	cctx, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
		ChainID: msg.ChainId,
		Nonce:   msg.Nonce,
	})
	if err != nil {
		return nil, cosmoserrors.Wrap(err, "CcxtByNonce failed")
	}
	if cctx == nil || cctx.CrossChainTx == nil {
		return nil, cosmoserrors.Wrapf(types.ErrCannotFindCctx, "no corresponding cctx found for chain %d, nonce %d", msg.ChainId, msg.Nonce)
	}
	// tracker submission is only allowed when the cctx is pending
	if !IsPending(*cctx.CrossChainTx) {
		// garbage tracker (for any reason) is harmful to outTx observation and should be removed
		k.RemoveOutTxTracker(ctx, msg.ChainId, msg.Nonce)
		return &types.MsgAddToOutTxTrackerResponse{IsRemoved: true}, nil
	}

	if msg.Proof == nil { // without proof, only certain accounts can send this message
		isAdmin := k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency)
		isObserver := k.zetaObserverKeeper.IsNonTombstonedObserver(ctx, msg.Creator)

		// Sender needs to be either the admin policy account or an observer
		if !(isAdmin || isObserver) {
			return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
		}
	}

	isProven := false
	if msg.Proof != nil { // verify proof when it is provided
		txBytes, err := k.lightclientKeeper.VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
		if err != nil {
			return nil, types.ErrProofVerificationFail.Wrapf(err.Error())
		}

		// get tss address
		var bitcoinChainID int64
		if chains.IsBitcoinChain(msg.ChainId) {
			bitcoinChainID = msg.ChainId
		}

		tss, err := k.zetaObserverKeeper.GetTssAddress(ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: bitcoinChainID,
		})
		if err != nil || tss == nil {
			reason := "tss response is nil"
			if err != nil {
				reason = err.Error()
			}
			return nil, observertypes.ErrTssNotFound.Wrapf("tss address not found %s", reason)
		}

		err = types.VerifyOutTxBody(*msg, txBytes, *tss)
		if err != nil {
			return nil, types.ErrTxBodyVerificationFail.Wrapf(err.Error())
		}
		isProven = true
	}

	tracker, found := k.GetOutTxTracker(ctx, msg.ChainId, msg.Nonce)
	hash := types.TxHashList{
		TxHash:   msg.TxHash,
		TxSigner: msg.Creator,
	}
	if !found {
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    "",
			ChainId:  chain.ChainId,
			Nonce:    msg.Nonce,
			HashList: []*types.TxHashList{&hash},
		})
		ctx.Logger().Info(fmt.Sprintf("Add tracker %s: , Block Height : %d ", getOutTrackerIndex(chain.ChainId, msg.Nonce), ctx.BlockHeight()))
		return &types.MsgAddToOutTxTrackerResponse{}, nil
	}

	var isDup = false
	for _, hash := range tracker.HashList {
		if strings.EqualFold(hash.TxHash, msg.TxHash) {
			isDup = true
			if isProven {
				hash.Proved = true
				k.SetOutTxTracker(ctx, tracker)
				k.Logger(ctx).Info("Proof'd outbound transaction")
				return &types.MsgAddToOutTxTrackerResponse{}, nil
			}
			break
		}
	}
	if !isDup {
		if isProven {
			hash.Proved = true
			tracker.HashList = append([]*types.TxHashList{&hash}, tracker.HashList...)
			k.Logger(ctx).Info("Proof'd outbound transaction")
		} else if len(tracker.HashList) < 2 {
			tracker.HashList = append(tracker.HashList, &hash)
		}
		k.SetOutTxTracker(ctx, tracker)
	}
	return &types.MsgAddToOutTxTrackerResponse{}, nil
}
