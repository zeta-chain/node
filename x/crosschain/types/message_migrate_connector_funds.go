package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgMigrateConnectorFunds = "MigrateConnectorFunds"

var _ sdk.Msg = &MsgMigrateConnectorFunds{}

func NewMsgMigrateConnectorFunds(
	creator string,
	chainID int64,
	newConnectorAddress string,
	amount sdkmath.Uint,
) *MsgMigrateConnectorFunds {
	return &MsgMigrateConnectorFunds{
		Creator:             creator,
		ChainId:             chainID,
		NewConnectorAddress: newConnectorAddress,
		Amount:              amount,
	}
}

func (msg *MsgMigrateConnectorFunds) Route() string {
	return RouterKey
}

func (msg *MsgMigrateConnectorFunds) Type() string {
	return TypeMsgMigrateConnectorFunds
}

func (msg *MsgMigrateConnectorFunds) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMigrateConnectorFunds) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateConnectorFunds) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !ethcommon.IsHexAddress(msg.NewConnectorAddress) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid new connector address")
	}

	if msg.Amount.LTE(sdkmath.ZeroUint()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}
	return nil
}
