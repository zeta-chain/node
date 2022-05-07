package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math"
	"strconv"
)

const TypeMsgZetaConversionRateVoter = "zeta_conversion_rate_voter"

var _ sdk.Msg = &MsgZetaConversionRateVoter{}

func NewMsgZetaConversionRateVoter(creator string, chain string, zetaConversionRate string, blockNumber uint64) *MsgZetaConversionRateVoter {
	return &MsgZetaConversionRateVoter{
		Creator:            creator,
		Chain:              chain,
		ZetaConversionRate: zetaConversionRate,
		BlockNumber:        blockNumber,
	}
}

func (msg *MsgZetaConversionRateVoter) Route() string {
	return RouterKey
}

func (msg *MsgZetaConversionRateVoter) Type() string {
	return TypeMsgZetaConversionRateVoter
}

func (msg *MsgZetaConversionRateVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgZetaConversionRateVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgZetaConversionRateVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.ZetaConversionRate) > 15 {
		return sdkerrors.Wrapf(ErrFloatParseError, "invalid float (%s)", msg.ZetaConversionRate)
	}

	v, err := strconv.ParseFloat(msg.ZetaConversionRate, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return sdkerrors.Wrapf(ErrFloatParseError, "invalid float (%s)", msg.ZetaConversionRate)
	}
	return nil
}
