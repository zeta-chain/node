package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgCreateTSSVoter{}

func NewMsgCreateTSSVoter(creator string, pubkey string, keygenZetaHeight int64, status common.ReceiveStatus) *MsgCreateTSSVoter {
	return &MsgCreateTSSVoter{
		Creator:          creator,
		TssPubkey:        pubkey,
		KeyGenZetaHeight: keygenZetaHeight,
		Status:           status,
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
	return nil
}

func (msg *MsgCreateTSSVoter) Digest() string {
	// We support only 1 keygen at a particular height
	return fmt.Sprintf("%d-%s", msg.KeyGenZetaHeight, msg.TssPubkey)
}
