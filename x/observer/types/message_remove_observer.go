package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveObserver = "remove_observer"

var _ sdk.Msg = &MsgRemoveObserver{}

func NewMsgRemoveObserver(
	creator string,
	observerAddress string,
) *MsgRemoveObserver {
	return &MsgRemoveObserver{
		Creator:         creator,
		ObserverAddress: observerAddress,
	}
}

func (msg *MsgRemoveObserver) Route() string {
	return RouterKey
}

func (msg *MsgRemoveObserver) Type() string {
	return TypeMsgRemoveObserver
}

func (msg *MsgRemoveObserver) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveObserver) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveObserver) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.ObserverAddress)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid observer address (%s)", err)
	}
	return nil
}
