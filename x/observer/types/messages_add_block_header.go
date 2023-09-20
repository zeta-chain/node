package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgAddBlockHeader{}

const (
	TypeMsgAddBlockHeader = "add_block_header"
)

func NewMsgAddBlockHeader(creator string, chainID int64, blockHash []byte, height int64, header []byte) *MsgAddBlockHeader {
	return &MsgAddBlockHeader{
		Creator:     creator,
		ChainId:     chainID,
		BlockHash:   blockHash,
		Height:      height,
		BlockHeader: header,
	}
}

func (msg *MsgAddBlockHeader) Route() string {
	return RouterKey
}

func (msg *MsgAddBlockHeader) Type() string {
	return TypeMsgAddBlockHeader
}

func (msg *MsgAddBlockHeader) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddBlockHeader) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddBlockHeader) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.BlockHash) > 32 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg.txhash; too long (%d)", len(msg.BlockHash))
	}
	if len(msg.BlockHeader) > 1024 { // on ethereum the block header is ~538 bytes in RLP encoding
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg.blockheader; too long (%d)", len(msg.BlockHeader))
	}
	if common.IsEthereum(msg.ChainId) {
		// RLP encoded block header
		var header types.Header
		err = rlp.DecodeBytes(msg.BlockHeader, &header)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header; cannot decode RLP (%s)", err)
		}
		if err = header.SanityCheck(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header; sanity check failed (%s)", err)
		}
		if bytes.Compare(msg.BlockHash, header.Hash().Bytes()) != 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header; tx hash mismatch")
		}
		if msg.Height != header.Number.Int64() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header; height mismatch")
		}
	} else {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain id (%d)", msg.ChainId)
	}

	return nil
}

func (msg *MsgAddBlockHeader) Digest() string {
	m := *msg
	m.Creator = ""
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}

func (msg *MsgAddBlockHeader) ParentHash() ([]byte, error) {
	var header types.Header
	err := rlp.DecodeBytes(msg.BlockHeader, &header)
	if err != nil {
		return nil, err
	}

	return header.ParentHash.Bytes(), nil
}
