package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateGatewayGasLimit = "update_gateway_gas_limit"

var _ sdk.Msg = &MsgUpdateGatewayGasLimit{}

func NewMsgUpdateGatewayGasLimit(creator string, newGasLimit sdkmath.Int) *MsgUpdateGatewayGasLimit {
	return &MsgUpdateGatewayGasLimit{
		Creator:     creator,
		NewGasLimit: newGasLimit,
	}
}

func (msg *MsgUpdateGatewayGasLimit) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGatewayGasLimit) Type() string {
	return TypeMsgUpdateGatewayGasLimit
}

func (msg *MsgUpdateGatewayGasLimit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGatewayGasLimit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGatewayGasLimit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.NewGasLimit.IsNil() || msg.NewGasLimit.IsZero() || msg.NewGasLimit.IsNegative() {
		return cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid gas limit (%s) - must be positive",
			msg.NewGasLimit.String(),
		)
	}

	return nil
}
