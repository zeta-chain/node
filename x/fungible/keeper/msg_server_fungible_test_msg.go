package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
)

func (k msgServer) FungibleTestMsg(goCtx context.Context, msg *types.MsgFungibleTestMsg) (*types.MsgFungibleTestMsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	addr, err := k.DeployZRC4Contract(ctx)
	if err != nil {
		return nil, err
	}

	alice := ethcommon.HexToAddress("0x9B4547Dd93e93c526a11c5123Ab74D42aFb41B10")
	bal1 := k.BalanceOfZRC4(ctx, addr, alice)
	if bal1 == nil {
		return nil, sdkerrors.Wrap(types.ErrBlanceQuery, fmt.Sprintf("zrc4 balance of %s", alice.String()))
	}
	_, err = k.DepositZRC4(ctx, addr, alice, big.NewInt(1000))
	if err != nil {
		return nil, err
	}
	bal2 := k.BalanceOfZRC4(ctx, addr, alice)
	if bal2 == nil {
		return nil, sdkerrors.Wrap(types.ErrBlanceQuery, fmt.Sprintf("zrc4 balance of %s", alice.String()))
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
			sdk.NewAttribute("action", "FungibleTestMsg"),
			sdk.NewAttribute("contract", addr.String()),
			sdk.NewAttribute("balance1", bal1.String()),
			sdk.NewAttribute("balance2", bal2.String()),
		),
	)

	return &types.MsgFungibleTestMsgResponse{}, nil
}
