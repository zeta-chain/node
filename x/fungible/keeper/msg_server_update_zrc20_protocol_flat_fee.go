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

func (k msgServer) UpdateZRC20ProtocolFlatFee(goCtx context.Context, msg *types.MsgUpdateZRC20ProtocolFlatFee) (*types.MsgUpdateZRC20ProtocolFlatFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != types.AdminAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only admin can deploy fungible coin")
	}
	zrc20 := ethcommon.HexToAddress(msg.Zrc20Address)
	if zrc20 == (ethcommon.Address{}) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, msg.Zrc20Address)
	}
	fee := big.NewInt(0)
	_, success := fee.SetString(msg.ProtocolFlatFee, 10)
	if !success {
		return nil, sdkerrors.Wrap(types.ErrInvalidAmount, msg.ProtocolFlatFee)
	}
	tx, err := k.ZRC20UpdateProtocolFlatFee(ctx, zrc20, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrContractCall, err.Error())
	}
	if tx.Failed() {
		return nil, sdkerrors.Wrap(types.ErrContractCall, tx.VmError)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("UpdateZRC20ProtocolFlatFee",
			sdk.NewAttribute("zrc20", msg.Zrc20Address),
			sdk.NewAttribute("fee", fmt.Sprintf("%d", msg.ProtocolFlatFee)),
			sdk.NewAttribute("txHash", tx.Hash),
			sdk.NewAttribute("VmError", tx.VmError),
			sdk.NewAttribute("GasUsed", fmt.Sprintf("%d", tx.GasUsed)),
			sdk.NewAttribute("ret", fmt.Sprintf("%x", tx.Ret)),
			sdk.NewAttribute("logs", fmt.Sprintf("%v", tx.Logs)),
		),
	)

	return &types.MsgUpdateZRC20ProtocolFlatFeeResponse{}, nil
}
