package types

import (
	cosmoserror "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateZRC20WithdrawFee = "update_zrc20_withdraw_fee"

var _ sdk.Msg = &MsgUpdateZRC20WithdrawFee{}

func NewMsgUpdateZRC20WithdrawFee(creator string, zrc20 string, newFee sdk.Uint) *MsgUpdateZRC20WithdrawFee {
	return &MsgUpdateZRC20WithdrawFee{
		Creator:        creator,
		Zrc20Address:   zrc20,
		NewWithdrawFee: newFee,
	}
}

func (msg *MsgUpdateZRC20WithdrawFee) Route() string {
	return RouterKey
}

func (msg *MsgUpdateZRC20WithdrawFee) Type() string {
	return TypeMsgUpdateZRC20WithdrawFee
}

func (msg *MsgUpdateZRC20WithdrawFee) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateZRC20WithdrawFee) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateZRC20WithdrawFee) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	if !ethcommon.IsHexAddress(msg.Zrc20Address) {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s)", msg.Zrc20Address)
	}
	if msg.NewWithdrawFee.IsNil() && msg.NewGasLimit.IsNil() {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidRequest, "nothing to update")
	}

	return nil
}
