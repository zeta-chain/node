package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

var _ sdk.Msg = &MsgSetNodeKeys{}

func NewMsgSetNodeKeys(creator string, pubkeySet common.PubKeySet, tssSignerAddress string) *MsgSetNodeKeys {
	return &MsgSetNodeKeys{
		Creator:           creator,
		PubkeySet:         &pubkeySet,
		TssSigner_Address: tssSignerAddress,
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
	accAddressCreator, err := sdk.AccAddressFromBech32(msg.TssSigner_Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid tss signer address (%s)", err)
	}
	pubkey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, msg.PubkeySet.Secp256k1.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPubKeySet, err.Error())
	}
	if bytes.Compare(accAddressCreator.Bytes(), pubkey.Address().Bytes()) != 0 {
		return sdkerrors.Wrapf(ErrInvalidPubKeySet, fmt.Sprintf("Creator : %s , PubkeySet %s", accAddressCreator.String(), pubkey.Address().String()))
	}
	_, err = sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
