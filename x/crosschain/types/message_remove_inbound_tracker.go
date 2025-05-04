package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveInboundTracker = "RemoveInboundTracker"

var _ sdk.Msg = &MsgRemoveInboundTracker{}

func NewMsgRemoveInboundTracker(creator string, chain int64, txHash string) *MsgRemoveInboundTracker {
	return &MsgRemoveInboundTracker{
		Creator: creator,
		ChainId: chain,
		TxHash:  txHash,
	}
}

func (msg *MsgRemoveInboundTracker) Route() string {
	return RouterKey
}

func (msg *MsgRemoveInboundTracker) Type() string {
	return TypeMsgRemoveInboundTracker
}

func (msg *MsgRemoveInboundTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveInboundTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveInboundTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChainId < 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	return nil
}
