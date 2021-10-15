package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgTxoutConfirmationVoter{}

func NewMsgTxoutConfirmationVoter(creator string, txoutId uint64, txHash string, mMint uint64, destinationAsset string, destinationAmount uint64, toAddress string, blockHeight uint64) *MsgTxoutConfirmationVoter {
	return &MsgTxoutConfirmationVoter{
		Creator:           creator,
		TxoutId:           txoutId,
		TxHash:            txHash,
		MMint:             mMint,
		DestinationAsset:  destinationAsset,
		DestinationAmount: destinationAmount,
		ToAddress:         toAddress,
		BlockHeight:       blockHeight,
	}
}

func (msg *MsgTxoutConfirmationVoter) Route() string {
	return RouterKey
}

func (msg *MsgTxoutConfirmationVoter) Type() string {
	return "TxoutConfirmationVoter"
}

func (msg *MsgTxoutConfirmationVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTxoutConfirmationVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTxoutConfirmationVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
