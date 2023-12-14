package types

import (
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateCoreParams(msg.CoreParams); err != nil {
		return cosmoserrors.Wrapf(ErrInvalidCoreParams, err.Error())
	}

	return nil
}
