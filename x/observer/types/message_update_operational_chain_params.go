package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateOperationalChainParams = "update_operational_chain_params"

var _ sdk.Msg = &MsgUpdateChainParams{}

func NewMsgUpdateOperationalChainParams(
	creator string,
	chainID int64,
	gasPriceTicker,
	inboundTicker,
	outboundTicker,
	watchUtxoTicker uint64,
	outboundScheduleInterval,
	outboundScheduleLookahead int64,
	confirmationParams ConfirmationParams,
) *MsgUpdateOperationalChainParams {
	return &MsgUpdateOperationalChainParams{
		Creator:                   creator,
		ChainId:                   chainID,
		GasPriceTicker:            gasPriceTicker,
		InboundTicker:             inboundTicker,
		OutboundTicker:            outboundTicker,
		WatchUtxoTicker:           watchUtxoTicker,
		OutboundScheduleInterval:  outboundScheduleInterval,
		OutboundScheduleLookahead: outboundScheduleLookahead,
		ConfirmationParams:        confirmationParams,
	}
}

func (msg *MsgUpdateOperationalChainParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateOperationalChainParams) Type() string {
	return TypeMsgUpdateOperationalChainParams
}

func (msg *MsgUpdateOperationalChainParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateOperationalChainParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateOperationalChainParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.ChainId < 0 {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidChainID, "chain id cannot be negative")
	}

	if msg.OutboundScheduleInterval < 0 {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, "outbound schedule interval cannot be negative")
	}

	if msg.OutboundScheduleLookahead < 0 {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, "outbound schedule lookahead cannot be negative")
	}

	if err := msg.ConfirmationParams.Validate(); err != nil {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}
