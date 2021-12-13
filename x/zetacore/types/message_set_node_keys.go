package types

import (
	"github.com/zeta-chain/zetacore/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetNodeKeys{}

func NewMsgSetNodeKeys(creator string, pubkeySet common.PubKeySet, validatorConsensusPubkey string) *MsgSetNodeKeys {
	return &MsgSetNodeKeys{
		Creator:                  creator,
		PubkeySet:                &pubkeySet,
		ValidatorConsensusPubkey: validatorConsensusPubkey,
	}
}

func (msg *MsgSetNodeKeys) Route() string {
	return RouterKey
}

func (msg *MsgSetNodeKeys) Type() string {
	return "SetNodeKeys"
}

func (msg *MsgSetNodeKeys) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetNodeKeys) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetNodeKeys) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
