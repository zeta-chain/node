package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateZRC20WithdrawFee = "update_zrc20_withdraw_fee"

var _ sdk.Msg = &MsgUpdateZRC20WithdrawFee{}

func NewMsgUpdateZRC20WithdrawFee(creator string, zrc20 string, newFee sdk.Uint) *MsgUpdateZRC20WithdrawFee {
	return &MsgUpdateZRC20WithdrawFee{
		Creator:        creator,
		Zrc20Address:   zrc20,
		NewWithdrawFee: newFee,
	}
}

func (msg *MsgUpdateZRC20WithdrawFee) Route() string {
	return RouterKey
}

func (msg *MsgUpdateZRC20WithdrawFee) Type() string {
	return TypeMsgUpdateSystemContract
}

func (msg *MsgUpdateZRC20WithdrawFee) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateZRC20WithdrawFee) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateZRC20WithdrawFee) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	if !ethcommon.IsHexAddress(msg.Zrc20Address) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s)", msg.Zrc20Address)
	}
	if msg.NewWithdrawFee.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid withdraw fee (%s)", msg.NewWithdrawFee)
	}

	return nil
}
