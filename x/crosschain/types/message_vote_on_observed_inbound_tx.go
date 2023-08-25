package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgVoteOnObservedInboundTx{}

func NewMsgVoteOnObservedInboundTx(
	creator,
	sender string,
	senderChain int64,
	txOrigin,
	receiver string,
	receiverChain int64,
	amount math.Uint,
	message,
	inTxHash string,
	inBlockHeight,
	gasLimit uint64,
	coinType common.CoinType,
	asset string,
) *MsgVoteOnObservedInboundTx {
	return &MsgVoteOnObservedInboundTx{
		Creator:       creator,
		Sender:        sender,
		SenderChainId: senderChain,
		TxOrigin:      txOrigin,
		Receiver:      receiver,
		ReceiverChain: receiverChain,
		Amount:        amount,
		Message:       message,
		InTxHash:      inTxHash,
		InBlockHeight: inBlockHeight,
		GasLimit:      gasLimit,
		CoinType:      coinType,
		Asset:         asset,
	}
}

func (msg *MsgVoteOnObservedInboundTx) Route() string {
	return RouterKey
}

func (msg *MsgVoteOnObservedInboundTx) Type() string {
	return common.InboundVoter.String()
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
	if msg.SenderChainId < 0 {
		return sdkerrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.SenderChainId)
	}

	if msg.ReceiverChain < 0 {
		return sdkerrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ReceiverChain)
	}

	// TODO: should parameterize the hardcoded max len
	// FIXME: should allow this observation and handle errors in the state machine
	// https://github.com/zeta-chain/node/issues/862
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
