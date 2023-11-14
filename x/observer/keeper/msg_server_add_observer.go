package keeper

import (
	"context"
	"math"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) AddObserver(goCtx context.Context, msg *types.MsgAddObserver) (*types.MsgAddObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_group2) {
		return &types.MsgAddObserverResponse{}, types.ErrNotAuthorizedPolicy
	}
	pubkey, err := common.NewPubKey(msg.ZetaclientGranteePubkey)
	if err != nil {
		return &types.MsgAddObserverResponse{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}
	granteeAddress, err := common.GetAddressFromPubkeyString(msg.ZetaclientGranteePubkey)
	if err != nil {
		return &types.MsgAddObserverResponse{}, cosmoserrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
	}
	k.DisableInboundOnly(ctx)
	// AddNodeAccountOnly flag usage
	// True: adds observer into the Node Account list but returns without adding to the observer list
	// False: adds observer to the observer list, and not the node account list
	// Inbound is disabled in both cases and needs to be enabled manually using an admin TX
	if msg.AddNodeAccountOnly {
		pubkeySet := common.PubKeySet{Secp256k1: pubkey, Ed25519: ""}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator:       msg.ObserverAddress,
			GranteeAddress: granteeAddress.String(),
			GranteePubkey:  &pubkeySet,
			NodeStatus:     types.NodeStatus_Active,
		})
		k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
		return &types.MsgAddObserverResponse{}, nil
	}

	observerMappers := k.GetAllObserverMappers(ctx)
	totalObserverCountCurrentBlock := uint64(0)
	for _, mapper := range observerMappers {
		mapper.ObserverList = append(mapper.ObserverList, msg.ObserverAddress)
		totalObserverCountCurrentBlock += uint64(len(mapper.ObserverList))
		k.SetObserverMapper(ctx, mapper)
	}
	k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: totalObserverCountCurrentBlock})
	EmitEventAddObserver(ctx, totalObserverCountCurrentBlock, msg.ObserverAddress, granteeAddress.String(), msg.ZetaclientGranteePubkey)
	return &types.MsgAddObserverResponse{}, nil
}
