package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveForeignCoin = "remove_foreign_coin"

var _ sdk.Msg = &MsgRemoveForeignCoin{}

func NewMsgRemoveForeignCoin(creator string, zrc20Address string) *MsgRemoveForeignCoin {
	return &MsgRemoveForeignCoin{
		Creator:      creator,
		ZRC20Address: zrc20Address,
	}
}

func (msg *MsgRemoveForeignCoin) Route() string {
	return RouterKey
}

func (msg *MsgRemoveForeignCoin) Type() string {
	return TypeMsgRemoveForeignCoin
}

func (msg *MsgRemoveForeignCoin) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveForeignCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveForeignCoin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
