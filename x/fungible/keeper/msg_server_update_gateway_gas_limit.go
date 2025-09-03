package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// UpdateGatewayGasLimit updates the gateway gas limit used by the ZetaChain protocol
func (k msgServer) UpdateGatewayGasLimit(
	goCtx context.Context,
	msg *types.MsgUpdateGatewayGasLimit,
) (*types.MsgUpdateGatewayGasLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// Get the current system contract to preserve other fields
	var protocolContracts types.SystemContract
	protocolContracts, found := k.GetSystemContract(ctx)
	if !found {
		protocolContracts = types.SystemContract{}
	}
	oldGasLimit := protocolContracts.GatewayGasLimit
	k.SetGatewayGasLimit(ctx, msg.NewGasLimit)

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventGatewayGasLimitUpdated{
			MsgTypeUrl:  sdk.MsgTypeURL(&types.MsgUpdateGatewayGasLimit{}),
			NewGasLimit: msg.NewGasLimit.String(),
			OldGasLimit: oldGasLimit.String(),
			Signer:      msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgUpdateGatewayGasLimitResponse{}, nil
}
