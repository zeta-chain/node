package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSendVoter{}

func NewMsgSendVoter(creator string, sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight string) *MsgSendVoter {
	return &MsgSendVoter{
		Creator:       creator,
		Sender:        sender,
		SenderChain:   senderChain,
		Receiver:      receiver,
		ReceiverChain: receiverChain,
		MBurnt:        mBurnt,
		MMint:         mMint,
		Message:       message,
		InTxHash:      inTxHash,
		InBlockHeight: inBlockHeight,
	}
}

func (msg *MsgSendVoter) Route() string {
	return RouterKey
}

func (msg *MsgSendVoter) Type() string {
	return "SendVoter"
}

func (msg *MsgSendVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSendVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
