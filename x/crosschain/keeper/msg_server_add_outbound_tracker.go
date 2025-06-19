package keeper

import (
	"context"
	"strings"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// MaxOutboundTrackerHashes is the maximum number of hashes that can be stored in the outbound transaction tracker

// AddOutboundTracker adds a new record to the outbound transaction tracker.
// only the admin policy account and the observer validators are authorized to broadcast this message without proof.
// If no pending cctx is found, the tracker is removed, if there is an existed tracker with the nonce & chainID.
func (k msgServer) AddOutboundTracker(
	goCtx context.Context,
	msg *types.MsgAddOutboundTracker,
) (*types.MsgAddOutboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check the chain is supported
	if _, found := k.GetObserverKeeper().GetSupportedChainFromChainID(ctx, msg.ChainId); !found {
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
		return nil, cosmoserrors.Wrapf(
			types.ErrCannotFindCctx,
			"no corresponding cctx found for chain %d, nonce %d",
			msg.ChainId,
			msg.Nonce,
		)
	}
	// tracker submission is only allowed when the cctx is pending
	if !IsPending(cctx.CrossChainTx) {
		// garbage tracker (for any reason) is harmful to outTx observation and should be removed if it exists
		// it if does not exist, RemoveOutboundTracker is a no-op
		k.RemoveOutboundTrackerFromStore(ctx, msg.ChainId, msg.Nonce)
		return &types.MsgAddOutboundTrackerResponse{IsRemoved: true}, nil
	}

	// check if the msg signer is from the emergency group policy address.
	// or an observer
	var (
		isAuthorizedPolicy = k.GetAuthorityKeeper().CheckAuthorization(ctx, msg) == nil
		isObserver         = k.GetObserverKeeper().CheckObserverCanVote(ctx, msg.Creator) == nil
	)

	if !isAuthorizedPolicy && !isObserver {
		return nil, cosmoserrors.Wrapf(authoritytypes.ErrUnauthorized, "Creator %s", msg.Creator)
	}

	// set the outbound hash from the last tracker and save it in the store
	// this value is helpful for explorer or front-end application to find the outbound hash, it has no on-chain utility
	// the hash will be replaced with the actual hash, if different, when the outbound transaction is observed and voted
	cctx.CrossChainTx.GetCurrentOutboundParam().Hash = msg.TxHash
	k.SetCrossChainTx(ctx, *cctx.CrossChainTx)

	// fetch the tracker
	// if the tracker does not exist, initialize a new one
	tracker, found := k.GetOutboundTracker(ctx, msg.ChainId, msg.Nonce)
	hash := types.TxHash{TxHash: msg.TxHash, TxSigner: msg.Creator}
	if !found {
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			Index:    "",
			ChainId:  msg.ChainId,
			Nonce:    msg.Nonce,
			HashList: []*types.TxHash{&hash},
		})
		return &types.MsgAddOutboundTrackerResponse{}, nil
	}

	// check if the hash is already in the tracker
	for _, hash := range tracker.HashList {
		if strings.EqualFold(hash.TxHash, msg.TxHash) {
			return &types.MsgAddOutboundTrackerResponse{}, nil
		}
	}

	// check if max hashes are reached
	if tracker.MaxReached() {
		return nil, types.ErrMaxTxOutTrackerHashesReached.Wrapf(
			"max hashes reached for chain %d, nonce %d, hash number: %d",
			msg.ChainId,
			msg.Nonce,
			len(tracker.HashList),
		)
	}

	// add the tracker to the list
	tracker.HashList = append(tracker.HashList, &hash)
	k.SetOutboundTracker(ctx, tracker)
	return &types.MsgAddOutboundTrackerResponse{}, nil
}
