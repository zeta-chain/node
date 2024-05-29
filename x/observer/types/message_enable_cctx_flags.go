package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgEnableCCTX = "enable_crosschain"
)

var _ sdk.Msg = &MsgEnableCCTX{}

func NewMsgEnableCCTX(creator string, enableInbound, enableOutbound bool) *MsgEnableCCTX {
	return &MsgEnableCCTX{
		Creator:        creator,
		EnableInbound:  enableInbound,
		EnableOutbound: enableOutbound,
	}
}

func (msg *MsgEnableCCTX) Route() string {
	return RouterKey
}

func (msg *MsgEnableCCTX) Type() string {
	return TypeMsgEnableCCTX
}

func (msg *MsgEnableCCTX) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgEnableCCTX) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgEnableCCTX) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !msg.EnableInbound && !msg.EnableOutbound {
		return cosmoserrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"at least one of EnableInbound or EnableOutbound must be true",
		)
	}
	return nil
}
