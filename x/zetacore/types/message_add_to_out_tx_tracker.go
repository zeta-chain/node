package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddToOutTxTracker = "add_to_out_tx_tracker"

var _ sdk.Msg = &MsgAddToOutTxTracker{}

func NewMsgAddToOutTxTracker(creator string, chain string, nonce uint64, txHash string) *MsgAddToOutTxTracker {
	return &MsgAddToOutTxTracker{
		Creator: creator,
		Chain:   chain,
		Nonce:   nonce,
		TxHash:  txHash,
	}
}

func (msg *MsgAddToOutTxTracker) Route() string {
	return RouterKey
}

func (msg *MsgAddToOutTxTracker) Type() string {
	return TypeMsgAddToOutTxTracker
}

func (msg *MsgAddToOutTxTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddToOutTxTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddToOutTxTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
