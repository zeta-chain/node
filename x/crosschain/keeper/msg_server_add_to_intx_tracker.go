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

	// emergency or observer group can submit tracker without proof
	isEmergencyGroup := k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency)
	isObserver := k.GetObserverKeeper().IsNonTombstonedObserver(ctx, msg.Creator)

	if !(isEmergencyGroup || isObserver) {
		// if not directly authorized, check the proof, if not provided, return unauthorized
		if msg.Proof == nil {
			return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
		}

		// verify the proof and tx body
		if err := verifyProofAndInTxBody(ctx, k, msg); err != nil {
			return nil, err
		}
	}

	// add the inTx tracker
	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})

	return &types.MsgAddToInTxTrackerResponse{}, nil
}

// verifyProofAndInTxBody verifies the proof and inbound tx body
func verifyProofAndInTxBody(ctx sdk.Context, k msgServer, msg *types.MsgAddToInTxTracker) error {
	txBytes, err := k.GetLightclientKeeper().VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
	if err != nil {
		return types.ErrProofVerificationFail.Wrapf(err.Error())
	}

	// get chain params and tss addresses to verify the inTx body
	chainParams, found := k.GetObserverKeeper().GetChainParamsByChainID(ctx, msg.ChainId)
	if !found || chainParams == nil {
		return types.ErrUnsupportedChain.Wrapf("chain params not found for chain %d", msg.ChainId)
	}
	tss, err := k.GetObserverKeeper().GetTssAddress(ctx, &observertypes.QueryGetTssAddressRequest{
		BitcoinChainId: msg.ChainId,
	})
	if err != nil || tss == nil {
		reason := "tss response is nil"
		if err != nil {
			reason = err.Error()
		}
		return observertypes.ErrTssNotFound.Wrapf("tss address not found %s", reason)
	}

	if err := types.VerifyInTxBody(*msg, txBytes, *chainParams, *tss); err != nil {
		return types.ErrTxBodyVerificationFail.Wrapf(err.Error())
	}

	return nil
}
