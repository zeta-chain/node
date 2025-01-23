package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/x/fungible/types"
)

const TypeMsgWhitelistERC20 = "whitelist_erc20"

var _ sdk.Msg = &MsgWhitelistERC20{}

func NewMsgWhitelistERC20(
	creator string,
	erc20Address string,
	chainID int64,
	name string,
	symbol string,
	decimals uint32,
	gasLimit int64,
	liquidityCap sdkmath.Uint,
) *MsgWhitelistERC20 {
	return &MsgWhitelistERC20{
		Creator:      creator,
		Erc20Address: erc20Address,
		ChainId:      chainID,
		Name:         name,
		Symbol:       symbol,
		Decimals:     decimals,
		GasLimit:     gasLimit,
		LiquidityCap: liquidityCap,
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
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err.Error())
	}
	if err := validateAssetAddress(msg.Erc20Address); err != nil {
		return cosmoserrors.Wrapf(types.ErrInvalidAddress, "invalid asset address (%s)", err.Error())
	}
	if msg.Decimals > 128 {
		return cosmoserrors.Wrapf(types.ErrInvalidDecimals, "invalid decimals (%d)", msg.Decimals)
	}
	if msg.GasLimit <= 0 {
		return cosmoserrors.Wrapf(types.ErrInvalidGasLimit, "invalid gas limit (%d)", msg.GasLimit)
	}
	if msg.LiquidityCap.IsNil() {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "liquidity cap is nil")
	}

	return nil
}

func validateAssetAddress(address string) error {
	if address == "" {
		return errors.New("empty asset address")
	}

	// if the address is an evm address, check if it is in checksum format
	if crypto.IsEVMAddress(address) && !crypto.IsChecksumAddress(address) {
		return errors.New("evm address is not in checksum format")
	}

	// currently no specific check is implemented for other address format
	return nil
}
