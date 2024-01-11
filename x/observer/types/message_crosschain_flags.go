package types

import (
	"errors"

	cosmoserrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateCrosschainFlags = "update_crosschain_flags"
)

var _ sdk.Msg = &MsgUpdateCrosschainFlags{}

func NewMsgUpdateCrosschainFlags(creator string, isInboundEnabled, isOutboundEnabled bool) *MsgUpdateCrosschainFlags {
	return &MsgUpdateCrosschainFlags{
		Creator:           creator,
		IsInboundEnabled:  isInboundEnabled,
		IsOutboundEnabled: isOutboundEnabled,
	}
}

func (msg *MsgUpdateCrosschainFlags) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCrosschainFlags) Type() string {
	return TypeMsgUpdateCrosschainFlags
}

func (msg *MsgUpdateCrosschainFlags) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCrosschainFlags) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCrosschainFlags) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.GasPriceIncreaseFlags != nil {
		if err := msg.GasPriceIncreaseFlags.Validate(); err != nil {
			return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	return nil
}

func (gpf GasPriceIncreaseFlags) Validate() error {
	if gpf.EpochLength <= 0 {
		return errors.New("epoch length must be positive")
	}
	if gpf.RetryInterval <= 0 {
		return errors.New("retry interval must be positive")
	}
	return nil
}

// GetRequiredGroup returns the required group policy for the message to execute the message
// Group 1 should only be able to stop or disable functiunalities in case of emergency
// this concerns disabling inbound and outbound txs or block header verification
// every other action requires group 2
func (msg *MsgUpdateCrosschainFlags) GetRequiredGroup() Policy_Type {
	if msg.IsInboundEnabled || msg.IsOutboundEnabled {
		return Policy_Type_group2
	}
	if msg.GasPriceIncreaseFlags != nil {
		return Policy_Type_group2
	}
	if msg.BlockHeaderVerificationFlags != nil && (msg.BlockHeaderVerificationFlags.IsEthTypeChainEnabled || msg.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled) {
		return Policy_Type_group2

	}
	return Policy_Type_group1
}
