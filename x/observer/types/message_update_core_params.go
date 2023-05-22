package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

const TypeMsgUpdateClientParams = "update_client_params"

var _ sdk.Msg = &MsgUpdateCoreParams{}

func NewMsgUpdateCoreParams(creator string, coreParams *CoreParams) *MsgUpdateCoreParams {
	return &MsgUpdateCoreParams{
		Creator:    creator,
		CoreParams: coreParams,
	}
}

func (msg *MsgUpdateCoreParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCoreParams) Type() string {
	return TypeMsgUpdateClientParams
}

func (msg *MsgUpdateCoreParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCoreParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCoreParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.CoreParams.ConfirmationCount == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "ConfirmationCount must be greater than 0")
	}
	if msg.CoreParams.GasPriceTicker <= 0 || msg.CoreParams.GasPriceTicker > 30 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "GasPriceTicker out of range")
	}
	if msg.CoreParams.InTxTicker <= 0 || msg.CoreParams.InTxTicker > 10 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "InTxTicker out of range")
	}
	if msg.CoreParams.OutTxTicker <= 0 || msg.CoreParams.OutTxTicker > 10 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutTxTicker out of range")
	}
	if common.GetChainFromChainID(msg.CoreParams.ChainId) == nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "ChainId not supported")
	}
	return nil
}
