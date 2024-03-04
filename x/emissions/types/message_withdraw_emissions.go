package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const MsgWithdrawEmissionType = "withdraw_emission"

var _ sdk.Msg = &MsgWithdrawEmission{}

// NewMsgWithdrawEmission creates a new MsgWithdrawEmission instance

func NewMsgWithdrawEmissions(creator string, amount sdkmath.Int) *MsgWithdrawEmission {
	return &MsgWithdrawEmission{Creator: creator, Amount: amount}
}

func (msg *MsgWithdrawEmission) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawEmission) Type() string {
	return MsgWithdrawEmissionType
}

func (msg *MsgWithdrawEmission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawEmission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawEmission) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Amount.IsNil() || !msg.Amount.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidAmount, "withdraw amount : (%s)", msg.Amount.String())
	}
	return nil
}
