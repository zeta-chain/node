package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// AddInboundTracker adds a new record to the inbound transaction tracker.
func (k msgServer) AddInboundTracker(
	goCtx context.Context,
	msg *types.MsgAddInboundTracker,
) (*types.MsgAddInboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, found := k.GetObserverKeeper().GetSupportedChainFromChainID(ctx, msg.ChainId); !found {
		return nil, observertypes.ErrSupportedChains
	}

	// only emergency group and observer can submit a tracker
	var (
		isAuthorizedPolicy = k.GetAuthorityKeeper().CheckAuthorization(ctx, msg) == nil
		isObserver         = k.GetObserverKeeper().CheckObserverCanVote(ctx, msg.Creator) == nil
	)

	if !isAuthorizedPolicy && !isObserver {
		return nil, errorsmod.Wrapf(authoritytypes.ErrUnauthorized, "Creator %s", msg.Creator)
	}

	// add the inTx tracker
	k.SetInboundTracker(ctx, types.InboundTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})

	return &types.MsgAddInboundTrackerResponse{}, nil
}
