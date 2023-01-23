package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddTokenEmission = "add_token_emission"

var _ sdk.Msg = &MsgAddTokenEmission{}

func NewMsgAddTokenEmission(creator string, category EmissionCategory, amount sdk.Dec) *MsgAddTokenEmission {
	return &MsgAddTokenEmission{
		Creator:  creator,
		Category: category,
		Amount:   amount,
	}
}

func (msg *MsgAddTokenEmission) Route() string {
	return RouterKey
}

func (msg *MsgAddTokenEmission) Type() string {
	return TypeMsgAddTokenEmission
}

func (msg *MsgAddTokenEmission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddTokenEmission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddTokenEmission) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
