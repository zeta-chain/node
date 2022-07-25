package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddToWatchList = "add_to_watch_list"

var _ sdk.Msg = &MsgAddToWatchList{}

func NewMsgAddToWatchList(creator string, chain string, nonce uint64, txHash string) *MsgAddToWatchList {
	return &MsgAddToWatchList{
		Creator: creator,
		Chain:   chain,
		Nonce:   nonce,
		TxHash:  txHash,
	}
}

func (msg *MsgAddToWatchList) Route() string {
	return RouterKey
}

func (msg *MsgAddToWatchList) Type() string {
	return TypeMsgAddToWatchList
}

func (msg *MsgAddToWatchList) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddToWatchList) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddToWatchList) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
