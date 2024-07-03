package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgDisableHeaderVerification = "disable_header_verification"
)

var _ sdk.Msg = &MsgDisableHeaderVerification{}

func NewMsgDisableHeaderVerification(creator string, chainIDs []int64) *MsgDisableHeaderVerification {
	return &MsgDisableHeaderVerification{
		Creator:     creator,
		ChainIdList: chainIDs,
	}
}

func (msg *MsgDisableHeaderVerification) Route() string {
	return RouterKey
}

func (msg *MsgDisableHeaderVerification) Type() string {
	return TypeMsgDisableHeaderVerification
}

func (msg *MsgDisableHeaderVerification) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableHeaderVerification) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableHeaderVerification) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.ChainIdList) == 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list cannot be empty")
	}

	if len(msg.ChainIdList) > MaxChainIDListLength {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list too long")
	}

	return nil
}
