package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateObserver = "update_observer"

var _ sdk.Msg = &MsgUpdateObserver{}

func NewMsgUpdateObserver(creator string, oldObserverAddress string, newObserverAddress string, updateReason ObserverUpdateReason) *MsgUpdateObserver {
	return &MsgUpdateObserver{
		Creator:            creator,
		OldObserverAddress: oldObserverAddress,
		NewObserverAddress: newObserverAddress,
		UpdateReason:       updateReason,
	}
}

func (msg *MsgUpdateObserver) Route() string {
	return RouterKey
}

func (msg *MsgUpdateObserver) Type() string {
	return TypeMsgUpdateObserver
}

func (msg *MsgUpdateObserver) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateObserver) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateObserver) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.OldObserverAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid old observer address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.NewObserverAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new observer address (%s)", err)
	}
	if msg.UpdateReason != ObserverUpdateReason_Tombstoned && msg.UpdateReason != ObserverUpdateReason_AdminUpdate {
		return errorsmod.Wrapf(ErrUpdateObserver, "invalid update reason (%s)", msg.UpdateReason)
	}
	if msg.UpdateReason == ObserverUpdateReason_Tombstoned && msg.OldObserverAddress != msg.Creator {
		return errorsmod.Wrapf(ErrUpdateObserver, "invalid old observer address (%s)", msg.OldObserverAddress)
	}
	return nil
}
