package types

import (
	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgAddBlockHeader{}

const (
	TypeMsgAddBlockHeader = "add_block_header"
)

func NewMsgAddBlockHeader(creator string, chainID int64, blockHash []byte, height int64, header common.HeaderData) *MsgAddBlockHeader {
	return &MsgAddBlockHeader{
		Creator:   creator,
		ChainId:   chainID,
		BlockHash: blockHash,
		Height:    height,
		Header:    header,
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
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if common.IsEthereumChain(msg.ChainId) {
		if len(msg.BlockHash) > 32 {
			return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg.txhash; too long (%d)", len(msg.BlockHash))
		}
	} else {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain id (%d)", msg.ChainId)
	}

	if err := msg.Header.Validate(msg.BlockHash, msg.Height); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header (%s)", err)
	}

	if _, err := msg.Header.ParentHash(); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "can't get parent hash (%s)", err)
	}

	return nil
}

func (msg *MsgAddBlockHeader) Digest() string {
	m := *msg
	m.Creator = ""
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
