package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/node/pkg/chains"
)

const TypeMsgUpdateChainInfo = "UpdateChainInfo"

var _ sdk.Msg = &MsgUpdateChainInfo{}

func NewMsgUpdateChainInfo(creator string, chain chains.Chain) *MsgUpdateChainInfo {
	return &MsgUpdateChainInfo{
		Creator: creator,
		Chain:   chain,
	}
}

func (msg *MsgUpdateChainInfo) Route() string {
	return RouterKey
}

func (msg *MsgUpdateChainInfo) Type() string {
	return TypeMsgUpdateChainInfo
}

func (msg *MsgUpdateChainInfo) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgUpdateChainInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateChainInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// the chain must be valid
	if err := msg.Chain.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain"+
			": %s", err.Error())
	}

	return nil
}
