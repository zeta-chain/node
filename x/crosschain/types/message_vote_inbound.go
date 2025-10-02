package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/solana"
)

// MaxMessageLength is the maximum length of a message in a cctx
// TODO: should parameterize the hardcoded max len
// FIXME: should allow this observation and handle errors in the state machine
// https://github.com/zeta-chain/node/issues/862
const MaxMessageLength = 10240

// InboundVoteOption is a function that sets some option on the inbound vote message
type InboundVoteOption func(*MsgVoteInbound)

// WithRevertOptions sets the revert options for inbound vote message
func WithRevertOptions(revertOptions RevertOptions) InboundVoteOption {
	return func(msg *MsgVoteInbound) {
		msg.RevertOptions = revertOptions
	}
}

// WithZEVMRevertOptions sets the revert options for the inbound vote message (ZEVM format)
// the function convert the type from abigen to type defined in proto
func WithZEVMRevertOptions(revertOptions gatewayzevm.RevertOptions) InboundVoteOption {
	return func(msg *MsgVoteInbound) {
		msg.RevertOptions = NewRevertOptionsFromZEVM(revertOptions)
	}
}

// WithEVMRevertOptions sets the revert options for the inbound vote message (EVM format)
// the function convert the type from abigen to type defined in proto
func WithEVMRevertOptions(revertOptions gatewayevm.RevertOptions) InboundVoteOption {
	return func(msg *MsgVoteInbound) {
		msg.RevertOptions = NewRevertOptionsFromEVM(revertOptions)
	}
}

// WithSOLRevertOptions sets the revert options for the inbound vote message (SOL format)
// the function convert the type from solana instruction to type defined in proto
func WithSOLRevertOptions(revertOptions solana.RevertOptions) InboundVoteOption {
	return func(msg *MsgVoteInbound) {
		msg.RevertOptions = NewRevertOptionsFromSOL(revertOptions)
	}
}

// WithCrossChainCall sets the cross chain call to true for the inbound vote message
func WithCrossChainCall(isCrossChainCall bool) InboundVoteOption {
	return func(msg *MsgVoteInbound) {
		msg.IsCrossChainCall = isCrossChainCall
	}
}

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
	eventIndex uint64,
	protocolContractVersion ProtocolContractVersion,
	isArbitraryCall bool,
	status InboundStatus,
	confirmationMode ConfirmationMode,
	options ...InboundVoteOption,
) *MsgVoteInbound {
	msg := &MsgVoteInbound{
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
		CallOptions: &CallOptions{
			GasLimit:        gasLimit,
			IsArbitraryCall: isArbitraryCall,
		},
		CoinType:                coinType,
		Asset:                   asset,
		EventIndex:              eventIndex,
		ProtocolContractVersion: protocolContractVersion,
		RevertOptions:           NewEmptyRevertOptions(),
		IsCrossChainCall:        false,
		Status:                  status,
		ConfirmationMode:        confirmationMode,
	}

	for _, option := range options {
		option(msg)
	}

	return msg
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

// EligibleForFastConfirmation determines if the inbound msg is eligible for fast confirmation
func (msg *MsgVoteInbound) EligibleForFastConfirmation() bool {
	// only asset CoinType is eligible for fast confirmation
	if !msg.CoinType.IsAsset() {
		return false
	}

	switch msg.ProtocolContractVersion {
	case ProtocolContractVersion_V1:
		// msg using protocol contract version V1 is not eligible for fast confirmation because:
		// 1. whether the receiver address is a contract or not is unknown
		// 2. it can be a depositAndCall (Gas or ZRC20) with empty payload
		// 3. it can be a message passing (CoinType_Zeta) calls 'onReceive'
		return false
	case ProtocolContractVersion_V2:
		// in protocol contract version V2, simple deposit is distinguished from depositAndCall/NoAssetCall
		return !msg.IsCrossChainCall
	default:
		return false
	}
}

// InboundTracker creates an InboundTracker for the inbound vote message
func (msg *MsgVoteInbound) InboundTracker() InboundTracker {
	return InboundTracker{
		ChainId:  msg.SenderChainId,
		TxHash:   msg.InboundHash,
		CoinType: msg.CoinType,
	}
}
