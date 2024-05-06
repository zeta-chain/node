package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveOutboundTracker = "RemoveOutboundTracker"

var _ sdk.Msg = &MsgRemoveOutboundTracker{}

func NewMsgRemoveOutboundTracker(creator string, chain int64, nonce uint64) *MsgRemoveOutboundTracker {
	return &MsgRemoveOutboundTracker{
		Creator: creator,
		ChainId: chain,
		Nonce:   nonce,
	}
}

func (msg *MsgRemoveOutboundTracker) Route() string {
	return RouterKey
}

func (msg *MsgRemoveOutboundTracker) Type() string {
	return TypeMsgRemoveOutboundTracker
}

func (msg *MsgRemoveOutboundTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveOutboundTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveOutboundTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChainId < 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	return nil
}
