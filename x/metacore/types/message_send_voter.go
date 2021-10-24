package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

var _ sdk.Msg = &MsgSendVoter{}

func NewMsgSendVoter(creator string, sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64) *MsgSendVoter {
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
	//TODO: add basic validation
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (msg *MsgSendVoter) Digest() string {
	m := *msg
	m.Creator = ""
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
