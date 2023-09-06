package types

import (
	cosmoserror "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateContractBytecode = "update_contract_bytecode"

var _ sdk.Msg = &MsgUpdateContractBytecode{}

func NewMsgUpdateContractBytecode(
	creator string, contractAddress ethcommon.Address, newBytecodeAddress ethcommon.Address,
) *MsgUpdateContractBytecode {
	return &MsgUpdateContractBytecode{
		Creator:            creator,
		ContractAddress:    contractAddress.Hex(),
		NewBytecodeAddress: newBytecodeAddress.Hex(),
	}
}

func (msg *MsgUpdateContractBytecode) Route() string {
	return RouterKey
}

func (msg *MsgUpdateContractBytecode) Type() string {
	return TypeMsgUpdateContractBytecode
}

func (msg *MsgUpdateContractBytecode) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateContractBytecode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateContractBytecode) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// check if the contract address is valid
	if ethcommon.HexToAddress(msg.ContractAddress) == (ethcommon.Address{}) {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.ContractAddress)
	}

	// check if the bytecode contract address is valid
	if ethcommon.HexToAddress(msg.NewBytecodeAddress) == (ethcommon.Address{}) {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.ContractAddress)
	}

	return nil
}
