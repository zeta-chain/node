package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common/ethereum"
)

const (
	TypeMsgProveOutboundTx = "prove_outbound_tx"
)

var _ sdk.Msg = &MsgProveOutboundTx{}

func NewMsgProveOutboundTx(creator string, txProof ethereum.Proof, receiptProof ethereum.Proof, blockHash string, txIndex int64) *MsgProveOutboundTx {
	return &MsgProveOutboundTx{
		Creator:      creator,
		TxProof:      &txProof,
		ReceiptProof: &receiptProof,
		BlockHash:    blockHash,
		TxIndex:      txIndex,
	}
}

func (msg *MsgProveOutboundTx) Route() string {
	return RouterKey
}

func (msg *MsgProveOutboundTx) Type() string {
	return TypeMsgProveOutboundTx
}

func (msg *MsgProveOutboundTx) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgProveOutboundTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgProveOutboundTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
