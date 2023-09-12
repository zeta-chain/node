package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateCrosschainFlags = "update_crosschain_flags"
)

var _ sdk.Msg = &MsgUpdateCrosschainFlags{}

func NewMsgUpdateCrosschainFlags(creator string, isInboundEnabled, isOutboundEnabled bool) *MsgUpdateCrosschainFlags {
	return &MsgUpdateCrosschainFlags{
		Creator:           creator,
		IsInboundEnabled:  isInboundEnabled,
		IsOutboundEnabled: isOutboundEnabled,
	}
}

func (msg *MsgUpdateCrosschainFlags) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCrosschainFlags) Type() string {
	return TypeMsgUpdateCrosschainFlags
}

func (msg *MsgUpdateCrosschainFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCrosschainFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCrosschainFlags) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
