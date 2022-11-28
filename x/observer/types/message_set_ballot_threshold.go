package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetBallotThreshold = "set_ballot_threshold"

var _ sdk.Msg = &MsgSetBallotThreshold{}

func NewMsgSetBallotThreshold(creator string, chain string, threshold string) *MsgSetBallotThreshold {
	return &MsgSetBallotThreshold{
		Creator:   creator,
		Chain:     chain,
		Threshold: threshold,
	}
}

func (msg *MsgSetBallotThreshold) Route() string {
	return RouterKey
}

func (msg *MsgSetBallotThreshold) Type() string {
	return TypeMsgSetBallotThreshold
}

func (msg *MsgSetBallotThreshold) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetBallotThreshold) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetBallotThreshold) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	threshold, err := sdk.NewDecFromStr(msg.Threshold)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid threshold (%s): %s", err, msg.Threshold)
	}
	if threshold.GT(sdk.OneDec()) || threshold.LT(sdk.ZeroDec()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid threshold: %s", msg.Threshold)
	}
	return nil
}
