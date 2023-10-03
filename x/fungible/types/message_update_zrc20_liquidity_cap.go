package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateZRC20LiquidityCap = "update_zrc20_liquidity_cap"

var _ sdk.Msg = &MsgUpdateZRC20LiquidityCap{}

func NewMsgUpdateZRC20LiquidityCap(creator string, zrc20 string, liquidityCap math.Uint) *MsgUpdateZRC20LiquidityCap {
	return &MsgUpdateZRC20LiquidityCap{
		Creator:      creator,
		Zrc20Address: zrc20,
		LiquidityCap: liquidityCap,
	}
}

func (msg *MsgUpdateZRC20LiquidityCap) Route() string {
	return RouterKey
}

func (msg *MsgUpdateZRC20LiquidityCap) Type() string {
	return TypeMsgUpdateSystemContract
}

func (msg *MsgUpdateZRC20LiquidityCap) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateZRC20LiquidityCap) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateZRC20LiquidityCap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !ethcommon.IsHexAddress(msg.Zrc20Address) {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.Zrc20Address)
	}

	return nil
}
