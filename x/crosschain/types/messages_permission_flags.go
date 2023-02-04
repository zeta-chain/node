package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePermissionFlags = "create_permission_flags"
	TypeMsgUpdatePermissionFlags = "update_permission_flags"
	TypeMsgDeletePermissionFlags = "delete_permission_flags"
)

var _ sdk.Msg = &MsgUpdatePermissionFlags{}

func NewMsgUpdatePermissionFlags(creator string, isInboundEnabled bool) *MsgUpdatePermissionFlags {
	return &MsgUpdatePermissionFlags{
		Creator:          creator,
		IsInboundEnabled: isInboundEnabled,
	}
}

func (msg *MsgUpdatePermissionFlags) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePermissionFlags) Type() string {
	return TypeMsgUpdatePermissionFlags
}

func (msg *MsgUpdatePermissionFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePermissionFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePermissionFlags) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
