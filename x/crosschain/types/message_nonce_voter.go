package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgNonceVoter{}

func NewMsgNonceVoter(creator string, chain int64, nonce uint64) *MsgNonceVoter {
	return &MsgNonceVoter{
		Creator: creator,
		ChainId: chain,
		Nonce:   nonce,
	}
}

func (msg *MsgNonceVoter) Route() string {
	return RouterKey
}

func (msg *MsgNonceVoter) Type() string {
	return common.NonceVoter.String()
}

func (msg *MsgNonceVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNonceVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNonceVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChainId < 0 {
		return sdkerrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}

	return nil
}
