package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgUpdateCrosschainFlags_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateCrosschainFlags
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateCrosschainFlags{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid gas price increase flags",
			msg: types.MsgUpdateCrosschainFlags{
				Creator: sample.AccAddress(),
				GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
					EpochLength:             -1,
					RetryInterval:           1,
					GasPriceIncreasePercent: 1,
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: types.MsgUpdateCrosschainFlags{
				Creator: sample.AccAddress(),
				GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
					EpochLength:             1,
					RetryInterval:           1,
					GasPriceIncreasePercent: 1,
				},
			},
		},
		{
			name: "gas price increase flags can be nil",
			msg: types.MsgUpdateCrosschainFlags{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGasPriceIncreaseFlags_Validate(t *testing.T) {
	tests := []struct {
		name        string
		gpf         types.GasPriceIncreaseFlags
		errContains string
	}{
		{
			name: "invalid epoch length",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             -1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 1,
			},
			errContains: "epoch length must be positive",
		},
		{
			name: "invalid retry interval",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             1,
				RetryInterval:           -1,
				GasPriceIncreasePercent: 1,
			},
			errContains: "retry interval must be positive",
		},
		{
			name: "valid",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 1,
			},
		},
		{
			name: "percent can be 0",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gpf.Validate()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateCrosschainFlags_GetRequiredPolicyType(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateCrosschainFlags
		want authoritytypes.PolicyType
	}{
		{
			name: "disabling outbound and inbound allows group 1",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:                      sample.AccAddress(),
				IsInboundEnabled:             false,
				IsOutboundEnabled:            false,
				BlockHeaderVerificationFlags: nil,
				GasPriceIncreaseFlags:        nil,
			},
			want: authoritytypes.PolicyType_groupEmergency,
		},
		{
			name: "disabling outbound and inbound and block header verification allows group 1",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:           sample.AccAddress(),
				IsInboundEnabled:  false,
				IsOutboundEnabled: false,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: false,
					IsBtcTypeChainEnabled: false,
				},
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupEmergency,
		},
		{
			name: "updating gas price increase flags asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:           sample.AccAddress(),
				IsInboundEnabled:  false,
				IsOutboundEnabled: false,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: false,
					IsBtcTypeChainEnabled: false,
				},
				GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
					EpochLength:             1,
					RetryInterval:           1,
					GasPriceIncreasePercent: 1,
					MaxPendingCctxs:         100,
				},
			},
			want: authoritytypes.PolicyType_groupAdmin,
		},
		{
			name: "enabling inbound asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:           sample.AccAddress(),
				IsInboundEnabled:  true,
				IsOutboundEnabled: false,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: false,
					IsBtcTypeChainEnabled: false,
				},
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupAdmin,
		},
		{
			name: "enabling outbound asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:           sample.AccAddress(),
				IsInboundEnabled:  false,
				IsOutboundEnabled: true,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: false,
					IsBtcTypeChainEnabled: false,
				},
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupAdmin,
		},
		{
			name: "enabling eth header verification asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:           sample.AccAddress(),
				IsInboundEnabled:  false,
				IsOutboundEnabled: false,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: true,
					IsBtcTypeChainEnabled: false,
				},
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupAdmin,
		},
		{
			name: "enabling btc header verification asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:           sample.AccAddress(),
				IsInboundEnabled:  false,
				IsOutboundEnabled: false,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: false,
					IsBtcTypeChainEnabled: true,
				},
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.want, tt.msg.GetRequiredPolicyType())
		})
	}
}

func TestMsgUpdateCrosschainFlags_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateCrosschainFlags
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateCrosschainFlags{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateCrosschainFlags{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				assert.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				assert.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgUpdateCrosschainFlags_Type(t *testing.T) {
	msg := types.MsgUpdateCrosschainFlags{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgUpdateCrosschainFlags, msg.Type())
}

func TestMsgUpdateCrosschainFlags_Route(t *testing.T) {
	msg := types.MsgUpdateCrosschainFlags{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateCrosschainFlags_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateCrosschainFlags{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
