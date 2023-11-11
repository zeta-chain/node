package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

var _ sdk.Msg = &MsgUpdateTssAddress{}

func NewMsgUpdateTssAddress(creator string, pubkey string) *MsgUpdateTssAddress {
	return &MsgUpdateTssAddress{
		Creator:   creator,
		TssPubkey: pubkey,
	}
}

func (msg *MsgUpdateTssAddress) Route() string {
	return RouterKey
}

func (msg *MsgUpdateTssAddress) Type() string {
	return "UpdateTssAddress"
}

func (msg *MsgUpdateTssAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateTssAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateTssAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, msg.TssPubkey)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid tss pubkey (%s)", err)
	}
	return nil
}
