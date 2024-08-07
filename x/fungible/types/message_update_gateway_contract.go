package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateGatewayContract = "update_gateway_contract"

var _ sdk.Msg = &MsgUpdateGatewayContract{}

func NewMsgUpdateGatewayContract(creator string, gatewayContractAddr string) *MsgUpdateGatewayContract {
	return &MsgUpdateGatewayContract{
		Creator:                   creator,
		NewGatewayContractAddress: gatewayContractAddr,
	}
}

func (msg *MsgUpdateGatewayContract) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGatewayContract) Type() string {
	return TypeMsgUpdateGatewayContract
}

func (msg *MsgUpdateGatewayContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGatewayContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGatewayContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	if ethcommon.HexToAddress(msg.NewGatewayContractAddress) == (ethcommon.Address{}) {
		return cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid gateway contract address (%s)",
			msg.NewGatewayContractAddress,
		)
	}

	return nil
}
