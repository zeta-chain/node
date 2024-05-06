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

// AddInboundTracker adds a new record to the inbound transaction tracker.
func (k msgServer) AddInboundTracker(goCtx context.Context, msg *types.MsgAddInboundTracker) (*types.MsgAddInboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.GetObserverKeeper().GetSupportedChainFromChainID(ctx, msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	// emergency or observer group can submit tracker without proof
	isEmergencyGroup := k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency)
	isObserver := k.GetObserverKeeper().IsNonTombstonedObserver(ctx, msg.Creator)

	// only emergency group and observer can submit tracker without proof
	// if the sender is not from the emergency group or observer, the inbound proof must be provided
	if !(isEmergencyGroup || isObserver) {
		if msg.Proof == nil {
			return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
		}

		// verify the proof and tx body
		if err := verifyProofAndInboundBody(ctx, k, msg); err != nil {
			return nil, err
		}
	}

	// add the inTx tracker
	k.SetInboundTracker(ctx, types.InboundTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})

	return &types.MsgAddInboundTrackerResponse{}, nil
}

// verifyProofAndInboundBody verifies the proof and inbound tx body
func verifyProofAndInboundBody(ctx sdk.Context, k msgServer, msg *types.MsgAddInboundTracker) error {
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
	if err != nil {
		return observertypes.ErrTssNotFound.Wrapf(err.Error())
	}
	if tss == nil {
		return observertypes.ErrTssNotFound.Wrapf("tss address nil")
	}

	if err := types.VerifyInboundBody(*msg, txBytes, *chainParams, *tss); err != nil {
		return types.ErrTxBodyVerificationFail.Wrapf(err.Error())
	}

	return nil
}
