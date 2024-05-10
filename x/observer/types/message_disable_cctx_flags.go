package types

import (
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgDisableCCTXFlags = "disable_crosschain_flags"
)

var _ sdk.Msg = &MsgDisableCCTXFlags{}

func NewMsgDisableCCTXFlags(creator string, disableOutbound, disableInbound bool) *MsgDisableCCTXFlags {
	return &MsgDisableCCTXFlags{
		Creator:         creator,
		DisableInbound:  disableInbound,
		DisableOutbound: disableOutbound,
	}
}

func (msg *MsgDisableCCTXFlags) Route() string {
	return RouterKey
}

func (msg *MsgDisableCCTXFlags) Type() string {
	return TypeMsgDisableCCTXFlags
}

func (msg *MsgDisableCCTXFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableCCTXFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableCCTXFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
