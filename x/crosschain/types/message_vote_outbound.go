package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
)

var _ sdk.Msg = &MsgVoteOutbound{}

func NewMsgVoteOutbound(
	creator,
	cctxIndex,
	outboundHash string,
	outboundBlockHeight,
	outboundGasUsed uint64,
	outboundEffectiveGasPrice math.Int,
	outboundEffectiveGasLimit uint64,
	valueReceived math.Uint,
	status chains.ReceiveStatus,
	chain int64,
	nonce uint64,
	coinType coin.CoinType,
	confirmationMode ConfirmationMode,
) *MsgVoteOutbound {
	return &MsgVoteOutbound{
		Creator:                           creator,
		CctxHash:                          cctxIndex,
		ObservedOutboundHash:              outboundHash,
		ObservedOutboundBlockHeight:       outboundBlockHeight,
		ObservedOutboundGasUsed:           outboundGasUsed,
		ObservedOutboundEffectiveGasPrice: outboundEffectiveGasPrice,
		ObservedOutboundEffectiveGasLimit: outboundEffectiveGasLimit,
		ValueReceived:                     valueReceived,
		Status:                            status,
		OutboundChain:                     chain,
		OutboundTssNonce:                  nonce,
		CoinType:                          coinType,
		ConfirmationMode:                  confirmationMode,
	}
}

func (msg *MsgVoteOutbound) Route() string {
	return RouterKey
}

func (msg *MsgVoteOutbound) Type() string {
	return authz.OutboundVoter.String()
}

func (msg *MsgVoteOutbound) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteOutbound) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteOutbound) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.OutboundChain < 0 {
		return cosmoserrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.OutboundChain)
	}

	return nil
}

func (msg *MsgVoteOutbound) Digest() string {
	m := *msg
	m.Creator = ""
	m.ConfirmationMode = ConfirmationMode_SAFE

	// Set status to ReceiveStatus_created to make sure both successful and failed votes are added to the same ballot
	m.Status = chains.ReceiveStatus_created

	// Outbound and reverted txs have different digest as ObservedOutboundHash is different so they are stored in different ballots
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
