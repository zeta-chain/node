package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

var _ sdk.Msg = &MsgCreateTSSVoter{}

func NewMsgCreateTSSVoter(creator string, chain string, address string, pubkey string) *MsgCreateTSSVoter {
	return &MsgCreateTSSVoter{
		Creator: creator,
		Chain:   chain,
		Address: address,
		Pubkey:  pubkey,
	}
}

func (msg *MsgCreateTSSVoter) Route() string {
	return RouterKey
}

func (msg *MsgCreateTSSVoter) Type() string {
	return "CreateTSSVoter"
}

func (msg *MsgCreateTSSVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTSSVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTSSVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid voter address (%s)", err)
	}
	return nil
}

func (msg *MsgCreateTSSVoter) Digest() string {
	m := *msg
	m.Creator = ""
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
