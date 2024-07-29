package keeper

import (
	"context"
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// UpdateGatewayContract updates the zevm gateway contract used by the ZetaChain protocol to read inbounds and process outbounds
func (k msgServer) UpdateGatewayContract(
	goCtx context.Context,
	msg *types.MsgUpdateGatewayContract,
) (*types.MsgUpdateGatewayContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// The SystemContract state variable tracks the contract addresses used by the protocol
	// This variable is planned to be renamed ProtocolContracts in the future:
	// https://github.com/zeta-chain/node/issues/2576
	var protocolContracts types.SystemContract
	protocolContracts, found := k.GetSystemContract(ctx)
	if !found {
		// protocolContracts has never been set before, set an empty structure
		protocolContracts = types.SystemContract{}
	}
	oldGateway := protocolContracts.Gateway

	// update address and save
	protocolContracts.Gateway = msg.NewGatewayContractAddress
	k.SetSystemContract(ctx, protocolContracts)

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventGatewayContractUpdated{
			MsgTypeUrl:         sdk.MsgTypeURL(&types.MsgUpdateGatewayContract{}),
			NewContractAddress: msg.NewGatewayContractAddress,
			OldContractAddress: oldGateway,
			Signer:             msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgUpdateGatewayContractResponse{}, nil
}
