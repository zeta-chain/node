package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

func (k msgServer) UpdateAuthorizations(
	goCtx context.Context,
	msg *types.MsgUpdateAuthorizations,
) (*types.MsgUpdateAuthorizationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check called by governance
	if k.govAddr.String() != msg.Signer {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority, expected %s, got %s",
			k.govAddr.String(),
			msg.Signer,
		)
	}

	list := k.UpdateAuthorizationList(ctx, msg.AddAuthorizationList, msg.RemoveAuthorizationList)
	err := k.SetAuthorizationList(ctx, list)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
