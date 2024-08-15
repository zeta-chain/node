package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeUpdateERC20CustodyPauseStatus = "UpdateERC20CustodyPauseStatus"

var _ sdk.Msg = &MsgUpdateERC20CustodyPauseStatus{}

func NewMsgUpdateERC20CustodyPauseStatus(
	creator string,
	chainID int64,
	pause bool,
) *MsgUpdateERC20CustodyPauseStatus {
	return &MsgUpdateERC20CustodyPauseStatus{
		Creator: creator,
		ChainId: chainID,
		Pause:   pause,
	}
}

func (msg *MsgUpdateERC20CustodyPauseStatus) Route() string {
	return RouterKey
}

func (msg *MsgUpdateERC20CustodyPauseStatus) Type() string {
	return TypeUpdateERC20CustodyPauseStatus
}

func (msg *MsgUpdateERC20CustodyPauseStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateERC20CustodyPauseStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateERC20CustodyPauseStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
