package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeploySystemContracts = "deploy_system_contract"

var _ sdk.Msg = &MsgDeploySystemContracts{}

func NewMsgDeploySystemContracts(creator string) *MsgDeploySystemContracts {
	return &MsgDeploySystemContracts{
		Creator: creator,
	}
}

func (msg *MsgDeploySystemContracts) Route() string {
	return RouterKey
}

func (msg *MsgDeploySystemContracts) Type() string {
	return TypeMsgDeployFungibleCoinZRC20
}

func (msg *MsgDeploySystemContracts) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeploySystemContracts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeploySystemContracts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
