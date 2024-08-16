package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgMigrateERC20CustodyFunds = "MigrateERC20CustodyFunds"

var _ sdk.Msg = &MsgMigrateERC20CustodyFunds{}

func NewMsgMigrateERC20CustodyFunds(
	creator string,
	chainID int64,
	newCustodyAddress string,
	erc20Address string,
	amount sdkmath.Uint,
) *MsgMigrateERC20CustodyFunds {
	return &MsgMigrateERC20CustodyFunds{
		Creator:           creator,
		ChainId:           chainID,
		NewCustodyAddress: newCustodyAddress,
		Erc20Address:      erc20Address,
		Amount:            amount,
	}
}

func (msg *MsgMigrateERC20CustodyFunds) Route() string {
	return RouterKey
}

func (msg *MsgMigrateERC20CustodyFunds) Type() string {
	return TypeMsgMigrateERC20CustodyFunds
}

func (msg *MsgMigrateERC20CustodyFunds) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMigrateERC20CustodyFunds) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateERC20CustodyFunds) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	switch {
	case !ethcommon.IsHexAddress(msg.NewCustodyAddress):
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new custody address")
	case !ethcommon.IsHexAddress(msg.Erc20Address):
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid erc20 address")
	case msg.Amount.IsZero():
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount cannot be zero")
	}

	return nil
}
