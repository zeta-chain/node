package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateGatewayGasLimit = "update_gateway_gas_limit"

	// GatewayGasLimitMax is a max value that can be set, it is arbitrarily chosen with a value that would never be set in practice (30M)
	GatewayGasLimitMax = uint64(30_000_000)
)

var _ sdk.Msg = &MsgUpdateGatewayGasLimit{}

func NewMsgUpdateGatewayGasLimit(creator string, newGasLimit uint64) *MsgUpdateGatewayGasLimit {
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
	if msg.NewGasLimit == 0 {
		return cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid gas limit (%d) - can't be zero",
			msg.NewGasLimit,
		)
	}
	if msg.NewGasLimit > GatewayGasLimitMax {
		return cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid gas limit (%d) - exceeds max allowed (%d)",
			msg.NewGasLimit,
			GatewayGasLimitMax,
		)
	}

	return nil
}
