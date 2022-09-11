package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFungibleTestMsg = "fungible_test_msg"

var _ sdk.Msg = &MsgFungibleTestMsg{}

func NewMsgFungibleTestMsg(creator string) *MsgFungibleTestMsg {
	return &MsgFungibleTestMsg{
		Creator: creator,
	}
}

func (msg *MsgFungibleTestMsg) Route() string {
	return RouterKey
}

func (msg *MsgFungibleTestMsg) Type() string {
	return TypeMsgFungibleTestMsg
}

func (msg *MsgFungibleTestMsg) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFungibleTestMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFungibleTestMsg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
