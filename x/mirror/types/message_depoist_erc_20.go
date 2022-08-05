package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDepoistERC20 = "depoist_erc_20"

var _ sdk.Msg = &MsgDepoistERC20{}

func NewMsgDepoistERC20(creator string, homeERC20ContractAddress string, recipientAddress string, amount string) *MsgDepoistERC20 {
	return &MsgDepoistERC20{
		Creator:                  creator,
		HomeERC20ContractAddress: homeERC20ContractAddress,
		RecipientAddress:         recipientAddress,
		Amount:                   amount,
	}
}

func (msg *MsgDepoistERC20) Route() string {
	return RouterKey
}

func (msg *MsgDepoistERC20) Type() string {
	return TypeMsgDepoistERC20
}

func (msg *MsgDepoistERC20) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDepoistERC20) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDepoistERC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
