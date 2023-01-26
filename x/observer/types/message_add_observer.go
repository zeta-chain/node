package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddObserver = "add_observer"

var _ sdk.Msg = &MsgAddObserver{}

func NewMsgAddObserver(creator string, chainID int64, observationType ObservationType) *MsgAddObserver {
	return &MsgAddObserver{
		Creator:         creator,
		ChainId:         chainID,
		ObservationType: observationType,
	}
}

func (msg *MsgAddObserver) Route() string {
	return RouterKey
}

func (msg *MsgAddObserver) Type() string {
	return TypeMsgAddObserver
}

func (msg *MsgAddObserver) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddObserver) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddObserver) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
