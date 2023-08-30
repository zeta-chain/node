package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

const TypeMsgAddObserver = "add_observer"

var _ sdk.Msg = &MsgAddObserver{}

func NewMsgAddObserver(creator string, observerAdresss string, zetaclientGranteePubKey string, addNodeAccountOnly bool) *MsgAddObserver {
	return &MsgAddObserver{
		Creator:                 creator,
		ObserverAddress:         observerAdresss,
		ZetaclientGranteePubkey: zetaclientGranteePubKey,
		AddNodeAccountOnly:      addNodeAccountOnly,
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
	_, err = sdk.AccAddressFromBech32(msg.ObserverAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid observer address (%s)", err)
	}
	_, err = common.NewPubKey(msg.ZetaclientGranteePubkey)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid zetaclient grantee pubkey (%s)", err)
	}
	_, err = common.GetAddressFromPubkeyString(msg.ZetaclientGranteePubkey)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid zetaclient grantee pubkey (%s)", err)
	}
	return nil
}
