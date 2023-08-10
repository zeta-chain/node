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

func NewMsgAddBlockHeader(creator string, chainId int64, txHash []byte, height int64, header []byte) *MsgAddBlockHeader {
	return &MsgAddBlockHeader{
		Creator:     creator,
		ChainId:     chainId,
		TxHash:      txHash,
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

	if len(msg.TxHash) > 32 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg.txhash; too long (%d)", len(msg.TxHash))
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
		if bytes.Compare(msg.TxHash, header.Hash().Bytes()) != 0 {
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
