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

const TypeMsgWhitelistAsset = "whitelist_asset"

var _ sdk.Msg = &MsgWhitelistAsset{}

func NewMsgWhitelistAsset(
	creator string,
	assetAddress string,
	chainID int64,
	name string,
	symbol string,
	decimals uint32,
	gasLimit int64,
	liquidityCap sdkmath.Uint,
) *MsgWhitelistAsset {
	return &MsgWhitelistAsset{
		Creator:      creator,
		AssetAddress: assetAddress,
		ChainId:      chainID,
		Name:         name,
		Symbol:       symbol,
		Decimals:     decimals,
		GasLimit:     gasLimit,
		LiquidityCap: liquidityCap,
	}
}

func (msg *MsgWhitelistAsset) Route() string {
	return RouterKey
}

func (msg *MsgWhitelistAsset) Type() string {
	return TypeMsgWhitelistAsset
}

func (msg *MsgWhitelistAsset) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWhitelistAsset) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWhitelistAsset) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err.Error())
	}
	if err := validateAssetAddress(msg.AssetAddress); err != nil {
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
