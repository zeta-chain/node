package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"github.com/zeta-chain/zetacore/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveCoreParams = "remove_core_params"

var _ sdk.Msg = &MsgRemoveCoreParams{}

func NewMsgRemoveCoreParams(creator string, chainID int64) *MsgRemoveCoreParams {
	return &MsgRemoveCoreParams{
		Creator: creator,
		ChainId: chainID,
	}
}

func (msg *MsgRemoveCoreParams) Route() string {
	return RouterKey
}

func (msg *MsgRemoveCoreParams) Type() string {
	return TypeMsgRemoveCoreParams
}

func (msg *MsgRemoveCoreParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveCoreParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveCoreParams) ValidateBasic() error {
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
