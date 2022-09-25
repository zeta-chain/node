package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

const TypeMsgDeployFungibleCoinZRC4 = "deploy_fungible_coin_zrc_4"

var _ sdk.Msg = &MsgDeployFungibleCoinZRC4{}

func NewMsgDeployFungibleCoinZRC4(creator string, eRC20 string, foreignChain string, decimals uint32, name string, symbol string, coinType common.CoinType) *MsgDeployFungibleCoinZRC4 {
	return &MsgDeployFungibleCoinZRC4{
		Creator:      creator,
		ERC20:        eRC20,
		ForeignChain: foreignChain,
		Decimals:     decimals,
		Name:         name,
		Symbol:       symbol,
		CoinType:     coinType,
	}
}

func (msg *MsgDeployFungibleCoinZRC4) Route() string {
	return RouterKey
}

func (msg *MsgDeployFungibleCoinZRC4) Type() string {
	return TypeMsgDeployFungibleCoinZRC4
}

func (msg *MsgDeployFungibleCoinZRC4) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeployFungibleCoinZRC4) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeployFungibleCoinZRC4) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
