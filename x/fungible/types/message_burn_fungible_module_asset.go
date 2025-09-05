package types

import (
	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgBurnFungibleModuleAsset = "burn_fungible_module_asset"

var _ sdk.Msg = &MsgBurnFungibleModuleAsset{}

func NewMsgBurnFungibleModuleAsset(
	creator string,
	zrc20 string,
) *MsgBurnFungibleModuleAsset {
	return &MsgBurnFungibleModuleAsset{
		Creator:      creator,
		Zrc20Address: zrc20,
	}
}

func (msg *MsgBurnFungibleModuleAsset) Route() string {
	return RouterKey
}

func (msg *MsgBurnFungibleModuleAsset) Type() string {
	return TypeMsgBurnFungibleModuleAsset
}

func (msg *MsgBurnFungibleModuleAsset) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurnFungibleModuleAsset) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnFungibleModuleAsset) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the zrc20 address is valid
	if !ethcommon.IsHexAddress(msg.Zrc20Address) {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s)", msg.Zrc20Address)
	}

	return nil
}
