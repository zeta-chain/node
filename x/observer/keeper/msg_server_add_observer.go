package keeper

import (
	"context"
	"math"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/node/pkg/crypto"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// AddObserver adds an observer address to the observer set
func (k msgServer) AddObserver(
	goCtx context.Context,
	msg *types.MsgAddObserver,
) (*types.MsgAddObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}
	pubkey, err := crypto.NewPubKey(msg.ZetaclientGranteePubkey)
	if err != nil {
		return &types.MsgAddObserverResponse{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}
	granteeAddress, err := crypto.GetAddressFromPubkeyString(msg.ZetaclientGranteePubkey)
	if err != nil {
		return &types.MsgAddObserverResponse{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}

	k.DisableInboundOnly(ctx)

	// AddNodeAccountOnly flag usage
	// True: adds observer into the Node Account list but returns without adding to the observer list
	// False: adds observer to the observer list, and not the node account list
	// Inbound is disabled in both cases and needs to be enabled manually using an admin TX
	if msg.AddNodeAccountOnly {
		pubkeySet := crypto.PubKeySet{Secp256k1: pubkey, Ed25519: ""}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator:       msg.ObserverAddress,
			GranteeAddress: granteeAddress.String(),
			GranteePubkey:  &pubkeySet,
			NodeStatus:     types.NodeStatus_Active,
		})
		k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
		return &types.MsgAddObserverResponse{}, nil
	}

	// Add observer to the observer set and update the observer count
	count, err := k.AddObserverToSet(ctx, msg.ObserverAddress)
	if err != nil {
		return &types.MsgAddObserverResponse{}, err
	}

	EmitEventAddObserver(
		ctx,
		count,
		msg.ObserverAddress,
		granteeAddress.String(),
		msg.ZetaclientGranteePubkey,
	)

	return &types.MsgAddObserverResponse{}, nil
}
