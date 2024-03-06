package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdatePolicies = "UpdatePolicies"

var _ sdk.Msg = &MsgUpdatePolicies{}

func NewMsgUpdatePolicies(signer string, policies Policies) *MsgUpdatePolicies {
	return &MsgUpdatePolicies{
		Signer:   signer,
		Policies: policies,
	}
}

func (msg *MsgUpdatePolicies) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePolicies) Type() string {
	return TypeMsgUpdatePolicies
}

func (msg *MsgUpdatePolicies) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePolicies) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePolicies) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer address (%s)", err)
	}

	if err := msg.Policies.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid policies (%s)", err)
	}

	return nil
}
