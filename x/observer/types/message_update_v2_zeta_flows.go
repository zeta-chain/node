package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateV2ZetaFlows = "update_v2_zeta_flows"
)

var _ sdk.Msg = &MsgUpdateV2ZetaFlows{}

func NewMsgUpdateV2ZetaFlows(creator string, isV2ZetaEnabled bool) *MsgUpdateV2ZetaFlows {
	return &MsgUpdateV2ZetaFlows{
		Creator:         creator,
		IsV2ZetaEnabled: isV2ZetaEnabled,
	}
}

func (msg *MsgUpdateV2ZetaFlows) Route() string {
	return RouterKey
}

func (msg *MsgUpdateV2ZetaFlows) Type() string {
	return TypeMsgUpdateV2ZetaFlows
}

func (msg *MsgUpdateV2ZetaFlows) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateV2ZetaFlows) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateV2ZetaFlows) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
