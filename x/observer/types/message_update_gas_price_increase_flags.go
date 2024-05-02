package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateGasPriceIncreaseFlags = "update_gas_price_increase_flags"
)

var _ sdk.Msg = &MsgUpdateGasPriceIncreaseFlags{}

func NewMsgUpdateGasPriceIncreaseFlags(creator string, flags GasPriceIncreaseFlags) *MsgUpdateGasPriceIncreaseFlags {
	return &MsgUpdateGasPriceIncreaseFlags{
		Creator:               creator,
		GasPriceIncreaseFlags: flags,
	}
}

func (msg *MsgUpdateGasPriceIncreaseFlags) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGasPriceIncreaseFlags) Type() string {
	return TypeMsgUpdateGasPriceIncreaseFlags
}

func (msg *MsgUpdateGasPriceIncreaseFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGasPriceIncreaseFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGasPriceIncreaseFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	err := msg.GasPriceIncreaseFlags.Validate()
	if err != nil {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (gpf GasPriceIncreaseFlags) Validate() error {
	if gpf.EpochLength <= 0 {
		return errors.New("epoch length must be positive")
	}
	if gpf.RetryInterval <= 0 {
		return errors.New("retry interval must be positive")
	}
	return nil
}
