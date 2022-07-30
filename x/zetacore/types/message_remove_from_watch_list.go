package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveFromWatchList = "remove_from_watch_list"

var _ sdk.Msg = &MsgRemoveFromWatchList{}

func NewMsgRemoveFromWatchList(creator string, chain string, nonce uint64) *MsgRemoveFromWatchList {
	return &MsgRemoveFromWatchList{
		Creator: creator,
		Chain:   chain,
		Nonce:   nonce,
	}
}

func (msg *MsgRemoveFromWatchList) Route() string {
	return RouterKey
}

func (msg *MsgRemoveFromWatchList) Type() string {
	return TypeMsgRemoveFromWatchList
}

func (msg *MsgRemoveFromWatchList) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveFromWatchList) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveFromWatchList) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
