package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgDisableFastConfirmation = "disable_fast_confirmation"

	// MaxChainIDListLength is the maximum number of chain IDs can be passed in one message
	// this is a value chosen arbitrarily to prevent abuse
	MaxChainIDListLength = 200
)

var _ sdk.Msg = &MsgDisableFastConfirmation{}

func NewMsgDisableFastConfirmation(creator string, chainIDs []int64) *MsgDisableFastConfirmation {
	return &MsgDisableFastConfirmation{
		Creator:     creator,
		ChainIdList: chainIDs,
	}
}

func (msg *MsgDisableFastConfirmation) Route() string {
	return RouterKey
}

func (msg *MsgDisableFastConfirmation) Type() string {
	return TypeMsgDisableFastConfirmation
}

func (msg *MsgDisableFastConfirmation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableFastConfirmation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableFastConfirmation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.ChainIdList) > MaxChainIDListLength {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list too long")
	}

	return nil
}
