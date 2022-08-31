package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveFromOutTxTracker = "remove_from_out_tx_tracker"

var _ sdk.Msg = &MsgRemoveFromOutTxTracker{}

func NewMsgRemoveFromOutTxTracker(creator string, chain string, nonce uint64) *MsgRemoveFromOutTxTracker {
	return &MsgRemoveFromOutTxTracker{
		Creator: creator,
		Chain:   chain,
		Nonce:   nonce,
	}
}

func (msg *MsgRemoveFromOutTxTracker) Route() string {
	return RouterKey
}

func (msg *MsgRemoveFromOutTxTracker) Type() string {
	return TypeMsgRemoveFromOutTxTracker
}

func (msg *MsgRemoveFromOutTxTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveFromOutTxTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveFromOutTxTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
