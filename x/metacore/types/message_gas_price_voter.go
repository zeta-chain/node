package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgGasPriceVoter{}

func NewMsgGasPriceVoter(creator string, chain string, price uint64, blockNumber uint64) *MsgGasPriceVoter {
	return &MsgGasPriceVoter{
		Creator:     creator,
		Chain:       chain,
		Price:       price,
		BlockNumber: blockNumber,
	}
}

func (msg *MsgGasPriceVoter) Route() string {
	return RouterKey
}

func (msg *MsgGasPriceVoter) Type() string {
	return "GasPriceVoter"
}

func (msg *MsgGasPriceVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGasPriceVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGasPriceVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
