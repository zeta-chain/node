package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateZRC20ProtocolFlatFee = "update_zrc_20_protocol_flat_fee"

var _ sdk.Msg = &MsgUpdateZRC20ProtocolFlatFee{}

func NewMsgUpdateZRC20ProtocolFlatFee(creator string, ZRC20 string, fee string) *MsgUpdateZRC20ProtocolFlatFee {
	return &MsgUpdateZRC20ProtocolFlatFee{
		Creator:         creator,
		Zrc20Address:    ZRC20,
		ProtocolFlatFee: fee,
	}
}

func (msg *MsgUpdateZRC20ProtocolFlatFee) Route() string {
	return RouterKey
}

func (msg *MsgUpdateZRC20ProtocolFlatFee) Type() string {
	return TypeMsgUpdateZRC20ProtocolFlatFee
}

func (msg *MsgUpdateZRC20ProtocolFlatFee) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateZRC20ProtocolFlatFee) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateZRC20ProtocolFlatFee) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
