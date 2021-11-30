package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgGasBalanceVoter{}

func NewMsgGasBalanceVoter(creator string, chain string, balance string, blockNumber uint64) *MsgGasBalanceVoter {
	return &MsgGasBalanceVoter{
		Creator:     creator,
		Chain:       chain,
		Balance:     balance,
		BlockNumber: blockNumber,
	}
}

func (msg *MsgGasBalanceVoter) Route() string {
	return RouterKey
}

func (msg *MsgGasBalanceVoter) Type() string {
	return "GasBalanceVoter"
}

func (msg *MsgGasBalanceVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGasBalanceVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGasBalanceVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
