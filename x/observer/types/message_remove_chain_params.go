package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"github.com/zeta-chain/zetacore/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveChainParams = "remove_chain_params"

var _ sdk.Msg = &MsgRemoveChainParams{}

func NewMsgRemoveChainParams(creator string, chainID int64) *MsgRemoveChainParams {
	return &MsgRemoveChainParams{
		Creator: creator,
		ChainId: chainID,
	}
}

func (msg *MsgRemoveChainParams) Route() string {
	return RouterKey
}

func (msg *MsgRemoveChainParams) Type() string {
	return TypeMsgRemoveChainParams
}

func (msg *MsgRemoveChainParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveChainParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveChainParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check if chain exists
	chain := common.GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidChainID, "invalid chain id (%d)", msg.ChainId)
	}

	return nil
}
