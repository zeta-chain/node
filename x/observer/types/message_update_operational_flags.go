package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateOperationalFlags = "update_operational_flags"
)

var _ sdk.Msg = &MsgUpdateOperationalFlags{}

func NewMsgUpdateOperationalFlags(creator string, flags OperationalFlags) *MsgUpdateOperationalFlags {
	return &MsgUpdateOperationalFlags{
		Creator:          creator,
		OperationalFlags: flags,
	}
}

func (msg *MsgUpdateOperationalFlags) Route() string {
	return RouterKey
}

func (msg *MsgUpdateOperationalFlags) Type() string {
	return TypeMsgUpdateOperationalFlags
}

func (msg *MsgUpdateOperationalFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateOperationalFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateOperationalFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := msg.OperationalFlags.Validate(); err != nil {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}
