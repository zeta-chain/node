package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgGasPriceVoter{}

func NewMsgGasPriceVoter(creator string, chain int64, price uint64, supply string, blockNumber uint64) *MsgGasPriceVoter {
	return &MsgGasPriceVoter{
		Creator:     creator,
		ChainId:     chain,
		Price:       price,
		BlockNumber: blockNumber,
		Supply:      supply,
	}
}

func (msg *MsgGasPriceVoter) Route() string {
	return RouterKey
}

func (msg *MsgGasPriceVoter) Type() string {
	return common.GasPriceVoter.String()
}

func (msg *MsgGasPriceVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGasPriceVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGasPriceVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChainId < 0 {
		return cosmoserrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	return nil
}
