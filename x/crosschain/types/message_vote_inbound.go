package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/coin"
)

// MaxMessageLength is the maximum length of a message in a cctx
// TODO: should parameterize the hardcoded max len
// FIXME: should allow this observation and handle errors in the state machine
// https://github.com/zeta-chain/node/issues/862
const MaxMessageLength = 10240

var _ sdk.Msg = &MsgVoteInbound{}

func NewMsgVoteInbound(
	creator,
	sender string,
	senderChain int64,
	txOrigin,
	receiver string,
	receiverChain int64,
	amount math.Uint,
	message,
	inboundHash string,
	inboundBlockHeight,
	gasLimit uint64,
	coinType coin.CoinType,
	asset string,
	eventIndex uint,
) *MsgVoteInbound {
	return &MsgVoteInbound{
		Creator:            creator,
		Sender:             sender,
		SenderChainId:      senderChain,
		TxOrigin:           txOrigin,
		Receiver:           receiver,
		ReceiverChain:      receiverChain,
		Amount:             amount,
		Message:            message,
		InboundHash:        inboundHash,
		InboundBlockHeight: inboundBlockHeight,
		GasLimit:           gasLimit,
		CoinType:           coinType,
		Asset:              asset,
		EventIndex:         uint64(eventIndex),
	}
}

func (msg *MsgVoteInbound) Route() string {
	return RouterKey
}

func (msg *MsgVoteInbound) Type() string {
	return authz.InboundVoter.String()
}

func (msg *MsgVoteInbound) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteInbound) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteInbound) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s): %s", err, msg.Creator)
	}

	if msg.SenderChainId < 0 {
		return cosmoserrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.SenderChainId)
	}

	if msg.ReceiverChain < 0 {
		return cosmoserrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ReceiverChain)
	}

	if len(msg.Message) > MaxMessageLength {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "message is too long: %d", len(msg.Message))
	}

	return nil
}

func (msg *MsgVoteInbound) Digest() string {
	m := *msg
	m.Creator = ""
	m.InboundBlockHeight = 0
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
