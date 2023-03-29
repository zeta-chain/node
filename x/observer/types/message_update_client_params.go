package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateClientParams = "update_client_params"

var _ sdk.Msg = &MsgUpdateClientParams{}

func NewMsgUpdateClientParams(creator string, chainId int64, clientParams *ClientParams) *MsgUpdateClientParams {
	return &MsgUpdateClientParams{
		Creator:      creator,
		ChainId:      chainId,
		ClientParams: clientParams,
	}
}

func (msg *MsgUpdateClientParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateClientParams) Type() string {
	return TypeMsgUpdateClientParams
}

func (msg *MsgUpdateClientParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateClientParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateClientParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ClientParams.ConfirmationCount == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ConfirmationCount must be greater than 0")
	}
	if msg.ClientParams.GasPriceTicker == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "GasPriceTicker must be greater than 0")
	}
	return nil
}
