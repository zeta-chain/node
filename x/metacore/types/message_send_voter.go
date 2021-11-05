package types

import (
	"github.com/Meta-Protocol/metacore/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
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

	receiverChain, err := common.ParseChain(msg.ReceiverChain)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidChainID, "invalid receiver chain (%s): %s", err, msg.ReceiverChain)
	}
	_, err = common.NewAddress(msg.Receiver, receiverChain)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s): %s", err, msg.Receiver)
	}

	if _, ok := big.NewInt(0).SetString(msg.MBurnt, 10); !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot convert mburnt to amount %s: %s", err, msg.MBurnt)
	}
	if msg.MMint != "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "mmint must be empty: %s",  msg.MMint)
	}
	// TODO: should parameterize the hardcoded max len
	if len(msg.Message) > 1024 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "message is too long: %d", len(msg.Message))
	}

	return nil
}

func (msg *MsgSendVoter) Digest() string {
	m := *msg
	m.Creator = ""
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
