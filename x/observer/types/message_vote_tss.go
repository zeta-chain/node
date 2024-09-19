package types

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/node/pkg/chains"
)

const TypeMsgVoteTSS = "VoteTSS"

var _ sdk.Msg = &MsgVoteTSS{}

func NewMsgVoteTSS(creator string, pubkey string, keygenZetaHeight int64, status chains.ReceiveStatus) *MsgVoteTSS {
	return &MsgVoteTSS{
		Creator:          creator,
		TssPubkey:        pubkey,
		KeygenZetaHeight: keygenZetaHeight,
		Status:           status,
	}
}

func (msg *MsgVoteTSS) Route() string {
	return RouterKey
}

func (msg *MsgVoteTSS) Type() string {
	return TypeMsgVoteTSS
}

func (msg *MsgVoteTSS) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteTSS) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteTSS) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// either success or observation failure
	if msg.Status != chains.ReceiveStatus_success && msg.Status != chains.ReceiveStatus_failed {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid status: %s", msg.Status)
	}

	return nil
}

func (msg *MsgVoteTSS) Digest() string {
	// We support only 1 keygen at a particular height
	return fmt.Sprintf("%d-%s-%s", msg.KeygenZetaHeight, msg.TssPubkey, "tss-keygen")
}
