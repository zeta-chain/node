package types

import (
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateChainParams = "update_chain_params"

var _ sdk.Msg = &MsgUpdateChainParams{}

func NewMsgUpdateChainParams(creator string, chainParams *ChainParams) *MsgUpdateChainParams {
	return &MsgUpdateChainParams{
		Creator:     creator,
		ChainParams: chainParams,
	}
}

func (msg *MsgUpdateChainParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateChainParams) Type() string {
	return TypeMsgUpdateChainParams
}

func (msg *MsgUpdateChainParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateChainParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateChainParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateChainParams(msg.ChainParams); err != nil {
		return cosmoserrors.Wrapf(ErrInvalidChainParams, err.Error())
	}

	return nil
}
