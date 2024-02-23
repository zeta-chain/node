package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateKeygen = "update_keygen"

var _ sdk.Msg = &MsgUpdateKeygen{}

func NewMsgUpdateKeygen(creator string, block int64) *MsgUpdateKeygen {
	return &MsgUpdateKeygen{
		Creator: creator,
		Block:   block,
	}
}

func (msg *MsgUpdateKeygen) Route() string {
	return RouterKey
}

func (msg *MsgUpdateKeygen) Type() string {
	return TypeMsgUpdateKeygen
}

func (msg *MsgUpdateKeygen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateKeygen) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateKeygen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
