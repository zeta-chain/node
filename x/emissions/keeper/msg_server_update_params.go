package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

// UpdateParams defines a governance operation for updating the x/emissions module parameters.
// The authority is hard-coded to the x/gov module account.
func (k msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	if msg.Authority != k.authority {
		return nil, errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetParams, err.Error())
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
