package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/proofs"
)

const TypeMsgAddOutboundTracker = "AddOutboundTracker"

var _ sdk.Msg = &MsgAddOutboundTracker{}

func NewMsgAddOutboundTracker(
	creator string,
	chain int64,
	nonce uint64,
	txHash string,
	proof *proofs.Proof,
	blockHash string,
	txIndex int64,
) *MsgAddOutboundTracker {
	return &MsgAddOutboundTracker{
		Creator:   creator,
		ChainId:   chain,
		Nonce:     nonce,
		TxHash:    txHash,
		Proof:     proof,
		BlockHash: blockHash,
		TxIndex:   txIndex,
	}
}

func (msg *MsgAddOutboundTracker) Route() string {
	return RouterKey
}

func (msg *MsgAddOutboundTracker) Type() string {
	return TypeMsgAddOutboundTracker
}

func (msg *MsgAddOutboundTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddOutboundTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddOutboundTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChainId < 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	return nil
}
