package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
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

func TestMsgUpdateCrosschainFlags_GetRequiredGroup(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateCrosschainFlags
		want types.Policy_Type
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
			want: types.Policy_Type_group1,
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
			want: types.Policy_Type_group1,
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
			want: types.Policy_Type_group2,
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
			want: types.Policy_Type_group2,
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
			want: types.Policy_Type_group2,
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
			want: types.Policy_Type_group2,
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
			want: types.Policy_Type_group2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.msg.GetRequiredGroup())
		})
	}
}
