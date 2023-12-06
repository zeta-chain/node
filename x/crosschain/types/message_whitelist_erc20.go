package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
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
	return types.RouterKey
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
	bz := types.ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWhitelistERC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	if ethcommon.HexToAddress(msg.Erc20Address) == (ethcommon.Address{}) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ERC20 contract address (%s)", msg.Erc20Address)
	}
	if msg.Decimals > 77 {
		return sdkerrors.Wrapf(types.ErrInvalidDecimals, "invalid decimals (%d), decimals must be less than 78", msg.Decimals)
	}
	if msg.GasLimit <= 0 {
		return sdkerrors.Wrapf(types.ErrInvalidGasLimit, "invalid gas limit (%d)", msg.GasLimit)
	}
	return nil
}
