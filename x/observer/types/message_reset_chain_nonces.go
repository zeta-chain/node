package types

import (
	"errors"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgResetChainNonces = "reset_chain_nonces"

var _ sdk.Msg = &MsgResetChainNonces{}

func NewMsgResetChainNonces(
	creator string,
	chainID int64,
	chainNonceLow int64,
	chainNonceHigh int64,
) *MsgResetChainNonces {
	return &MsgResetChainNonces{
		Creator:        creator,
		ChainId:        chainID,
		ChainNonceLow:  chainNonceLow,
		ChainNonceHigh: chainNonceHigh,
	}
}

func (msg *MsgResetChainNonces) Route() string {
	return RouterKey
}

func (msg *MsgResetChainNonces) Type() string {
	return TypeMsgResetChainNonces
}

func (msg *MsgResetChainNonces) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgResetChainNonces) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgResetChainNonces) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.ChainNonceLow < 0 {
		return errors.New("chain nonce low must be greater or equal 0")
	}

	if msg.ChainNonceHigh < 0 {
		return errors.New("chain nonce high must be greater or equal 0")
	}

	if msg.ChainNonceLow > msg.ChainNonceHigh {
		return errors.New("chain nonce low must be less or equal than chain nonce high")
	}

	return nil
}
