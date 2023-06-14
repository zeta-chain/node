package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgVoteOnObservedOutboundTx{}

func NewMsgReceiveConfirmation(creator string, sendHash string, outTxHash string, outBlockHeight uint64, mMint sdk.Uint, status common.ReceiveStatus, chain int64, nonce uint64, coinType common.CoinType) *MsgVoteOnObservedOutboundTx {
	return &MsgVoteOnObservedOutboundTx{
		Creator:                  creator,
		CctxHash:                 sendHash,
		ObservedOutTxHash:        outTxHash,
		ObservedOutTxBlockHeight: outBlockHeight,
		ZetaMinted:               mMint,
		Status:                   status,
		OutTxChain:               chain,
		OutTxTssNonce:            nonce,
		CoinType:                 coinType,
	}
}

func (msg *MsgVoteOnObservedOutboundTx) Route() string {
	return RouterKey
}

func (msg *MsgVoteOnObservedOutboundTx) Type() string {
	return common.OutboundVoter.String()
}

func (msg *MsgVoteOnObservedOutboundTx) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteOnObservedOutboundTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteOnObservedOutboundTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.OutTxChain < 0 {
		return sdkerrors.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.OutTxChain)
	}
	return nil
}

func (msg *MsgVoteOnObservedOutboundTx) Digest() string {
	m := *msg
	m.Creator = ""
	// Set status to ReceiveStatus_Created to make sure both successfull and failed votes are added to the same ballot
	m.Status = common.ReceiveStatus_Created
	// Outbound and reverted txs have different digest as ObservedOutTxHash is different so they are stored in different ballots
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
