package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAbortStuckCCTX = "AbortStuckCCTX"

var _ sdk.Msg = &MsgAbortStuckCCTX{}

func NewMsgAbortStuckCCTX(creator string, cctxIndex string) *MsgAbortStuckCCTX {
	return &MsgAbortStuckCCTX{
		Creator:   creator,
		CctxIndex: cctxIndex,
	}
}

func (msg *MsgAbortStuckCCTX) Route() string {
	return RouterKey
}

func (msg *MsgAbortStuckCCTX) Type() string {
	return TypeMsgAbortStuckCCTX
}

func (msg *MsgAbortStuckCCTX) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAbortStuckCCTX) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAbortStuckCCTX) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
