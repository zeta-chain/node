package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/node/pkg/proofs"
)

var _ sdk.Msg = &MsgVoteBlockHeader{}

const (
	TypeMsgVoteBlockHeader = "vote_block_header"
)

func NewMsgVoteBlockHeader(
	creator string,
	chainID int64,
	blockHash []byte,
	height int64,
	header proofs.HeaderData,
) *MsgVoteBlockHeader {
	return &MsgVoteBlockHeader{
		Creator:   creator,
		ChainId:   chainID,
		BlockHash: blockHash,
		Height:    height,
		Header:    header,
	}
}

func (msg *MsgVoteBlockHeader) Route() string {
	return RouterKey
}

func (msg *MsgVoteBlockHeader) Type() string {
	return TypeMsgVoteBlockHeader
}

func (msg *MsgVoteBlockHeader) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteBlockHeader) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteBlockHeader) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if len(msg.BlockHash) != 32 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block hash length (%d)", len(msg.BlockHash))
	}

	if err := msg.Header.Validate(msg.BlockHash, msg.ChainId, msg.Height); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header (%s)", err)
	}

	if _, err := msg.Header.ParentHash(); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "can't get parent hash (%s)", err)
	}

	return nil
}

func (msg *MsgVoteBlockHeader) Digest() string {
	m := *msg
	m.Creator = ""
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
