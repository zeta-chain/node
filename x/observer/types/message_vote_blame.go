package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

const TypeMsgVoteBlame = "vote_blame"

var _ sdk.Msg = &MsgVoteBlame{}

func NewMsgVoteBlameMsg(creator string, chainID int64, blameInfo Blame) *MsgVoteBlame {
	return &MsgVoteBlame{
		Creator:   creator,
		ChainId:   chainID,
		BlameInfo: blameInfo,
	}
}

func (m *MsgVoteBlame) Route() string {
	return RouterKey
}

func (m *MsgVoteBlame) Type() string {
	return TypeMsgVoteBlame
}

func (m *MsgVoteBlame) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

func (m *MsgVoteBlame) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgVoteBlame) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgVoteBlame) Digest() string {
	msg := *m
	msg.Creator = ""
	// Generate an Identifier for the ballot corresponding to specific blame data
	hash := crypto.Keccak256Hash([]byte(msg.String()))
	return hash.Hex()
}
