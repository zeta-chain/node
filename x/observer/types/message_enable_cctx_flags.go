package types

import (
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgEnableCCTXFlags = "enable_crosschain_flags"
)

var _ sdk.Msg = &MsgEnableCCTXFlags{}

func NewMsgEnableCCTXFlags(creator string, isInboundEnabled, isOutboundEnabled bool) *MsgEnableCCTXFlags {
	return &MsgEnableCCTXFlags{
		Creator:        creator,
		EnableInbound:  isInboundEnabled,
		EnableOutbound: isOutboundEnabled,
	}
}

func (msg *MsgEnableCCTXFlags) Route() string {
	return RouterKey
}

func (msg *MsgEnableCCTXFlags) Type() string {
	return TypeMsgUpdateCrosschainFlags
}

func (msg *MsgEnableCCTXFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgEnableCCTXFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgEnableCCTXFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
