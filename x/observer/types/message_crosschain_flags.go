package types

import (
	"errors"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

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

// GetRequiredPolicyType returns the required policy type for the message to execute the message
// Group emergency should only be able to stop or disable functionalities in case of emergency
// this concerns disabling inbound and outbound txs or block header verification
// every other action requires group admin
// TODO: add separate message for each group
// https://github.com/zeta-chain/node/issues/1562
func (msg *MsgUpdateCrosschainFlags) GetRequiredPolicyType() authoritytypes.PolicyType {
	if msg.IsInboundEnabled || msg.IsOutboundEnabled {
		return authoritytypes.PolicyType_groupAdmin
	}
	if msg.GasPriceIncreaseFlags != nil {
		return authoritytypes.PolicyType_groupAdmin
	}
	if msg.BlockHeaderVerificationFlags != nil && (msg.BlockHeaderVerificationFlags.IsEthTypeChainEnabled || msg.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled) {
		return authoritytypes.PolicyType_groupAdmin

	}
	return authoritytypes.PolicyType_groupEmergency
}
