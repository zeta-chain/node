package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgPauseZrc20 = "pause_zrc20"

var _ sdk.Msg = &MsgPauseZRC20{}

func NewMsgPauseZRC20(creator string, zrc20 []string) *MsgPauseZRC20 {
	return &MsgPauseZRC20{
		Creator:        creator,
		Zrc20Addresses: zrc20,
	}
}

func (msg *MsgPauseZRC20) Route() string {
	return RouterKey
}

func (msg *MsgPauseZRC20) Type() string {
	return TypeMsgPauseZrc20
}

func (msg *MsgPauseZRC20) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPauseZRC20) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPauseZRC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Zrc20Addresses) == 0 {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, "no zrc20 to update")
	}

	// check if all zrc20 addresses are valid
	for _, zrc20 := range msg.Zrc20Addresses {
		if !ethcommon.IsHexAddress(zrc20) {
			return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid zrc20 contract address (%s)", zrc20)
		}
	}
	return nil
}
