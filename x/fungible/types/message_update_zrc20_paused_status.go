package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateZRC20PausedStatus = "update_zrc20_withdraw_fee"

var _ sdk.Msg = &MsgUpdateZRC20PausedStatus{}

func NewMsgUpdateZRC20PausedStatus(creator string, zrc20 []string, action UpdatePausedStatusAction) *MsgUpdateZRC20PausedStatus {
	return &MsgUpdateZRC20PausedStatus{
		Creator:        creator,
		Zrc20Addresses: zrc20,
		Action:         action,
	}
}

func (msg *MsgUpdateZRC20PausedStatus) Route() string {
	return RouterKey
}

func (msg *MsgUpdateZRC20PausedStatus) Type() string {
	return TypeMsgUpdateZRC20PausedStatus
}

func (msg *MsgUpdateZRC20PausedStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateZRC20PausedStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateZRC20PausedStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Action != UpdatePausedStatusAction_PAUSE && msg.Action != UpdatePausedStatusAction_UNPAUSE {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid action (%d)", msg.Action)
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
