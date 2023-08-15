package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

const TypeMsgAddObserver = "add_observer"

var _ sdk.Msg = &MsgAddObserver{}

func NewMsgAddObserver(creator string, observerAdresss, zetaclientGranteeAddress, zetaclientGranteePubKey string) *MsgAddObserver {
	return &MsgAddObserver{
		Creator:                  creator,
		ObserverAddress:          observerAdresss,
		ZetaclientGranteeAddress: zetaclientGranteeAddress,
		ZetaclientGranteePubkey:  zetaclientGranteePubKey,
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
	_, err = sdk.AccAddressFromBech32(msg.ZetaclientGranteeAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid zetaclient grantee address (%s)", err)
	}
	_, err = common.NewPubKey(msg.ZetaclientGranteePubkey)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid zetaclient grantee pubkey (%s)", err)
	}
	// https://github.com/zeta-chain/node/issues/988
	return nil
}
