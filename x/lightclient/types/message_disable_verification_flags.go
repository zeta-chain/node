package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

const (
	TypeMsgDisableVerificationFlags = "disable_verification_flags"
)

var _ sdk.Msg = &MsgDisableVerificationFlags{}

func NewMsgDisableVerificationFlags(creator string, chainIDs []int64) *MsgDisableVerificationFlags {
	return &MsgDisableVerificationFlags{
		Creator:     creator,
		ChainIdList: chainIDs,
	}

}

func (msg *MsgDisableVerificationFlags) Route() string {
	return RouterKey
}

func (msg *MsgDisableVerificationFlags) Type() string {
	return TypeMsgDisableVerificationFlags
}

func (msg *MsgDisableVerificationFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDisableVerificationFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDisableVerificationFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	changelistForHeaderSupport := chains.ChainListForHeaderSupport()
	if len(msg.ChainIdList) == 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list cannot be empty")
	}
	if len(msg.ChainIdList) > len(changelistForHeaderSupport) {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "chain id list cannot be greater than supported chains")
	}
	for _, chainID := range msg.ChainIdList {
		if !chains.ChainIDInChainList(chainID, changelistForHeaderSupport) {
			return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain id header not supported (%d)", chainID)
		}
	}

	return nil
}
