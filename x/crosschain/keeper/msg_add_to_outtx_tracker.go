package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AddToOutTxTracker adds a new record to the outbound transaction tracker.
// only the admin policy account and the observer validators are authorized to broadcast this message.
func (k msgServer) AddToOutTxTracker(goCtx context.Context, msg *types.MsgAddToOutTxTracker) (*types.MsgAddToOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}
	adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group1)
	isAdmin := msg.Creator == adminPolicyAccount

	isObserver, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)
	if err != nil {
		ctx.Logger().Error("Error while checking if the account is an observer", err)
	}

	isProven := false
	if !(isAdmin || isObserver) && msg.Proof != nil {
		txx, err := k.VerifyProof(ctx, msg.Proof, msg.BlockHash, msg.TxIndex, msg.ChainId)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
		}
		err = k.VerifyOutTxTrackerProof(ctx, txx, msg.Nonce)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
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
		} else {
			tracker.HashList = append(tracker.HashList, &hash)
		}
		k.SetOutTxTracker(ctx, tracker)
	}
	return &types.MsgAddToOutTxTrackerResponse{}, nil
}

func (k Keeper) VerifyOutTxTrackerProof(ctx sdk.Context, txx ethtypes.Transaction, nonce uint64) error {
	signer := ethtypes.NewLondonSigner(txx.ChainId())
	sender, err := ethtypes.Sender(signer, &txx)
	if err != nil {
		return err
	}
	res, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return err
	}
	tssAddr := eth.HexToAddress(res.Eth)
	if tssAddr == (eth.Address{}) {
		return fmt.Errorf("tss address not found")
	}
	if sender != tssAddr {
		return fmt.Errorf("sender is not tss address")
	}
	if txx.Nonce() != nonce {
		return fmt.Errorf("nonce mismatch")
	}
	return nil
}
