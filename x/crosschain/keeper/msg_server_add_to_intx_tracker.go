package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AddToInTxTracker adds a new record to the inbound transaction tracker.
func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	isAdmin := k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency)
	isObserver := k.zetaObserverKeeper.IsNonTombstonedObserver(ctx, msg.Creator)

	isProven := false
	if !(isAdmin || isObserver) && msg.Proof != nil {
		txBytes, err := k.lightclientKeeper.VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
		if err != nil {
			return nil, types.ErrProofVerificationFail.Wrapf(err.Error())
		}

		// get chain params and tss addresses to verify the inTx body
		chainParams, found := k.zetaObserverKeeper.GetChainParamsByChainID(ctx, msg.ChainId)
		if !found || chainParams == nil {
			return nil, types.ErrUnsupportedChain.Wrapf("chain params not found for chain %d", msg.ChainId)
		}
		tss, err := k.zetaObserverKeeper.GetTssAddress(ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: msg.ChainId,
		})
		if err != nil || tss == nil {
			reason := "tss response is nil"
			if err != nil {
				reason = err.Error()
			}
			return nil, observertypes.ErrTssNotFound.Wrapf("tss address not found %s", reason)
		}

		if err := types.VerifyInTxBody(*msg, txBytes, *chainParams, *tss); err != nil {
			return nil, types.ErrTxBodyVerificationFail.Wrapf(err.Error())
		}

		isProven = true
	}

	// Sender needs to be either the admin policy account or an observer
	if !(isAdmin || isObserver || isProven) {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})
	return &types.MsgAddToInTxTrackerResponse{}, nil
}
