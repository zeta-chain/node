package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgEnableHeaderVerification = "enable_header_verification"

	// MaxChainIDListLength is the maximum number of chain IDs that can be enabled for header verification
	// this is a value chosen arbitrarily to prevent abuse
	MaxChainIDListLength = 200
)

var _ sdk.Msg = &MsgEnableHeaderVerification{}

func NewMsgEnableHeaderVerification(creator string, chainIDs []int64) *MsgEnableHeaderVerification {
	return &MsgEnableHeaderVerification{
		Creator:     creator,
		ChainIdList: chainIDs,
	}
}

func (msg *MsgEnableHeaderVerification) Route() string {
	return RouterKey
}

func (msg *MsgEnableHeaderVerification) Type() string {
	return TypeMsgEnableHeaderVerification
}

func (msg *MsgEnableHeaderVerification) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgEnableHeaderVerification) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgEnableHeaderVerification) ValidateBasic() error {
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
