package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeployERC20Mirror = "deploy_erc_20_mirror"

var _ sdk.Msg = &MsgDeployERC20Mirror{}

func NewMsgDeployERC20Mirror(creator string, homeChain string, homeERC20ContractAddress string, name string, symbol string, decimals string) *MsgDeployERC20Mirror {
	return &MsgDeployERC20Mirror{
		Creator:                  creator,
		HomeChain:                homeChain,
		HomeERC20ContractAddress: homeERC20ContractAddress,
		Name:                     name,
		Symbol:                   symbol,
		Decimals:                 decimals,
	}
}

func (msg *MsgDeployERC20Mirror) Route() string {
	return RouterKey
}

func (msg *MsgDeployERC20Mirror) Type() string {
	return TypeMsgDeployERC20Mirror
}

func (msg *MsgDeployERC20Mirror) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeployERC20Mirror) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeployERC20Mirror) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
