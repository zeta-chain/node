package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
				Creator:               sample.AccAddress(),
				IsInboundEnabled:      false,
				IsOutboundEnabled:     false,
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
				GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
					EpochLength:             1,
					RetryInterval:           1,
					GasPriceIncreasePercent: 1,
					MaxPendingCctxs:         100,
				},
			},
			want: authoritytypes.PolicyType_groupOperational,
		},
		{
			name: "enabling inbound asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:               sample.AccAddress(),
				IsInboundEnabled:      true,
				IsOutboundEnabled:     false,
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupOperational,
		},
		{
			name: "enabling outbound asserts group 2",
			msg: types.MsgUpdateCrosschainFlags{
				Creator:               sample.AccAddress(),
				IsInboundEnabled:      false,
				IsOutboundEnabled:     true,
				GasPriceIncreaseFlags: nil,
			},
			want: authoritytypes.PolicyType_groupOperational,
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
		msg    *types.MsgUpdateCrosschainFlags
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.NewMsgUpdateCrosschainFlags(
				signer,
				true,
				true,
			),
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.NewMsgUpdateCrosschainFlags(
				"invalid",
				true,
				true,
			),
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
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
	require.Equal(t, types.TypeMsgUpdateCrosschainFlags, msg.Type())
}

func TestMsgUpdateCrosschainFlags_Route(t *testing.T) {
	msg := types.MsgUpdateCrosschainFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateCrosschainFlags_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateCrosschainFlags{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
