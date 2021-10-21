package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSend{}

func NewMsgCreateSend(creator string, index string, sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64, outTxHash string, outBlockHeight uint64) *MsgCreateSend {
	return &MsgCreateSend{
		Creator:        creator,
		Index:          index,
		Sender:         sender,
		SenderChain:    senderChain,
		Receiver:       receiver,
		ReceiverChain:  receiverChain,
		MBurnt:         mBurnt,
		MMint:          mMint,
		Message:        message,
		InTxHash:       inTxHash,
		InBlockHeight:  inBlockHeight,
		OutTxHash:      outTxHash,
		OutBlockHeight: outBlockHeight,
	}
}

func (msg *MsgCreateSend) Route() string {
	return RouterKey
}

func (msg *MsgCreateSend) Type() string {
	return "CreateSend"
}

func (msg *MsgCreateSend) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateSend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
