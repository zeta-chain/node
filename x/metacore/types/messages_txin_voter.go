package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTxinVoter{}

func NewMsgCreateTxinVoter(creator string, txHash string, sourceAsset string, sourceAmount uint64, mBurnt uint64, destinationAsset string, fromAddress string, toAddress string, blockHeight uint64) *MsgCreateTxinVoter {
	return &MsgCreateTxinVoter{
		Creator:          creator,
		Index:            fmt.Sprintf("%s-%s", txHash, creator),
		TxHash:           txHash,
		SourceAsset:      sourceAsset,
		SourceAmount:     sourceAmount,
		MBurnt:           mBurnt,
		DestinationAsset: destinationAsset,
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		BlockHeight:      blockHeight,
	}
}

func (msg *MsgCreateTxinVoter) Route() string {
	return RouterKey
}

func (msg *MsgCreateTxinVoter) Type() string {
	return "CreateTxinVoter"
}

func (msg *MsgCreateTxinVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTxinVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTxinVoter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	wantedIndex := fmt.Sprintf("%s-%s", msg.TxHash, msg.Creator)
	if msg.Index != wantedIndex {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "msg.Index must be (%s), instead got (%s)", wantedIndex, msg.Index)
	}
	// TODO: Validate the addresses, amounts, and asset format
	return nil
}
