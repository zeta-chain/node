package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeployGasPriceOracle = "deploy_gas_price_oracle"

var _ sdk.Msg = &MsgDeployGasPriceOracle{}

func NewMsgDeployGasPriceOracle(creator string) *MsgDeployGasPriceOracle {
	return &MsgDeployGasPriceOracle{
		Creator: creator,
	}
}

func (msg *MsgDeployGasPriceOracle) Route() string {
	return RouterKey
}

func (msg *MsgDeployGasPriceOracle) Type() string {
	return TypeMsgDeployGasPriceOracle
}

func (msg *MsgDeployGasPriceOracle) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeployGasPriceOracle) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeployGasPriceOracle) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
