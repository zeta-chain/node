package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDisableFastConfirmation = "disable_fast_confirmation"

var _ sdk.Msg = &MsgDisableFastConfirmation{}

func NewMsgDisableFastConfirmation(creator string, chainID int64) *MsgDisableFastConfirmation {
	return &MsgDisableFastConfirmation{
		Creator: creator,
		ChainId: chainID,
	}
}

func (msg *MsgDisableFastConfirmation) Route() string {
	return RouterKey
}

func (msg *MsgDisableFastConfirmation) Type() string {
	return TypeMsgDisableFastConfirmation
}

func (msg *MsgDisableFastConfirmation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableFastConfirmation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableFastConfirmation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
