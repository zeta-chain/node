package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// UpdateZRC20Name updates the name and/or the symbol of a zrc20 token
func (k msgServer) UpdateZRC20Name(
	goCtx context.Context,
	msg *types.MsgUpdateZRC20Name,
) (*types.MsgUpdateZRC20NameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check signer permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check the zrc20 is valid
	zrc20Addr := ethcommon.HexToAddress(msg.Zrc20Address)
	if zrc20Addr == (ethcommon.Address{}) {
		return nil, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid zrc20 contract address (%s)",
			msg.Zrc20Address,
		)
	}

	// check the zrc20 exists
	fc, found := k.GetForeignCoins(ctx, msg.Zrc20Address)
	if !found {
		return nil, cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"no foreign coin match requested zrc20 address (%s)",
			msg.Zrc20Address,
		)
	}

	// call the contract methods and update the object
	if msg.Name != "" {
		if err := k.ZRC20SetName(ctx, zrc20Addr, msg.Name); err != nil {
			return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to update zrc20 name (%s)", err.Error())
		}
		fc.Name = msg.Name
	}

	if msg.Symbol != "" {
		if err = k.ZRC20SetSymbol(ctx, zrc20Addr, msg.Symbol); err != nil {
			return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to update zrc20 symbol (%s)", err.Error())
		}
		fc.Symbol = msg.Symbol
	}

	// save the object
	k.SetForeignCoins(ctx, fc)

	return &types.MsgUpdateZRC20NameResponse{}, nil
}
