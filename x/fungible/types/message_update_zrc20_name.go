package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateZRC20Name = "update_zrc20_name"

var _ sdk.Msg = &MsgUpdateZRC20Name{}

func NewMsgUpdateZRC20Name(creator, zrc20, name, symbol string) *MsgUpdateZRC20Name {
	return &MsgUpdateZRC20Name{
		Creator:      creator,
		Zrc20Address: zrc20,
		Name:         name,
		Symbol:       symbol,
	}
}

func (msg *MsgUpdateZRC20Name) Route() string {
	return RouterKey
}

func (msg *MsgUpdateZRC20Name) Type() string {
	return TypeMsgUpdateZRC20Name
}

func (msg *MsgUpdateZRC20Name) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateZRC20Name) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateZRC20Name) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !ethcommon.IsHexAddress(msg.Zrc20Address) {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.Zrc20Address)
	}

	if msg.Name == "" && msg.Symbol == "" {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "nothing to update")
	}

	return nil
}
