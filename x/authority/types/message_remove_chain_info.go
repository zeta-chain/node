package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveChainInfo = "RemoveChainInfo"

var _ sdk.Msg = &MsgRemoveChainInfo{}

func NewMsgRemoveChainInfo(creator string, chainID int64) *MsgRemoveChainInfo {
	return &MsgRemoveChainInfo{
		Creator: creator,
		ChainId: chainID,
	}
}

func (msg *MsgRemoveChainInfo) Route() string {
	return RouterKey
}

func (msg *MsgRemoveChainInfo) Type() string {
	return TypeMsgRemoveChainInfo
}

func (msg *MsgRemoveChainInfo) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveChainInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveChainInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
