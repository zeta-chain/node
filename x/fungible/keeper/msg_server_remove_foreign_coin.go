package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// RemoveForeignCoin removes a coin from the list of foreign coins in the
// module's state.
//
// Authorized: admin policy group 2.
func (k msgServer) RemoveForeignCoin(goCtx context.Context, msg *types.MsgRemoveForeignCoin) (*types.MsgRemoveForeignCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}
	index := msg.Name
	_, found := k.GetForeignCoins(ctx, index)
	if !found {
		return nil, cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "foreign coin not found")
	}
	k.RemoveForeignCoins(ctx, index)
	return &types.MsgRemoveForeignCoinResponse{}, nil
}
