package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSendVoter{}

func NewMsgCreateSendVoter(creator string, index string, sender string, senderChainId string, receiver string, receiverChainId string, mBurnt string, message string, txHash string, blockHeight uint64) *MsgCreateSendVoter {
	return &MsgCreateSendVoter{
		Creator:         creator,
		Index:           index,
		Sender:          sender,
		SenderChainId:   senderChainId,
		Receiver:        receiver,
		ReceiverChainId: receiverChainId,
		MBurnt:          mBurnt,
		Message:         message,
		TxHash:          txHash,
		BlockHeight:     blockHeight,
	}
}

func (msg *MsgCreateSendVoter) Route() string {
	return RouterKey
}

func (msg *MsgCreateSendVoter) Type() string {
	return "CreateSendVoter"
}

func (msg *MsgCreateSendVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateSendVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSendVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
