package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

const (
	TypeMsgDisableHeaderVerification = "disable_header_verification"
)

var _ sdk.Msg = &MsgDisableHeaderVerification{}

func NewMsgDisableHeaderVerification(creator string, chainIDs []int64) *MsgDisableHeaderVerification {
	return &MsgDisableHeaderVerification{
		Creator:     creator,
		ChainIdList: chainIDs,
	}

}

func (msg *MsgDisableHeaderVerification) Route() string {
	return RouterKey
}

func (msg *MsgDisableHeaderVerification) Type() string {
	return TypeMsgDisableHeaderVerification
}

func (msg *MsgDisableHeaderVerification) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableHeaderVerification) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableHeaderVerification) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	chainListForHeaderSupport := chains.ChainListForHeaderSupport()
	if len(msg.ChainIdList) == 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list cannot be empty")
	}
	if len(msg.ChainIdList) > len(chainListForHeaderSupport) {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list cannot be greater than supported chains")
	}
	for _, chainID := range msg.ChainIdList {
		if !chains.ChainIDInChainList(chainID, chainListForHeaderSupport) {
			return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain id header not supported (%d)", chainID)
		}
	}

	return nil
}
