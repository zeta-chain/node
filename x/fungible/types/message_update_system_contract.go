package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateSystemContract = "update_system_contract"

var _ sdk.Msg = &MsgUpdateSystemContract{}

func NewMsgUpdateSystemContract(creator string, systemContractAddr string) *MsgUpdateSystemContract {
	return &MsgUpdateSystemContract{
		Creator:                  creator,
		NewSystemContractAddress: systemContractAddr,
	}
}

func (msg *MsgUpdateSystemContract) Route() string {
	return RouterKey
}

func (msg *MsgUpdateSystemContract) Type() string {
	return TypeMsgUpdateSystemContract
}

func (msg *MsgUpdateSystemContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateSystemContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateSystemContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	if ethcommon.HexToAddress(msg.NewSystemContractAddress) == (ethcommon.Address{}) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s)", msg.NewSystemContractAddress)
	}

	return nil
}
