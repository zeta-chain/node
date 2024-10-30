package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/node/x/fungible/types"
)

const TypeMsgWhitelistERC20 = "whitelist_erc20"

var _ sdk.Msg = &MsgWhitelistERC20{}

func NewMsgWhitelistERC20(
	creator string, erc20Address string, chainID int64, name string,
	symbol string, decimals uint32, gasLimit int64) *MsgWhitelistERC20 {
	return &MsgWhitelistERC20{
		Creator:      creator,
		Erc20Address: erc20Address,
		ChainId:      chainID,
		Name:         name,
		Symbol:       symbol,
		Decimals:     decimals,
		GasLimit:     gasLimit,
	}
}

func (msg *MsgWhitelistERC20) Route() string {
	return RouterKey
}

func (msg *MsgWhitelistERC20) Type() string {
	return TypeMsgWhitelistERC20
}

func (msg *MsgWhitelistERC20) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWhitelistERC20) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWhitelistERC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Erc20Address == "" {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "empty asset address")
	}
	if msg.Decimals > 128 {
		return cosmoserrors.Wrapf(types.ErrInvalidDecimals, "invalid decimals (%d)", msg.Decimals)
	}
	if msg.GasLimit <= 0 {
		return cosmoserrors.Wrapf(types.ErrInvalidGasLimit, "invalid gas limit (%d)", msg.GasLimit)
	}
	return nil
}
