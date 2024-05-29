package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgDisableCCTX = "disable_crosschain"
)

var _ sdk.Msg = &MsgDisableCCTX{}

func NewMsgDisableCCTX(creator string, disableOutbound, disableInbound bool) *MsgDisableCCTX {
	return &MsgDisableCCTX{
		Creator:         creator,
		DisableInbound:  disableInbound,
		DisableOutbound: disableOutbound,
	}
}

func (msg *MsgDisableCCTX) Route() string {
	return RouterKey
}

func (msg *MsgDisableCCTX) Type() string {
	return TypeMsgDisableCCTX
}

func (msg *MsgDisableCCTX) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableCCTX) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableCCTX) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !msg.DisableInbound && !msg.DisableOutbound {
		return cosmoserrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"at least one of DisableInbound or DisableOutbound must be true",
		)
	}
	return nil
}
