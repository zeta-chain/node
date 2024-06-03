package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddAuthorization = "AddAuthorization"

var _ sdk.Msg = &MsgAddAuthorization{}

func NewMsgAddAuthorization(creator string, msgURL string, authorizedPolicy PolicyType) *MsgAddAuthorization {
	return &MsgAddAuthorization{
		Creator:          creator,
		MsgUrl:           msgURL,
		AuthorizedPolicy: authorizedPolicy,
	}
}

func (msg *MsgAddAuthorization) Route() string {
	return RouterKey
}

func (msg *MsgAddAuthorization) Type() string {
	return TypeMsgAddAuthorization
}

func (msg *MsgAddAuthorization) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddAuthorization) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddAuthorization) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// the authorized policy must be valid
	if err := msg.AuthorizedPolicy.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authorized policy: %s", err.Error())
	}

	// the message URL must be valid
	if err := ValidateMsgURL(msg.MsgUrl); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg url: %s", err.Error())
	}

	return nil
}

func ValidateMsgURL(url string) error {
	if len(url) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "message URL cannot be empty")
	}
	return nil
}
