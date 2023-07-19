package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

const TypeMsgAddBlameVote = "add_blame_vote"

var _ sdk.Msg = &MsgAddBlameVote{}

func NewMsgAddBlameVoteMsg(creator string, chainId int64, blameInfo *Blame) *MsgAddBlameVote {
	return &MsgAddBlameVote{
		Creator:   creator,
		ChainId:   chainId,
		BlameInfo: blameInfo,
	}
}

func (m *MsgAddBlameVote) Route() string {
	return RouterKey
}

func (m *MsgAddBlameVote) Type() string {
	return TypeMsgAddBlameVote
}

func (m *MsgAddBlameVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if m.ChainId < 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidChainID, "chain id (%d)", m.ChainId)
	}
	return nil
}

func (m *MsgAddBlameVote) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgAddBlameVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAddBlameVote) Digest() string {
	m.Creator = ""
	// Generate an Identifier for the ballot corresponding to specific blame data
	hash := crypto.Keccak256Hash([]byte(m.String()))
	return hash.Hex()
}
