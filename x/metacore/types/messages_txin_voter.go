package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTxinVoter{}

func NewMsgCreateTxinVoter(creator string, index string, txHash string, sourceAsset string, sourceAmount string, mBurnt string, destinationAsset string, fromAddress string, toAddress string, blockHeight string, signer string, signature string) *MsgCreateTxinVoter {
	return &MsgCreateTxinVoter{
		Creator:          creator,
		Index:            index,
		TxHash:           txHash,
		SourceAsset:      sourceAsset,
		SourceAmount:     sourceAmount,
		MBurnt:           mBurnt,
		DestinationAsset: destinationAsset,
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		BlockHeight:      blockHeight,
		Signer:           signer,
		Signature:        signature,
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
	//TODO: validate the signature.
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
