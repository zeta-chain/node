package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetZetaGasPool = "set_zeta_gas_pool"

var _ sdk.Msg = &MsgSetZetaGasPool{}

func NewMsgSetZetaGasPool(creator string, address string, poolType string, chain string) *MsgSetZetaGasPool {
	return &MsgSetZetaGasPool{
		Creator:  creator,
		Address:  address,
		PoolType: poolType,
		Chain:    chain,
	}
}

func (msg *MsgSetZetaGasPool) Route() string {
	return RouterKey
}

func (msg *MsgSetZetaGasPool) Type() string {
	return TypeMsgSetZetaGasPool
}

func (msg *MsgSetZetaGasPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetZetaGasPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetZetaGasPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
