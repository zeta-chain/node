package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateVerificationFlags = "update_verification_flags"
)

var _ sdk.Msg = &MsgUpdateVerificationFlags{}

func NewMsgUpdateVerificationFlags(creator string, ethTypeChainEnabled, btcTypeChainEnabled bool) *MsgUpdateVerificationFlags {
	return &MsgUpdateVerificationFlags{
		Creator: creator,
		VerificationFlags: VerificationFlags{
			EthTypeChainEnabled: ethTypeChainEnabled,
			BtcTypeChainEnabled: btcTypeChainEnabled,
		},
	}
}

func (msg *MsgUpdateVerificationFlags) Route() string {
	return RouterKey
}

func (msg *MsgUpdateVerificationFlags) Type() string {
	return TypeMsgUpdateVerificationFlags
}

func (msg *MsgUpdateVerificationFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateVerificationFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateVerificationFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
