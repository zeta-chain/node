package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgVoteOnObservedInboundTx{}

func NewMsgSendVoter(creator string, sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64, gasLimit uint64) *MsgVoteOnObservedInboundTx {
	return &MsgVoteOnObservedInboundTx{
		Creator:       creator,
		Sender:        sender,
		SenderChain:   senderChain,
		Receiver:      receiver,
		ReceiverChain: receiverChain,
		ZetaBurnt:     mBurnt,
		Message:       message,
		InTxHash:      inTxHash,
		InBlockHeight: inBlockHeight,
		GasLimit:      gasLimit,
	}
}

func (msg *MsgVoteOnObservedInboundTx) Route() string {
	return RouterKey
}

func (msg *MsgVoteOnObservedInboundTx) Type() string {
	return "InBoundTXVoter"
}

func (msg *MsgVoteOnObservedInboundTx) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteOnObservedInboundTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteOnObservedInboundTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s): %s", err, msg.Creator)
	}
	senderChain, err := common.ParseChain(msg.SenderChain)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidChainID, "invalid sender chain (%s): %s", err, msg.SenderChain)
	}
	_, err = common.NewAddress(msg.Sender, senderChain)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s): %s", err, msg.Sender)
	}
	recvChain, err := common.ParseChain(msg.ReceiverChain)
	if err != nil {
		return fmt.Errorf("cannot parse receiver chain %s", msg.ReceiverChain)
	}
	_, err = common.NewAddress(msg.Receiver, recvChain)
	if err != nil {
		return fmt.Errorf("cannot parse receiver addr %s", msg.Receiver)
	}

	// TODO: should parameterize the hardcoded max len
	// FIXME: should allow this observation and handle errors in the state machine
	if len(msg.Message) > 10240 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "message is too long: %d", len(msg.Message))
	}

	return nil
}

func (msg *MsgVoteOnObservedInboundTx) Digest() string {
	m := *msg
	m.Creator = ""
	m.InBlockHeight = 0
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
