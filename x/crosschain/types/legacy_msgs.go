package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/authz"
)

// MsgWhitelistERC20

var _ sdk.Msg = &MsgWhitelistERC20{}

func (msg *MsgWhitelistERC20) Route() string {
	return RouterKey
}

func (msg *MsgWhitelistERC20) Type() string {
	return "whitelist_erc20"
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
	return nil
}

// MsgVoteOnObservedInboundTx

var _ sdk.Msg = &MsgVoteOnObservedInboundTx{}

func (msg *MsgVoteOnObservedInboundTx) Route() string {
	return RouterKey
}

func (msg *MsgVoteOnObservedInboundTx) Type() string {
	return authz.InboundVoter.String()
}

func (msg *MsgVoteOnObservedInboundTx) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteOnObservedInboundTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteOnObservedInboundTx) ValidateBasic() error {
	return nil
}

// MsgVoteOnObservedOutboundTx

var _ sdk.Msg = &MsgVoteOnObservedOutboundTx{}

func (msg *MsgVoteOnObservedOutboundTx) Route() string {
	return RouterKey
}

func (msg *MsgVoteOnObservedOutboundTx) Type() string {
	return authz.OutboundVoter.String()
}

func (msg *MsgVoteOnObservedOutboundTx) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteOnObservedOutboundTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteOnObservedOutboundTx) ValidateBasic() error {
	return nil
}

// MsgAddToInTxTracker

var _ sdk.Msg = &MsgAddToInTxTracker{}

func (msg *MsgAddToInTxTracker) Route() string {
	return RouterKey
}

func (msg *MsgAddToInTxTracker) Type() string {
	return "AddToInTxTracker"
}

func (msg *MsgAddToInTxTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddToInTxTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddToInTxTracker) ValidateBasic() error {
	return nil
}

// MsgAddToOutTxTracker

var _ sdk.Msg = &MsgAddToOutTxTracker{}

func (msg *MsgAddToOutTxTracker) Route() string {
	return RouterKey
}

func (msg *MsgAddToOutTxTracker) Type() string {
	return "AddToOutTxTracker"
}

func (msg *MsgAddToOutTxTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddToOutTxTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddToOutTxTracker) ValidateBasic() error {
	return nil
}

// MsgRemoveFromOutTxTracker

var _ sdk.Msg = &MsgRemoveFromOutTxTracker{}

func (msg *MsgRemoveFromOutTxTracker) Route() string {
	return RouterKey
}

func (msg *MsgRemoveFromOutTxTracker) Type() string {
	return "RemoveFromOutTxTracker"
}

func (msg *MsgRemoveFromOutTxTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveFromOutTxTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveFromOutTxTracker) ValidateBasic() error {
	return nil
}

// MsgGasPriceVoter

var _ sdk.Msg = &MsgGasPriceVoter{}

func (msg *MsgGasPriceVoter) Route() string {
	return RouterKey
}

func (msg *MsgGasPriceVoter) Type() string {
	return "GasPriceVoter"
}

func (msg *MsgGasPriceVoter) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGasPriceVoter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGasPriceVoter) ValidateBasic() error {
	return nil
}

// MsgMigrateERC20CustodyFunds

var _ sdk.Msg = &MsgMigrateERC20CustodyFunds{}

func (msg *MsgMigrateERC20CustodyFunds) Route() string {
	return RouterKey
}

func (msg *MsgMigrateERC20CustodyFunds) Type() string {
	return "MigrateERC20CustodyFunds"
}

func (msg *MsgMigrateERC20CustodyFunds) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMigrateERC20CustodyFunds) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateERC20CustodyFunds) ValidateBasic() error {
	return nil
}

// MsgUpdateERC20CustodyPauseStatus

var _ sdk.Msg = &MsgUpdateERC20CustodyPauseStatus{}

func (msg *MsgUpdateERC20CustodyPauseStatus) Route() string {
	return RouterKey
}

func (msg *MsgUpdateERC20CustodyPauseStatus) Type() string {
	return "UpdateERC20CustodyPauseStatus"
}

func (msg *MsgUpdateERC20CustodyPauseStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateERC20CustodyPauseStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateERC20CustodyPauseStatus) ValidateBasic() error {
	return nil
}
