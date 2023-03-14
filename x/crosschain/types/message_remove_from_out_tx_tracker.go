package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveFromOutTxTracker = "RemoveFromTracker"

var _ sdk.Msg = &MsgRemoveFromOutTxTracker{}

func NewMsgRemoveFromOutTxTracker(creator string, chain int64, nonce uint64) *MsgRemoveFromOutTxTracker {
	return &MsgRemoveFromOutTxTracker{
		Creator: creator,
		ChainId: chain,
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
	if msg.ChainId < 0 {
		return sdkerrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	return nil
}
