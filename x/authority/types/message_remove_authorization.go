package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeRemoveAuthorization = "RemoveAuthorization"

var _ sdk.Msg = &MsgRemoveAuthorization{}

func NewMsgRemoveAuthorization(creator string, msgURL string) *MsgRemoveAuthorization {
	return &MsgRemoveAuthorization{
		Creator: creator,
		MsgUrl:  msgURL,
	}
}

func (msg *MsgRemoveAuthorization) Route() string {
	return RouterKey
}

func (msg *MsgRemoveAuthorization) Type() string {
	return TypeRemoveAuthorization
}

func (msg *MsgRemoveAuthorization) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveAuthorization) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveAuthorization) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateMsgURL(msg.MsgUrl); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg url: %s", err.Error())
	}

	return nil
}
